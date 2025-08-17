package visual

import (
	"EGirl/helpers"
	"EGirl/memory"
	"EGirl/modules"
)

func init() {
	var modFunction func()
	var xPtr, yPtr, zPtr uintptr
	var fixPtrs func() error
	fixPtrs = func() error {
		return memory.GlobalManager.ReadPtrConfigs([]memory.PtrConfig{
			{memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x090A54C8, 0x10, 0x20, 0xC8}, 0x598, &xPtr},
			{memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x0090A54C8, 0x18, 0x20, 0xC8, 0x220}, 0x4, &yPtr},
			{memory.GlobalManager.BaseModule.ModBaseAddr, []uintptr{0x0905C8B8, 0x10, 0x128, 0x0, 0x110}, 0x20, &zPtr},
		})
	}
	modFunction = func() {
		defer helpers.PanicDisplay()

		modules.RegisterModule(&modules.Module{
			Author:      "spot",
			Version:     "v1.0.0",
			Description: "meow coords",
			OnTick: func(module *modules.Module) {
				if err := fixPtrs(); err != nil {
					return
				}
				var x, y, z float32

				helpers.Check(memory.GlobalManager.Read(xPtr, &x))
				helpers.Check(memory.GlobalManager.Read(yPtr, &y))
				helpers.Check(memory.GlobalManager.Read(zPtr, &z))

				helpers.LogF("%v %v %v\n", x, y, z)
			},
			Enabled: true,
		})
	}
	modules.AfterStartup(modFunction)
}
