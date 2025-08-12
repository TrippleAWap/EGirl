package modules

import (
	"EGirl/helpers"
	"strings"
	"time"
)

type GameState struct{}

type ModArrayListRendererDelegate func(*Module) string
type ModIsRelevantDelegate func(*Module) bool
type ModVoidDelegate func(*Module)

type ModuleOptionType byte

const (
	Toggle ModuleOptionType = iota
	FloatSlider
	IntSlider
)

type ToggleI struct {
	State bool
}

type ModuleOption struct {
	Type  ModuleOptionType
	Value interface{}
}

type Module struct {
	Name        string
	Author      string
	Version     string
	Category    string
	Description string

	KeyBind rune

	Options map[string]ModuleOption

	ArrayListRenderer ModArrayListRendererDelegate
	IsRelevant        ModIsRelevantDelegate

	OnTick    ModVoidDelegate
	OnEnable  ModVoidDelegate
	OnDisable ModVoidDelegate
}

var (
	modules           []*Module
	AfterStartupFuncs []func()
)

func RegisterHandles() {
	go func() {
		defer helpers.PanicDisplay()
		for {
			for _, m := range modules {
				if m.OnTick == nil {
					continue
				}
				m.OnTick(m)
			}
			time.Sleep(time.Millisecond * 20)
		}
	}()

	go initializeKeyHook()
}

func RegisterModule(module *Module) {
	if module == nil {
		panic("module is nil")
	}
	frame := helpers.GetLastCallerFrame(1)
	parts := strings.Split(strings.TrimSuffix(frame.File, ".go"), "/")
	moduleName := FormatName(parts[len(parts)-1])
	categoryName := FormatName(parts[len(parts)-2])
	module.Category = categoryName
	module.Name = moduleName
	helpers.LogF(
		"\n\tRegistering '%s'\n\t\tAuthor '%s'\n\t\tDescription '%s'\n\t\tCategory '%s'\n",
		module.Name,
		module.Author,
		module.Description,
		module.Category,
	)
	modules = append(modules, module)

}

func FormatName(name string) string {
	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), " ", "")
}

func AfterStartup(f func()) {
	AfterStartupFuncs = append(AfterStartupFuncs, f)
}
