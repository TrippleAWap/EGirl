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
