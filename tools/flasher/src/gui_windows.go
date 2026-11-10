//go:build windows

package main

import (
	"errors"
	"os/exec"
	"syscall"
	"unsafe"
)

func guiIsAdmin() bool { v, _, _ := shell.NewProc("IsUserAnAdmin").Call(); return v != 0 }
func guiElevate(exe, notify, nonce string) error {
	if guiIsAdmin() {
		cmd := exec.Command(exe, "--gui-worker", "--notify", notify, "--nonce", nonce)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if e := cmd.Start(); e != nil {
			return e
		}
		go cmd.Wait()
		return nil
	}
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exe)
	args, _ := syscall.UTF16PtrFromString(`--gui-worker --notify "` + notify + `" --nonce "` + nonce + `"`)
	r, _, _ := shell.NewProc("ShellExecuteW").Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(args)), 0, 0)
	if r <= 32 {
		return errors.New("Administrator-Anfrage wurde nicht bestätigt.")
	}
	return nil
}
func guiOpenBrowser(u string) error {
	cmd := exec.Command("explorer.exe", u)
	if e := cmd.Start(); e != nil {
		return e
	}
	go cmd.Wait()
	return nil
}
func guiShowError(e error) {
	text, _ := syscall.UTF16PtrFromString(e.Error())
	title, _ := syscall.UTF16PtrFromString("GA36 Flasher")
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), 0x10)
}
