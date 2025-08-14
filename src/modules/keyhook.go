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
	WM_KEYUP       = 257
)

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

var keystates = make([]bool, 0xFE)

func keyboardCallback(nCode int, wParam uintptr, lParam uintptr) uintptr {
	go UpdateOnTargetWindow()
	if onTargetWindow && nCode >= 0 {
		kb := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		r := rune(kb.VkCode)
		shouldToggle := false
		if wParam == WM_KEYUP && keystates[kb.VkCode] != true {
			keystates[kb.VkCode] = true
		}
		if wParam == WM_KEYDOWN && keystates[kb.VkCode] != false {
			keystates[kb.VkCode] = false
			shouldToggle = true
		}
		if shouldToggle {
			for _, m := range modules {
				if m.IsRelevant() && m.KeyBind == r {
					m.SetActive(!m.Enabled)
				}
			}
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
