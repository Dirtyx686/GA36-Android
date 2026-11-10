package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func newImageReader(r io.Reader) (*gzip.Reader, error) { return gzip.NewReader(r) }
func runGUI(args []string) error {
	exe, e := os.Executable()
	if e != nil {
		return e
	}
	exe, e = filepath.EvalSymlinks(exe)
	if e != nil {
		return e
	}
	if len(args) > 0 && (args[0] == "--gui-demo" || args[0] == "--gui-demo-headless") {
		a, e := newGUIApp(true, exe)
		if e != nil {
			return e
		}
		defer os.RemoveAll(a.temp)
		u, e := a.start()
		if e != nil {
			return e
		}
		if args[0] == "--gui-demo-headless" {
			fmt.Println(u)
		} else if e = guiOpenBrowser(u); e != nil {
			a.stop()
			return e
		}
		<-a.done
		return nil
	}
	if len(args) > 0 && args[0] == "--gui-worker" {
		if !guiIsAdmin() {
			return errors.New("Der Geräte-Dienst benötigt Administratorrechte.")
		}
		if len(args) != 5 || args[1] != "--notify" || args[3] != "--nonce" || len(args[4]) != 64 {
			return errors.New("Ungültiger Startaufruf")
		}
		if e = validateNotify(args[2]); e != nil {
			return e
		}
		a, e := newGUIApp(false, exe)
		if e != nil {
			return e
		}
		defer os.RemoveAll(a.temp)
		u, e := a.start()
		if e != nil {
			return e
		}
		a.defaultImage()
		body, _ := json.Marshal(map[string]string{"url": u})
		req, _ := http.NewRequest("POST", args[2], bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+args[4])
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{Proxy: nil}}
		resp, e := client.Do(req)
		if e != nil {
			a.stop()
			return e
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			a.stop()
			return errors.New("Lokaler Startkanal nicht bestätigt")
		}
		<-a.done
		return nil
	}
	if len(args) > 0 && !(len(args) == 1 && args[0] == "--gui") {
		return errors.New("Unbekannte GUI-Option")
	}
	// The normal-user process opens the browser. Only the background device service elevates.
	ln, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		return e
	}
	base := "http://" + ln.Addr().String()
	nonce := randomHex(32)
	var mu sync.Mutex
	ready := ""
	message := "Bitte die Administrator-Anfrage bestätigen."
	state := "starting"
	done := make(chan struct{})
	var once sync.Once
	finish := func() { once.Do(func() { close(done) }) }
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if r.Host != strings.TrimPrefix(base, "http://") {
			http.Error(w, "Host", 403)
			return
		}
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			io.WriteString(w, bootstrapHTML)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+nonce {
			http.Error(w, "Unauthorized", 401)
			return
		}
		if r.URL.Path == "/status" && r.Method == "GET" {
			mu.Lock()
			defer mu.Unlock()
			jsonReply(w, 200, map[string]string{"status": state, "message": message, "url": ready})
			return
		}
		if r.URL.Path == "/ready" && r.Method == "POST" {
			var q struct {
				URL string `json:"url"`
			}
			if e := readJSON(r, &q); e != nil {
				apiError(w, e)
				return
			}
			u, e := url.Parse(q.URL)
			if e != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" || u.Port() == "" || u.Path != "/" || len(u.Fragment) != 64 {
				http.Error(w, "Invalid endpoint", 400)
				return
			}
			mu.Lock()
			ready = q.URL
			state = "ready"
			message = "Bereit"
			mu.Unlock()
			jsonReply(w, 200, map[string]bool{"ok": true})
			go func() { time.Sleep(6 * time.Second); finish() }()
			return
		}
		http.NotFound(w, r)
	})
	srv := &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	go srv.Serve(ln)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	if e = guiOpenBrowser(base + "/#" + nonce); e != nil {
		return e
	}
	if e = guiElevate(exe, base+"/ready", nonce); e != nil {
		mu.Lock()
		state = "error"
		message = e.Error()
		mu.Unlock()
	}
	select {
	case <-done:
		return nil
	case <-time.After(90 * time.Second):
		mu.Lock()
		if state != "ready" {
			state = "error"
			message = "Der Geräte-Dienst wurde nicht gestartet. Bitte den Flasher erneut öffnen und die Anmeldung bestätigen."
		}
		mu.Unlock()
		time.Sleep(30 * time.Second)
		return nil
	}
}

const bootstrapHTML = `<!doctype html><html lang="de"><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>GA36 Flasher starten</title><style>body{margin:0;background:#f3f5fa;color:#17223a;font:16px system-ui;display:grid;place-items:center;height:100vh}.box{background:white;padding:48px;border-radius:24px;box-shadow:0 20px 80px #16336112;max-width:480px}small{font-weight:700;color:#6855e9}h1{font-size:32px}p{line-height:1.6;color:#56637b}.dot{display:inline-block;width:12px;height:12px;background:#7764ef;border-radius:50%;margin-right:10px}</style><div class="box"><small>GA36 · FLASHER</small><h1>Wir machen alles bereit.</h1><p id="message"><span class="dot"></span>Bitte die Administrator-Anfrage bestätigen. Die Oberfläche öffnet sich gleich.</p><p>Der Geräte-Dienst läuft lokal auf deinem PC. Es wird noch keine Karte beschrieben.</p></div><script>const t=location.hash.slice(1);setInterval(async()=>{try{const r=await fetch('/status',{headers:{Authorization:'Bearer '+t}});const d=await r.json();if(d.status==='ready')location.replace(d.url);else document.getElementById('message').textContent=d.message}catch(e){}},700)</script></html>`
