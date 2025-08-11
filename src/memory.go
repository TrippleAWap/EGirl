package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

func getBaseModule(pid int) (*windows.ModuleEntry32, error) {
	baseModule := windows.ModuleEntry32{
		Size: uint32(unsafe.Sizeof(windows.ModuleEntry32{})),
	}
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, uint32(pid))
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	if err := windows.Module32First(handle, &baseModule); err != nil {
		return nil, err
	}
	return &baseModule, nil
}

type MemoryManager struct {
	hProcess      windows.Handle
	pId           int
	processMemory []byte
	baseModule    *windows.ModuleEntry32
	memoryPatches map[uintptr][]byte
}

func (m *MemoryManager) OpenProcess(pid int) error {
	m.pId = pid
	handle, err := windows.OpenProcess(windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_VM_OPERATION, false, uint32(pid))
	if err != nil {
		return err
	}
	m.hProcess = handle
	m.baseModule, err = getBaseModule(pid)
	if err != nil {
		return err
	}
	return nil
}

func (m *MemoryManager) Read(address uintptr, output interface{}) error {
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Ptr {
		return errors.New("output must be a pointer to struct")
	}
	elem := rv.Elem()
	if !elem.CanSet() {
		return errors.New("output is not settable")
	}

	var size uintptr
	k := elem.Kind()
	switch k {
	case reflect.Slice:
		size = uintptr(elem.Len())
	case reflect.Array:
		size = uintptr(elem.Len())
	default:
		size = elem.Type().Size()
	}

	//LogF("0%X, %d, %s\n", address, size, k.String())

	buffer := make([]byte, size)
	var read uintptr
	if err := windows.ReadProcessMemory(m.hProcess, address, &buffer[0], size, &read); err != nil {
		return err
	}
	if read != size {
		return windows.ERROR_PARTIAL_COPY
	}

	v := reflect.ValueOf(buffer).Slice(0, int(size))
	switch k {
	case reflect.Uintptr:
		if err := setUintptrFromBytes(v, elem); err != nil {
			return err
		}
	case reflect.Float64:
		if err := setFloat64FromBytes(v, elem); err != nil {
			return err
		}
	case reflect.Float32:
		if err := setFloat32FromBytes(v, elem); err != nil {
			return err
		}
	default:
		elem.Set(v)
	}
	return nil
}

func setUintptrFromBytes(v reflect.Value, elem reflect.Value) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Uint8 {
		return errors.New("v must be a []byte")
	}
	if elem.Kind() != reflect.Uintptr {
		return errors.New("elem must be uintptr")
	}

	b := v.Bytes()
	var x uint64
	switch len(b) {
	case 1, 2, 4, 8:
		buf := make([]byte, 8)
		copy(buf, b)
		x = binary.LittleEndian.Uint64(buf)
	default:
		return errors.New("invalid length for uintptr conversion")
	}

	ptrVal := uintptr(x)
	elem.Set(reflect.ValueOf(ptrVal))
	return nil
}

func setFloat32FromBytes(v reflect.Value, elem reflect.Value) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Uint8 {
		return errors.New("v must be a []byte")
	}
	if elem.Kind() != reflect.Float32 {
		return errors.New("elem must be float32")
	}

	data := v.Bytes()
	expected := int(unsafe.Sizeof(float32(0)))
	if len(data) != expected {
		return errors.New(
			"invalid length for float32 conversion: expected " +
				strconv.Itoa(expected) +
				", got " + strconv.Itoa(len(data)))
	}

	bits := binary.LittleEndian.Uint32(data)
	f := math.Float32frombits(bits)
	elem.SetFloat(float64(f))
	return nil
}

func setFloat64FromBytes(v reflect.Value, elem reflect.Value) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Uint8 {
		return errors.New("v must be a []byte")
	}
	if elem.Kind() != reflect.Float64 {
		return errors.New("elem must be float64")
	}

	data := v.Bytes()
	expected := int(unsafe.Sizeof(float64(0)))
	if len(data) != expected {
		return errors.New(
			"invalid length for float64 conversion: expected " +
				strconv.Itoa(expected) +
				", got " + strconv.Itoa(len(data)))
	}

	bits := binary.LittleEndian.Uint64(data)
	f := math.Float64frombits(bits)
	elem.SetFloat(f)
	return nil
}

func (m *MemoryManager) Write(address uintptr, data interface{}) error {
	val := reflect.ValueOf(data)
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		val = val.Elem()
	}
	size := val.Type().Size()

	var oldProtection uint32
	if err := windows.VirtualProtectEx(m.hProcess, address, size, windows.PAGE_EXECUTE_READWRITE, &oldProtection); err != nil {
		return err
	}
	valueP := val
	var write uintptr
	if err := windows.WriteProcessMemory(m.hProcess, address, (*byte)(unsafe.Pointer(&valueP)), size, &write); err != nil {
		return err
	}
	if err := windows.VirtualProtectEx(m.hProcess, address, size, oldProtection, &oldProtection); err != nil {
		return err
	}

	if write != size {
		return windows.ERROR_PARTIAL_COPY
	}
	if _, ok := m.memoryPatches[address]; !ok {
		output := make([]byte, size)
		if err := m.Read(address, &output); err != nil {
			return err
		}
		m.memoryPatches[address] = output
	}
	return nil
}

func (m *MemoryManager) LoadProcessMemory() error {
	baseModule, err := getBaseModule(m.pId)
	if err != nil {
		return err
	}
	var processMemoryV = make([]byte, baseModule.ModBaseSize)
	err = m.Read(baseModule.ModBaseAddr, &processMemoryV)
	if err != nil {
		return err
	}
	m.processMemory = processMemoryV
	m.memoryPatches = make(map[uintptr][]byte)
	return nil
}

// Scan finds a signature in the process memory and returns the address of the first match.
func (m *MemoryManager) Scan(pattern []byte, mask string) uintptr {
	var result uintptr = 0

scanForSegment:
	for i := 0; i <= len(m.processMemory)-len(pattern); i++ {
		for j := 0; j < len(pattern); j++ {
			if (pattern[j] == '?' && mask[j] == '?') || m.processMemory[i+j] == pattern[j] {
				if j == len(pattern)-1 {
					result = uintptr(i) + m.baseModule.ModBaseAddr
					break scanForSegment
				}
			} else {
				break
			}
		}
	}
	return result
}

func (m *MemoryManager) ReadPointer(address uintptr, pointers []uintptr) (uintptr, error) {
	var r uintptr
	addr := address
	for _, pointer := range pointers {
		err := m.Read(addr+pointer, &r)
		if err != nil {
			return 0, err
		}
		//LogF("[0x%X+0x%X] => 0x%X\n", addr, pointer, r)
		addr = r
	}
	return r, nil
}

func (m *MemoryManager) Restore(address uintptr) error {
	if _, ok := m.memoryPatches[address]; !ok {
		return fmt.Errorf("address not patched: 0x%x", address)
	}
	if err := m.Write(address, m.memoryPatches[address]); err != nil {
		return err
	}
	delete(m.memoryPatches, address)
	return nil
}

func (m *MemoryManager) Cleanup() error {
	if m.hProcess != 0 {
		for address := range m.memoryPatches {
			if err := m.Restore(address); err != nil {
				fmt.Println(err)
			}
		}
		if err := windows.CloseHandle(m.hProcess); err != nil {
			return err
		}
		m.hProcess = 0
	}
	return nil
}
