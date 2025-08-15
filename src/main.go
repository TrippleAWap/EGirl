package main

import (
	"EGirl/helpers"
	"EGirl/memory"
	"EGirl/modules"
	_ "EGirl/modules"
	_ "EGirl/modules/visual"
	"github.com/bi-zone/go-fileversion"
	"golang.org/x/sys/windows"
	"os"
	"os/signal"
	"runtime"
	"slices"
	"sync"
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
)

func init() {
	helpers.GetProjectRoot()
	NAME = helpers.GetProjectRoot()
}

func main() {
	go debug()
	defer helpers.PanicDisplay()
	targetPID := -1
	var err error
	helpers.LogF("Starting %s v%s (%s)...\n", NAME, VERSION, BRANCH)
	helpers.LogF("Getting process ID for Minecraft.Windows.exe...\n")
	for targetPID < 0 {
		targetPID, err = helpers.GetProcessID("Minecraft.Windows.exe")
		if err != nil {
			helpers.LogF("We've encountered an error while getting the process id. | %+v\n", err.Error())
			os.Exit(1)
		}
		time.Sleep(time.Millisecond * 100)
	}
	helpers.LogF("Minecraft.Windows.exe found with PID: %d\n", targetPID)
	baseModule, err := memory.GetBaseModule(targetPID)
	if err != nil {
		helpers.LogF("We've encountered an error while getting the base module. | %+v\n", err.Error())
		main()
		return
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
	defer memory.GlobalManager.Cleanup()
	if err := memory.GlobalManager.OpenProcess(targetPID); err != nil {
		helpers.LogF("We've encountered an error while opening the process handle. | %+v\n", err.Error())
		os.Exit(1)
	}

	if err := memory.GlobalManager.LoadProcessMemory(); err != nil {
		helpers.LogF("We've encountered an error while loading the process memory. | %+v\n", err.Error())
		os.Exit(1)
	}
	helpers.LogF("Process memory loaded successfully.\n")
	for _, f := range modules.AfterStartupFuncs {
		go f()
	}
	modules.RegisterHandles()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		_ = <-c
	}()

	wg.Wait()
}

func debug() {
	for {
		helpers.LogF("%v\n", runtime.NumGoroutine())
		helpers.LogF("%v\n", runtime.NumCPU())

		time.Sleep(time.Second)
	}
}
