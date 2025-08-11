package modules

import (
	"EGirl/helpers"
	"EGirl/memory"
	"math"
)

func init() {
	helpers.LogF("%+v\n", memory.InterfaceToBytes(float32(math.MaxFloat32)))
}
