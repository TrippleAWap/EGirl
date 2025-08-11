package main

import (
	"EGirl/helpers"
	_ "EGirl/modules"
	"github.com/bi-zone/go-fileversion"
	"golang.org/x/sys/windows"
	"os"
	"slices"
	"time"
)

var (
	NAME         = ""
	GAME_VERSION = ""
)

const (
	VERSION = "0.1.1"
	BRANCH  = "development"
)

var (
	SupportedVersions = []string{
		"1.21.100.6",
	}
	Pointers = map[string]map[string]func(*MemoryManager) (uintptr, error){
		"1.21.100.6": {
			"Brightness": func(memManager *MemoryManager) (uintptr, error) {
				v, err := memManager.ReadPointer(memManager.baseModule.ModBaseAddr, []uintptr{0x9083788, 0x188, 0x40, 0x1B0})
				v += 0x18
				return v, err
			},
		},
	}
)

func init() {
	helpers.GetProjectRoot()
	NAME = helpers.GetProjectRoot()
}

func main() {
	defer PanicDisplay()
	targetPID := -1
	var err error
	helpers.LogF("Starting %s v%s (%s)...\n", NAME, VERSION, BRANCH)
	helpers.LogF("Getting process ID for Minecraft.Windows.exe...\n")
	for targetPID < 0 {
		targetPID, err = getProcessID("Minecraft.Windows.exe")
		if err != nil {
			helpers.LogF("We've encountered an error while getting the process id. | %+v\n", err.Error())
			os.Exit(1)
		}
		time.Sleep(time.Millisecond * 100)
	}
	helpers.LogF("Minecraft.Windows.exe found with PID: %d\n", targetPID)
	baseModule, err := getBaseModule(targetPID)
	if err != nil {
		helpers.LogF("We've encountered an error while getting the base module. | %+v\n", err.Error())
		os.Exit(1)
	}
	f, err := fileversion.New(windows.UTF16ToString(baseModule.ExePath[:]))
	if err != nil {
		helpers.LogF("We've encountered an error while reading the file properties! | %+v\n", err.Error())
		os.Exit(1)
	}
	GAME_VERSION = f.FileVersion()
	supportedVersion := slices.Contains(SupportedVersions, GAME_VERSION)
	supportedVersionStr := "\x1b[35mUNSUPPORTED\x1b[0m"
	if supportedVersion {
		supportedVersionStr = "\x1b[32mSUPPORTED\x1b[0m"
	}
	helpers.LogF("Version: %s %s\n", GAME_VERSION, supportedVersionStr)
	if !supportedVersion {
		return
	}
	helpers.LogF("Base Module: 0x%X\n", baseModule.ModBaseAddr)
	memManager := MemoryManager{}
	defer memManager.Cleanup()
	if err := memManager.OpenProcess(targetPID); err != nil {
		helpers.LogF("We've encountered an error while opening the process handle. | %+v\n", err.Error())
		os.Exit(1)
	}

	if err := memManager.LoadProcessMemory(); err != nil {
		helpers.LogF("We've encountered an error while loading the process memory. | %+v\n", err.Error())
		os.Exit(1)
	}
	helpers.LogF("Process memory loaded successfully.\n")

	// read brightness;
	func() {
		//for ptrName, ptrFunc := range Pointers[GAME_VERSION] {
		//	helpers.LogF("Resolving '%s' for '%s'", ptrName, GAME_VERSION)
		//	ptr, err := ptrFunc(&memManager)
		//	if err != nil {
		//		helpers.LogF("failed to get '%s' pointer: %+v\n", ptrName, err)
		//		return
		//	}
		//}
		brightnessPtr, err := Pointers[GAME_VERSION]["Brightness"](&memManager)
		if err != nil {
			helpers.LogF("failed to get '%s' pointer: %+v\n", "brightnessPtr", err)
			return
		}
		helpers.LogF("found brightnessPtr, 0x%X\n", brightnessPtr)
		var brightness float32
		if err := memManager.Read(brightnessPtr, &brightness); err != nil {
			helpers.LogF("failed to get brightness value: %+v\n", err)
			return
		}
		helpers.LogF("found brightness value, %f\n", brightness)
		helpers.LogF("toggling full-bright!\n")
		if err := memManager.Write(brightnessPtr, &brightness); err != nil {
			helpers.LogF("failed to set brightness value: %+v\n", err)
		}
		if err := memManager.Read(brightnessPtr, &brightness); err != nil {
			helpers.LogF("failed to get brightness value: %+v\n", err)
			return
		}
		helpers.LogF("updated brightness value, %f\n", brightness)
	}()
}
