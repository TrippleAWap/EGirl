package memory

import (
	"EGirl/helpers"
	"fmt"
	"reflect"
	"time"
)

func (m *Manager) SmoothWrite(address uintptr, start, step any, steps int, duration time.Duration) error {
	startV := reflect.ValueOf(start)
	stepV := reflect.ValueOf(step)
	helpers.LogF("%v+%v*%v\n", start, step, steps)
	for i := 0; i < steps; i++ {
		var outV reflect.Value

		switch startV.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			base := startV.Int()
			delta := stepV.Int()
			v := base + delta*int64(i)
			outV = reflect.ValueOf(v).Convert(startV.Type())

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			base := startV.Uint()
			delta := stepV.Uint()
			v := base + delta*uint64(i)
			outV = reflect.ValueOf(v).Convert(startV.Type())

		case reflect.Float32, reflect.Float64:
			base := startV.Float()
			delta := stepV.Float()
			v := base + delta*float64(i)
			outV = reflect.ValueOf(v).Convert(startV.Type())

		default:
			return fmt.Errorf("SmoothWrite: unsupported kind %s", startV.Kind())
		}
		if err := m.Write(address, outV.Interface()); err != nil {
			return err
		}
		time.Sleep(duration / time.Duration(steps))
	}

	return nil
}
