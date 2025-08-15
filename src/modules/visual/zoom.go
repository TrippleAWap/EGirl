package visual

import (
	"EGirl/helpers"
	"EGirl/memory"
	"EGirl/modules"
	"math"
	"time"
)

func init() {
	var modFunction func()
	modFunction = func() {
		defer helpers.PanicDisplay()
		FOVPtr, err := memory.GlobalManager.ReadPointer(memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x9024e70, 0x18, 0x20, 0x1A8})
		FOVPtr += 0x18

		if err != nil {
			time.Sleep(time.Millisecond * 20)
			modFunction()
			return
		}
		modules.RegisterModule(&modules.Module{
			Author:      "spot",
			Version:     "v1.0.0",
			Description: "meow zoom",
			KeyBind:     'C',
			OnDisable: func(module *modules.Module) {
				var fovOriginal float32
				var fovCurrent float32

				if err := memory.GlobalManager.Original(FOVPtr, &fovOriginal); err != nil {
					helpers.LogF("failed to restore FOV value: %+v\n", err)
					return
				}
				if err := memory.GlobalManager.Read(FOVPtr, &fovCurrent); err != nil {
					helpers.LogF("failed to read FOV value: %+v\n", err)
					return
				}
				diff := float32(math.Abs(float64(fovOriginal - fovCurrent)))
				helpers.LogF("%v, %v-%v\n", diff, fovOriginal, fovCurrent)
				smoothingMS := 200
				if err := memory.GlobalManager.SmoothWrite(FOVPtr, fovCurrent, diff/float32(smoothingMS), smoothingMS, time.Millisecond); err != nil {
					helpers.LogF("failed to smooth FOV value: %+v\n", err)
				}
				if err := memory.GlobalManager.Restore(FOVPtr); err != nil {
					helpers.LogF("failed to restore FOV value: %+v\n", err)
				}
			},
			OnTick: func(module *modules.Module) {
				targetFov := float32(10)
				var fovCurrent float32
				if err := memory.GlobalManager.Read(FOVPtr, &fovCurrent); err != nil {
					helpers.LogF("failed to read FOV value: %+v\n", err)
					return
				}
				if fovCurrent != targetFov {
					diff := targetFov - fovCurrent
					smoothingMS := 200
					if err := memory.GlobalManager.SmoothWrite(FOVPtr, fovCurrent, diff/float32(smoothingMS), smoothingMS, time.Millisecond); err != nil {
						helpers.LogF("failed to move FOV value: %+v\n", err)
					}
				}
				if !modules.KeyMap[module.KeyBind] {
					module.SetActive(false)
				}
			},
			Enabled: false,
		})
	}
	modules.AfterStartup(modFunction)
}
