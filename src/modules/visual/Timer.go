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
		TimerPTR, err := memory.GlobalManager.ReadPointer(memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x090821C8, 0x230, 0x20, 0x40, 0xB50})
		TimerPTR += 0x80

		if err != nil {
			time.Sleep(time.Millisecond * 20)
			modFunction()
			return
		}
		modules.RegisterModule(&modules.Module{
			Author:      "spot",
			Version:     "v1.0.0",
			Description: "MEOW MEOW TIMER",
			OnDisable: func(module *modules.Module) {
				if err := memory.GlobalManager.Restore(TimerPTR); err != nil {
					helpers.LogF("failed to restore timer value: %+v\n", err)
				} else {
					helpers.LogF("Restored timer!\n")
				}
			},
			OnEnable: func(module *modules.Module) {
				if err := memory.GlobalManager.Write(TimerPTR, float32(100)); err != nil {
					helpers.LogF("failed to set timer value: %+v\n", err)
				}
			},
			Enabled: false,
		})
	}
	modules.AfterStartup(modFunction)
}
