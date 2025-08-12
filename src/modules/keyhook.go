package modules

import (
	"EGirl/helpers"
	"golang.org/x/sys/windows"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procSetWindowsHookEx = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx   = user32.NewProc("CallNextHookEx")
	procGetMessage       = user32.NewProc("GetMessageW")
	keyboardHook         uintptr
	targetWindowHWND     windows.HWND
	onTargetWindow       = true
)

func UpdateOnTargetWindow() {
	for {
		if targetWindowHWND == 0 {
			targetWindowHWNDi, err := helpers.FindWindow("Minecraft")
			if err != nil {
				panic(err)
			}
			targetWindowHWND = targetWindowHWNDi
		}
		foreground := windows.GetForegroundWindow()
		onTargetWindow = targetWindowHWND == foreground
		time.Sleep(time.Millisecond * 5)
	}
}

const (
	WH_KEYBOARD    = 2
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 256
)

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

func keyboardCallback(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if onTargetWindow && nCode >= 0 && wParam == WM_KEYDOWN {
		kb := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		r := rune(kb.VkCode)
		if printableASCII(r) {
			keyChar := string(r)
			helpers.LogF("%s\n", keyChar)
		}
	}
	ret, _, _ := procCallNextHookEx.Call(keyboardHook, uintptr(nCode), wParam, lParam)
	return ret
}

// TODO: migrate to raw input, maybe... our current solution works!
// https://learn.microsoft.com/en-us/windows/win32/inputdev/raw-input
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getrawinputdevicelist
func initializeKeyHook() {
	//go UpdateOnTargetWindow()
	keyboardHook, _, _ = procSetWindowsHookEx.Call(
		WH_KEYBOARD_LL,
		syscall.NewCallback(keyboardCallback),
		0,
		0,
	)
	var msg struct {
		HWND   uintptr
		UINT   uint32
		WPARAM int16
		LPARAM int64
		DWORD  uint32
		POINT  struct{ X, Y int64 }
	}
	for {
		_, _, _ = procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
	}
}
func printableASCII(r rune) bool {
	return r >= 32 && r <= 126
}
