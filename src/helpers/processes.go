package helpers

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

func GetProcessID(name string) (int, error) {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return -1, err
	}
	defer windows.CloseHandle(handle)

	process := windows.ProcessEntry32{Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{}))}
	for windows.Process32Next(handle, &process) == nil {
		processName := windows.UTF16ToString(process.ExeFile[:])
		if processName == name {
			return int(process.ProcessID), nil
		}
	}

	return -1, nil
}

var (
	targetWindowName = ""
	findWindow       = user32.NewProc("FindWindowW")
)

func FindWindow(windowName string) (windows.HWND, error) {
	v := windows.StringToUTF16Ptr(windowName)
	hwnd, _, err := findWindow.Call(0, uintptr((unsafe.Pointer)(v)))
	if err != nil && err.Error() != "The operation completed successfully." {
		return 0, err
	}
	return windows.HWND(hwnd), nil
}

func GetThreadIDs(pid int) ([]uint32, error) {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, uint32(pid))
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)
	var r []uint32
	thread := windows.ThreadEntry32{Size: uint32(unsafe.Sizeof(windows.ThreadEntry32{}))}
	for windows.Thread32Next(handle, &thread) == nil {
		if thread.OwnerProcessID != uint32(pid) {
			continue
		}
		r = append(r, thread.ThreadID)
	}

	return r, nil
}
