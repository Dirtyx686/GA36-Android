//go:build linux

package main

import (
	"errors"
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

func guiIsAdmin() bool { return os.Geteuid() == 0 }
func guiElevate(exe, notify, nonce string) error {
	if guiIsAdmin() {
		cmd := exec.Command(exe, "--gui-worker", "--notify", notify, "--nonce", nonce)
		if e := cmd.Start(); e != nil {
			return e
		}
		go cmd.Wait()
		return nil
	}
	p, e := exec.LookPath("pkexec")
	if e != nil {
		return errors.New("Für den grafischen Start wird pkexec benötigt. Alternativ den Flasher im Terminal mit sudo starten.")
	}
	cmd := exec.Command(p, exe, "--gui-worker", "--notify", notify, "--nonce", nonce)
	if e = cmd.Start(); e != nil {
		return e
	}
	go cmd.Wait()
	return nil
}
func guiOpenBrowser(u string) error {
	if os.Geteuid() != 0 {
		cmd := exec.Command("xdg-open", u)
		if e := cmd.Start(); e != nil {
			return e
		}
		go cmd.Wait()
		return nil
	}
	id := os.Getenv("SUDO_UID")
	if id == "" {
		id = os.Getenv("PKEXEC_UID")
	}
	n, e := strconv.Atoi(id)
	if e != nil || n <= 0 {
		return errors.New("Bitte den Flasher als normaler Benutzer starten. Die Anmeldung erfolgt anschließend grafisch.")
	}
	usr, e := user.LookupId(id)
	if e != nil {
		return e
	}
	cmd := exec.Command("runuser", "-u", usr.Username, "--", "xdg-open", u)
	if e = cmd.Start(); e != nil {
		return e
	}
	go cmd.Wait()
	return nil
}
func guiShowError(e error) {
	if z, err := exec.LookPath("zenity"); err == nil {
		exec.Command(z, "--error", "--title=GA36 Flasher", "--text="+e.Error()).Run()
	}
}
