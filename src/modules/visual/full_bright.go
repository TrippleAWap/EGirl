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
		brightnessPtr, err := memory.GlobalManager.ReadPointer(memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x9083788, 0x188, 0x40, 0x1B0})
		brightnessPtr += 0x18
		if err != nil {
			time.Sleep(time.Millisecond * 20)
			modFunction()
			return
		}
		modules.RegisterModule(&modules.Module{
			Author:      "TrippleAWap",
			Version:     "v1.0.0",
			Description: "meow meow fullbright",
			KeyBind:     'K',
			OnTick: func(module *modules.Module) {
				if err := memory.GlobalManager.Write(brightnessPtr, float32(10)); err != nil {
					helpers.LogF("failed to set brightness value: %+v\n", err)
				}
			},
		})
	}
	modules.AfterStartup(modFunction)
}
