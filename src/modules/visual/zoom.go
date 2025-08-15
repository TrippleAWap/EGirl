package visual

import (
	"EGirl/helpers"
	"EGirl/memory"
	"EGirl/modules"
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
				var fovOA float32

				if err := memory.GlobalManager.Original(FOVPtr, &fovOA); err != nil {
					helpers.LogF("failed to restore FOV value: %+v\n", err)
				} else {
					helpers.LogF("Restored FOV!\n")
				}
			},
			OnTick: func(module *modules.Module) {
				if err := memory.GlobalManager.Write(FOVPtr, float32(10)); err != nil {
					helpers.LogF("failed to set FOV value: %+v\n", err)
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
