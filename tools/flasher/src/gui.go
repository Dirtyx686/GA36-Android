package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//go:embed web/*
var guiFiles embed.FS

const guiVersion = "2.0-Teststand"
const maxUpload = int64(4*1024*1024*1024 + 2*1024*1024)

type guiImage struct {
	Manifest Manifest `json:"manifest"`
	Name     string   `json:"name"`
	Path     string   `json:"-"`
	OwnedDir string   `json:"-"`
}
type guiJob struct {
	ID         string    `json:"id"`
	Kind       string    `json:"kind"`
	Phase      string    `json:"phase"`
	Status     string    `json:"status"`
	Current    int64     `json:"current"`
	Total      int64     `json:"total"`
	Error      string    `json:"error,omitempty"`
	Incomplete bool      `json:"incomplete"`
	Started    time.Time `json:"started"`
}
type guiEvent struct {
	At   time.Time `json:"at"`
	Text string    `json:"text"`
}
type guiSnapshot struct {
	Version string     `json:"version"`
	Demo    bool       `json:"demo"`
	Busy    bool       `json:"busy"`
	Image   *guiImage  `json:"image"`
	Job     guiJob     `json:"job"`
	Events  []guiEvent `json:"events"`
}
type guiCard struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Serial string `json:"serial"`
	Bytes  int64  `json:"bytes"`
	Ready  bool   `json:"ready"`
	Reason string `json:"reason,omitempty"`
}
type guiInboxView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type guiApp struct {
	mu                     sync.Mutex
	demo                   bool
	token, base, exe, temp string
	image                  *guiImage
	job                    guiJob
	events                 []guiEvent
	cancel                 context.CancelFunc
	server                 *http.Server
	lastSeen               time.Time
	done                   chan struct{}
	stopOnce               sync.Once
}

func randomHex(n int) string {
	b := make([]byte, n)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func newGUIApp(demo bool, exe string) (*guiApp, error) {
	t, e := os.MkdirTemp("", "ga36-gui-")
	if e != nil {
		return nil, e
	}
	a := &guiApp{demo: demo, token: randomHex(32), exe: exe, temp: t, lastSeen: time.Now(), done: make(chan struct{})}
	if demo {
		var m Manifest
		json.Unmarshal(releaseJSON, &m)
		a.image = &guiImage{Manifest: m, Name: "Android 7 · Demo-Image"}
		root := filepath.Join(t, "demo-inbox")
		os.MkdirAll(filepath.Join(root, "APKS"), 0700)
		os.MkdirAll(filepath.Join(root, "OBB"), 0700)
		os.WriteFile(filepath.Join(root, "GA36-APPS.TXT"), []byte(inboxMarker), 0600)
		os.WriteFile(filepath.Join(root, "INSTALL.TXT"), []byte("DEMO: Keine echte Karte angeschlossen.\n"), 0600)
	}
	a.event("Grafische Oberfläche bereit.")
	return a, nil
}
func (a *guiApp) event(text string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, guiEvent{time.Now(), text})
	if len(a.events) > 150 {
		a.events = a.events[len(a.events)-150:]
	}
}
func (a *guiApp) snapshot() guiSnapshot {
	a.mu.Lock()
	defer a.mu.Unlock()
	var img *guiImage
	if a.image != nil {
		v := *a.image
		img = &v
	}
	events := append([]guiEvent{}, a.events...)
	return guiSnapshot{guiVersion, a.demo, a.job.Status == "running", img, a.job, events}
}
func (a *guiApp) begin(kind, phase string) (context.Context, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.job.Status == "running" {
		return nil, errors.New("Ein Vorgang läuft bereits.")
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	a.job = guiJob{ID: randomHex(8), Kind: kind, Phase: phase, Status: "running", Started: time.Now()}
	return ctx, nil
}
func (a *guiApp) progress(phase string, n, total int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.job.Phase = phase
	a.job.Current = n
	a.job.Total = total
}
func (a *guiApp) finish(e error) {
	a.mu.Lock()
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
	if e != nil {
		a.job.Status = "error"
		a.job.Error = e.Error()
	} else {
		a.job.Status = "done"
		a.job.Error = ""
	}
	kind := a.job.Kind
	a.mu.Unlock()
	if e != nil {
		a.event(kind + ": " + e.Error())
	} else {
		a.event(kind + ": abgeschlossen.")
	}
}
func (a *guiApp) opaque(s string) string {
	h := sha256.Sum256([]byte(a.token + "|" + s))
	return hex.EncodeToString(h[:12])
}
func (a *guiApp) disks() ([]Disk, error) {
	if a.demo {
		return []Disk{{Path: "demo://sd64", Model: "microSD · Kartenleser", Serial: "DEMO-64", Bytes: 64000000000, Safe: true, Identity: "demo64"}, {Path: "demo://sd32", Model: "microSD · zweiter Kartenleser", Serial: "DEMO-32", Bytes: 32000000000, Safe: true, Identity: "demo32"}}, nil
	}
	a.mu.Lock()
	paths := []string{a.exe}
	if a.image != nil && a.image.Path != "" {
		paths = append(paths, a.image.Path)
	}
	a.mu.Unlock()
	return listDisks(paths...)
}
func (a *guiApp) findCard(id string) (Disk, error) {
	ds, e := a.disks()
	if e != nil {
		return Disk{}, e
	}
	for _, d := range ds {
		if d.Safe && a.opaque(d.Identity) == id {
			return d, nil
		}
	}
	return Disk{}, errors.New("Zielkarte nicht mehr unverändert erreichbar. Bitte neu laden.")
}
func (a *guiApp) inboxes() ([]Inbox, error) {
	if a.demo {
		root := filepath.Join(a.temp, "demo-inbox")
		return []Inbox{{Root: root, Description: "microSD · DEMO-64 · GA36_APPS", Open: func() (string, func(), error) { return root, func() {}, nil }}}, nil
	}
	return listInboxes(a.exe)
}
func (a *guiApp) findInbox(id string) (Inbox, error) {
	vs, e := a.inboxes()
	if e != nil {
		return Inbox{}, e
	}
	for _, v := range vs {
		if a.opaque(v.Description) == id {
			return v, nil
		}
	}
	return Inbox{}, errors.New("GA36_APPS nicht gefunden. Kartenleser neu verbinden und Liste aktualisieren.")
}
func jsonReply(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func apiError(w http.ResponseWriter, e error) {
	jsonReply(w, 400, map[string]string{"error": e.Error()})
}
func readJSON(r *http.Request, v any) error {
	d := json.NewDecoder(io.LimitReader(r.Body, 8193))
	d.DisallowUnknownFields()
	return d.Decode(v)
}
func (a *guiApp) serveHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'self'")
	if a.base != "" && r.Host != strings.TrimPrefix(a.base, "http://") {
		http.Error(w, "Host rejected", 403)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/api/") {
		if r.Method != "GET" {
			http.Error(w, "Method", 405)
			return
		}
		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			name = "index.html"
		}
		if name != "index.html" && name != "app.js" && name != "style.css" {
			http.NotFound(w, r)
			return
		}
		b, e := fs.ReadFile(guiFiles, "web/"+name)
		if e != nil {
			http.NotFound(w, r)
			return
		}
		switch name {
		case "index.html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case "app.js":
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		case "style.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		}
		w.Write(b)
		return
	}
	got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if len(got) != len(a.token) || subtle.ConstantTimeCompare([]byte(got), []byte(a.token)) != 1 {
		http.Error(w, "Unauthorized", 401)
		return
	}
	if origin := r.Header.Get("Origin"); origin != "" && origin != a.base {
		http.Error(w, "Origin rejected", 403)
		return
	}
	a.mu.Lock()
	a.lastSeen = time.Now()
	a.mu.Unlock()
	switch r.URL.Path {
	case "/api/state":
		if r.Method != "GET" {
			http.Error(w, "Method", 405)
			return
		}
		jsonReply(w, 200, a.snapshot())
	case "/api/cards":
		if r.Method != "GET" {
			http.Error(w, "Method", 405)
			return
		}
		ds, e := a.disks()
		if e != nil {
			apiError(w, e)
			return
		}
		s := a.snapshot()
		var m Manifest
		json.Unmarshal(releaseJSON, &m)
		if s.Image != nil {
			m = s.Image.Manifest
		}
		out := []guiCard{}
		for _, d := range ds {
			if !d.Safe {
				continue
			}
			ready := eligible(d, m)
			reason := ""
			if !ready {
				reason = "Zu klein oder nicht unterstützt"
			}
			out = append(out, guiCard{a.opaque(d.Identity), d.Model, d.Path, d.Serial, d.Bytes, ready, reason})
		}
		jsonReply(w, 200, out)
	case "/api/inboxes":
		if r.Method != "GET" {
			http.Error(w, "Method", 405)
			return
		}
		vs, e := a.inboxes()
		if e != nil {
			apiError(w, e)
			return
		}
		out := []guiInboxView{}
		for _, v := range vs {
			out = append(out, guiInboxView{a.opaque(v.Description), v.Description})
		}
		jsonReply(w, 200, out)
	case "/api/image":
		if r.Method != "POST" {
			http.Error(w, "Method", 405)
			return
		}
		a.uploadImage(w, r)
	case "/api/flash":
		if r.Method != "POST" {
			http.Error(w, "Method", 405)
			return
		}
		a.startFlash(w, r)
	case "/api/apks", "/api/obbs":
		if r.Method != "POST" {
			http.Error(w, "Method", 405)
			return
		}
		a.uploadImports(w, r, r.URL.Path == "/api/obbs")
	case "/api/log":
		if r.Method != "GET" {
			http.Error(w, "Method", 405)
			return
		}
		v, e := a.findInbox(r.URL.Query().Get("id"))
		if e != nil {
			apiError(w, e)
			return
		}
		root, close, e := v.Open()
		if e != nil {
			apiError(w, e)
			return
		}
		defer close()
		if e = checkInbox(root); e != nil {
			apiError(w, e)
			return
		}
		f, e := os.Open(filepath.Join(root, "INSTALL.TXT"))
		if e != nil {
			apiError(w, e)
			return
		}
		defer f.Close()
		s, e := f.Stat()
		if e != nil {
			apiError(w, e)
			return
		}
		if s.Size() > 256*1024 {
			f.Seek(s.Size()-256*1024, 0)
		}
		b, e := io.ReadAll(io.LimitReader(f, 256*1024))
		if e != nil {
			apiError(w, e)
			return
		}
		jsonReply(w, 200, map[string]string{"text": string(b)})
	case "/api/cancel":
		if r.Method != "POST" {
			http.Error(w, "Method", 405)
			return
		}
		a.mu.Lock()
		if a.cancel != nil {
			a.cancel()
		}
		a.mu.Unlock()
		jsonReply(w, 200, map[string]bool{"ok": true})
	case "/api/quit":
		if r.Method != "POST" {
			http.Error(w, "Method", 405)
			return
		}
		if a.snapshot().Busy {
			jsonReply(w, 409, map[string]string{"error": "Bitte zuerst den laufenden Vorgang beenden."})
			return
		}
		jsonReply(w, 200, map[string]bool{"ok": true})
		go a.stop()
	default:
		http.NotFound(w, r)
	}
}
func (a *guiApp) inspectImage(path, name string, m Manifest, owned string) error {
	if e := validManifest(m); e != nil {
		return e
	}
	if m.ImageBytes != baseImageBytes {
		return errors.New("Dieses Image verwendet ein anderes Kartenformat. Unterstützt wird derzeit das GA36-R01-Layout.")
	}
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()
	if e = validateImage(f, m); e != nil {
		return fmt.Errorf("%w Passende IMAGE-INFOS.json auswählen.", e)
	}
	gz, e := newImageReader(f)
	if e != nil {
		return e
	}
	buf := make([]byte, 4096)
	_, e = io.ReadFull(gz, buf)
	gz.Close()
	if e != nil {
		return e
	}
	h := sha256.Sum256(buf)
	if hex.EncodeToString(h[:]) != m.HeaderSHA {
		return errors.New("Boot-Metadaten stimmen nicht mit den Prüfdaten überein.")
	}
	if _, e = adjustedHeader(buf, baseImageBytes); e != nil {
		return e
	}
	a.mu.Lock()
	old := a.image
	a.image = &guiImage{m, name, path, owned}
	a.mu.Unlock()
	if old != nil && old.OwnedDir != "" && old.OwnedDir != owned {
		os.RemoveAll(old.OwnedDir)
	}
	a.event("Image geprüft: " + name)
	return nil
}
func (a *guiApp) uploadImage(w http.ResponseWriter, r *http.Request) {
	ctx, e := a.begin("Image", "Image wird eingelesen")
	if e != nil {
		apiError(w, e)
		return
	}
	_ = ctx
	defer r.Body.Close()
	dir, e := os.MkdirTemp(a.temp, "image-")
	if e != nil {
		a.finish(e)
		apiError(w, e)
		return
	}
	ok := false
	defer func() {
		if !ok {
			os.RemoveAll(dir)
		}
	}()
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)
	mr, e := r.MultipartReader()
	var m Manifest
	json.Unmarshal(releaseJSON, &m)
	name := ""
	path := filepath.Join(dir, "image.gz")
	if e == nil {
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				e = err
				break
			}
			if p.FormName() == "manifest" {
				b, err := io.ReadAll(io.LimitReader(p, 65537))
				if err != nil || len(b) > 65536 {
					e = errors.New("Prüfdatei zu groß")
					break
				}
				if err = json.Unmarshal(b, &m); err != nil {
					e = err
					break
				}
			} else if p.FormName() == "image" && name == "" {
				name = filepath.Base(p.FileName())
				f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
				if err != nil {
					e = err
					break
				}
				_, err = io.Copy(f, p)
				closeErr := f.Close()
				if err != nil {
					e = err
					break
				}
				if closeErr != nil {
					e = closeErr
					break
				}
			} else {
				e = errors.New("Unerwartete Upload-Datei")
				break
			}
			p.Close()
		}
	}
	if e == nil && name == "" {
		e = errors.New("Kein Image ausgewählt")
	}
	if e != nil {
		a.finish(e)
		apiError(w, e)
		return
	}
	ok = true
	jsonReply(w, 202, map[string]bool{"ok": true})
	go func() {
		a.progress("Prüfsummen und Kartenformat prüfen", 0, 0)
		e := a.inspectImage(path, name, m, dir)
		if e != nil {
			os.RemoveAll(dir)
		}
		a.finish(e)
	}()
}
func (a *guiApp) startFlash(w http.ResponseWriter, r *http.Request) {
	var q struct {
		CardID  string `json:"card_id"`
		Confirm string `json:"confirm"`
	}
	if e := readJSON(r, &q); e != nil {
		apiError(w, e)
		return
	}
	if q.Confirm != "ERASE:"+q.CardID {
		apiError(w, errors.New("Löschbestätigung fehlt."))
		return
	}
	s := a.snapshot()
	if s.Image == nil {
		apiError(w, errors.New("Zuerst ein gültiges Image auswählen."))
		return
	}
	d, e := a.findCard(q.CardID)
	if e != nil {
		apiError(w, e)
		return
	}
	if !eligible(d, s.Image.Manifest) {
		apiError(w, errors.New("Karte ist zu klein."))
		return
	}
	ctx, e := a.begin("System schreiben", "Image prüfen")
	if e != nil {
		apiError(w, e)
		return
	}
	jsonReply(w, 202, map[string]bool{"ok": true})
	go func() {
		if a.demo {
			for _, phase := range []string{"Image prüfen · Demo", "Schreiben · Demo", "Zurücklesen und prüfen · Demo"} {
				for n := int64(0); n <= 100; n += 4 {
					select {
					case <-ctx.Done():
						a.finish(errors.New("Simulation abgebrochen."))
						return
					case <-time.After(55 * time.Millisecond):
					}
					a.progress(phase, n, 100)
				}
			}
			a.event("DEMO: Keine echte Karte beschrieben.")
			a.finish(nil)
			return
		}
		img, e := os.Open(s.Image.Path)
		if e != nil {
			a.finish(e)
			return
		}
		defer img.Close()
		if e = validateImage(img, s.Image.Manifest); e != nil {
			a.finish(e)
			return
		}
		fresh, e := a.findCard(q.CardID)
		if e != nil || fresh.Identity != d.Identity {
			a.finish(errors.New("Karte wurde verändert. Bitte neu laden."))
			return
		}
		dev, e := openDevice(fresh, s.Image.Path, a.exe)
		if e != nil {
			a.finish(e)
			return
		}
		defer dev.Close()
		e = flash(ctx, img, s.Image.Manifest, dev, func(phase string, n, total int64) {
			a.progress(phase, n, total)
			if n > 0 {
				a.mu.Lock()
				a.job.Incomplete = true
				a.mu.Unlock()
			}
		})
		if e == nil {
			a.progress("Datenpartition an Kartengröße anpassen", 0, 0)
			e = configureDataPartition(dev, d.Bytes, s.Image.Manifest.HeaderSHA)
		}
		if e == nil {
			e = dev.Finish()
		}
		if e == nil {
			e = dev.Close()
		}
		if e == nil {
			a.mu.Lock()
			a.job.Incomplete = false
			a.mu.Unlock()
			a.event("Schreiben und Zurücklesen bestanden. Karte sicher auswerfen; für APKs/OBBs den Leser neu verbinden.")
		}
		a.finish(e)
	}()
}
func (a *guiApp) uploadImports(w http.ResponseWriter, r *http.Request, obb bool) {
	kind := "APKs vorbereiten"
	if obb {
		kind = "OBBs vorbereiten"
	}
	ctx, e := a.begin(kind, "Dateien werden hochgeladen")
	if e != nil {
		apiError(w, e)
		return
	}
	dir, e := os.MkdirTemp(a.temp, "imports-")
	if e != nil {
		a.finish(e)
		apiError(w, e)
		return
	}
	ok := false
	defer func() {
		if !ok {
			os.RemoveAll(dir)
		}
	}()
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)
	mr, e := r.MultipartReader()
	id := ""
	paths := []string{}
	names := []string{}
	if e == nil {
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				e = err
				break
			}
			if part.FormName() == "inbox_id" {
				b, _ := io.ReadAll(io.LimitReader(part, 256))
				id = string(b)
				continue
			}
			if part.FormName() != "files" || len(paths) >= 64 {
				e = errors.New("Ungültige Auswahl oder mehr als64Dateien.")
				break
			}
			name := filepath.Base(part.FileName())
			var dest string
			if obb {
				if obbRE.FindStringSubmatch(name) == nil {
					e = fmt.Errorf("Kein gültiger OBB-Dateiname: %s", name)
					break
				}
				dest = filepath.Join(dir, name)
			} else {
				if !strings.HasSuffix(strings.ToLower(name), ".apk") {
					e = errors.New("Bitte APK-Dateien auswählen.")
					break
				}
				dest = filepath.Join(dir, fmt.Sprintf("%03d.apk", len(paths)))
			}
			f, err := os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
			if err != nil {
				e = err
				break
			}
			n, err := io.Copy(f, io.LimitReader(part, 0x100000000))
			f.Close()
			if err != nil {
				e = err
				break
			}
			if n > 0xffffffff {
				e = errors.New("Einzeldatei ist größer als die FAT32-Grenze.")
				break
			}
			paths = append(paths, dest)
			names = append(names, name)
			part.Close()
		}
	}
	if e == nil && len(paths) == 0 {
		e = errors.New("Keine Dateien ausgewählt.")
	}
	var volume Inbox
	if e == nil {
		volume, e = a.findInbox(id)
	}
	if e != nil {
		a.finish(e)
		apiError(w, e)
		return
	}
	ok = true
	jsonReply(w, 202, map[string]bool{"ok": true})
	go func() {
		defer os.RemoveAll(dir)
		root, close, e := volume.Open()
		if e != nil {
			a.finish(e)
			return
		}
		defer close()
		if e = checkInbox(root); e != nil {
			a.finish(e)
			return
		}
		a.progress("Prüfen und auf die Karte übertragen", 0, int64(len(paths)))
		if obb {
			pkg, err := queueOBB(root, dir)
			if err == nil {
				a.event("OBBs bereit für nächsten Konsolenstart: " + pkg)
				a.progress("OBB-Daten geprüft und vorbereitet", int64(len(paths)), int64(len(paths)))
			}
			a.finish(err)
			return
		}
		failures := 0
		for i, p := range paths {
			if ctx.Err() != nil {
				a.finish(errors.New("Weitere Dateien werden nicht mehr übernommen; bereits vorbereitete Aufträge bleiben erhalten."))
				return
			}
			info, err := queueAPK(root, p)
			if err != nil {
				a.event("Nicht übernommen: " + names[i] + " – " + err.Error())
				failures++
			} else {
				a.event("APK bereit für nächsten Konsolenstart: " + info.Package)
			}
			a.progress("APKs prüfen und vorbereiten", int64(i+1), int64(len(paths)))
		}
		if failures > 0 {
			a.finish(fmt.Errorf("%d von %d APKs wurden abgelehnt. Einzelheiten im Protokoll.", failures, len(paths)))
		} else {
			a.finish(nil)
		}
	}()
}
func (a *guiApp) stop() {
	a.stopOnce.Do(func() {
		if a.server != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			a.server.Shutdown(ctx)
		}
		close(a.done)
	})
}
func (a *guiApp) start() (string, error) {
	ln, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		return "", e
	}
	a.base = "http://" + ln.Addr().String()
	a.server = &http.Server{Handler: http.HandlerFunc(a.serveHTTP), ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 64 * 1024}
	go func() { a.server.Serve(ln); a.stop() }()
	go func() {
		t := time.NewTicker(15 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-a.done:
				return
			case <-t.C:
				a.mu.Lock()
				idle := time.Since(a.lastSeen) > 15*time.Minute && a.job.Status != "running"
				a.mu.Unlock()
				if idle {
					a.stop()
					return
				}
			}
		}
	}()
	return a.base + "/#" + a.token, nil
}
func (a *guiApp) defaultImage() {
	if a.demo {
		return
	}
	var m Manifest
	json.Unmarshal(releaseJSON, &m)
	dir := filepath.Dir(a.exe)
	if b, e := os.ReadFile(filepath.Join(dir, "IMAGE-INFOS.json")); e == nil {
		if e = json.Unmarshal(b, &m); e != nil {
			a.event("Prüfdatei konnte nicht gelesen werden.")
			return
		}
	}
	if e := validManifest(m); e != nil {
		return
	}
	path := filepath.Join(dir, m.File)
	if _, e := os.Stat(path); e != nil {
		return
	}
	if _, e := a.begin("Image", "Vorhandenes Image prüfen"); e != nil {
		return
	}
	go func() { e := a.inspectImage(path, m.File, m, ""); a.finish(e) }()
}
func validateNotify(raw string) error {
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" || u.Port() == "" || u.Path != "/ready" || u.RawQuery != "" || u.User != nil {
		return errors.New("Ungültiger lokaler Startkanal")
	}
	return nil
}
