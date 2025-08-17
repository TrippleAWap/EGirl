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
		TimerPTR, err := memory.GlobalManager.ReadPointer(memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x907BF28, 0x130, 0x38, 0x70, 0x418, 0x2E8})
		if err != nil {
			time.Sleep(time.Millisecond * 20)
			modFunction()
			return
		}
		modules.RegisterModule(&modules.Module{
			Author:      "spot & TrippleAWap",
			Version:     "v1.0.0",
			KeyBind:     'L',
			Description: "This module is basically useless as of the 'server-authoritative-movement' release.",
			OnDisable: func(module *modules.Module) {
				if err := memory.GlobalManager.Restore(TimerPTR); err != nil {
					helpers.LogF("failed to restore timer value: %+v\n", err)
				} else {
					helpers.LogF("Restored timer!\n")
				}
			},
			OnEnable: func(module *modules.Module) {
				var timerF float32
				if err := memory.GlobalManager.Read(TimerPTR, &timerF); err != nil {
					helpers.LogF("failed to read timer value: %+v\n", err)
				} else {
					helpers.LogF("Timer value: %f\n", timerF)
				}
				if err := memory.GlobalManager.Write(TimerPTR, float32(12)); err != nil {
					helpers.LogF("failed to set timer value: %+v\n", err)
				}
			},
			Enabled: false,
		})
	}
	modules.AfterStartup(modFunction)
}
