package helpers

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	user32     = syscall.NewLazyDLL("user32.dll")
	MessageBox = user32.NewProc("MessageBoxW")
)

func PanicDisplay() {
	if r := recover(); r != nil {

		stackBuf := make([]byte, 2048)
		stackSize := runtime.Stack(stackBuf, true)

		message := fmt.Sprintf("Panic Occured: %v\n\n%s", r, stackBuf[:stackSize])
		title := "Whoops, something went wrong!"

		titlePtr, _ := syscall.UTF16PtrFromString(title)
		messagePtr, _ := syscall.UTF16PtrFromString(message)

		_, _, err := MessageBox.Call(
			0,
			uintptr(unsafe.Pointer(messagePtr)),
			uintptr(unsafe.Pointer(titlePtr)),
			0x10,
		)
		if err != nil && err.Error() != "The operation completed successfully." {
			panic(err)
		}
		os.Exit(1)
	}
}
