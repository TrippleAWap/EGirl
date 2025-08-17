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

		CordXptr, err := memory.GlobalManager.ReadPointer(memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x090A54C8, 0x10, 0x20, 0xC8})
		helpers.Check(err)
		CordXptr += 0x598

		CordYptr, err := memory.GlobalManager.ReadPointer(memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x0090A54C8, 0x18, 0x20, 0xC8, 0x220})
		helpers.Check(err)
		CordYptr += 4

		CordZptr, err := memory.GlobalManager.ReadPointer(memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x0905C8B8, 0x10, 0x128, 0x0, 0x110})
		helpers.Check(err)
		CordZptr += 0x20
		if err != nil {
			time.Sleep(time.Millisecond * 20)
			modFunction()
			return
		}
		modules.RegisterModule(&modules.Module{
			Author:      "TrippleAWap",
			Version:     "v1.0.0",
			Description: "meow meow fullbright",
			OnTick: func(module *modules.Module) {
				var y float32
				var x float32
				var z float32

				helpers.Check(memory.GlobalManager.Read(CordXptr, &x))
				helpers.Check(memory.GlobalManager.Read(CordYptr, &y))
				helpers.Check(memory.GlobalManager.Read(CordZptr, &z))

				helpers.LogF("%v %v %v\n", x, y, z)
			},
			Enabled: true,
		})
	}
	modules.AfterStartup(modFunction)
}
