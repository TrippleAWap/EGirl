package helpers

import (
	"unsafe"
)

// typedef struct tagCURSORINFO {
// DWORD   cbSize;
// DWORD   flags;
// HCURSOR hCursor;
// POINT   ptScreenPos;
// } CURSORINFO, *PCURSORINFO, *LPCURSORINFO;

type CURSORINFO struct {
	cbSize      uint32
	flags       uint32
	hCursor     uintptr
	ptScreenPos struct {
		x int32
		y int32
	}
}

var (
	GetCursorInfo = user32.NewProc("GetCursorInfo")
)

const (
	CURSOR_SHOWING = 0x00000001
)

func IsMouseVisible() bool {
	ci := CURSORINFO{}
	ci.cbSize = uint32(unsafe.Sizeof(ci))
	ret, _, _ := GetCursorInfo.Call(uintptr(unsafe.Pointer(&ci)))
	if ret == 0 {
		return false
	}
	return ci.flags&CURSOR_SHOWING != 0
}
