package modules

import (
	"EGirl/helpers"
	"golang.org/x/sys/windows"
	"strings"
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
	onTargetWindow       = true
	getForegroundWindow  = user32.NewProc("GetForegroundWindow")
	getWindowTextW       = user32.NewProc("GetWindowTextW")
	getWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
)

func GetForegroundWindowTitle() (string, error) {
	hwnd, _, _ := getForegroundWindow.Call()
	if hwnd == 0 {
		return "", nil
	}

	textLen, _, _ := getWindowTextLengthW.Call(hwnd)
	if textLen == 0 {
		return "", nil
	}

	buf := make([]uint16, textLen+1)
	getWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return windows.UTF16ToString(buf), nil
}

func UpdateOnTargetWindow() {
	for {
		title, err := GetForegroundWindowTitle()
		if err != nil {
			panic(err)
		}
		title = strings.ToLower(strings.TrimSpace(title))
		onTargetWindow = strings.Contains(title, "minecraft")
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

var KeyMap = make([]bool, 0xFE)

func keyboardCallback(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if onTargetWindow && nCode >= 0 && !helpers.IsMouseVisible() {
		kb := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		r := rune(kb.VkCode)
		shouldToggle := false
		if wParam == WM_KEYUP && KeyMap[kb.VkCode] != false {
			KeyMap[kb.VkCode] = false
		}
		if wParam == WM_KEYDOWN && KeyMap[kb.VkCode] != true {
			KeyMap[kb.VkCode] = true
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
	go UpdateOnTargetWindow()
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
