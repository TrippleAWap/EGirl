package helpers

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

func LogF(format string, a ...any) {
	frame := GetLastCallerFrame(1)
	if frame != nil {
		format = FormatCaller(*frame) + " " + format
	}
	_, _ = fmt.Fprintf(os.Stdout, format, a...)
}
func GetLastCallerFrame(skipFrames int) *runtime.Frame {
	pc, _, _, ok := runtime.Caller(1 + skipFrames)
	if !ok {
		return nil
	}
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()
	return &frame
}
func FormatCaller(frame runtime.Frame) string {
	ret := frame.Function
	firstPart := strings.IndexRune(ret, '.') + 1
	ret = ret[firstPart:]
	ret = RemoveAll(ret, '(', ')', '*')
	ret = strings.ReplaceAll(ret, ".", "] [")
	ret = "[" + ret + "]"
	projectRoot := GetProjectRoot()
	relPath, _ := filepath.Rel(projectRoot, frame.File)
	ret += fmt.Sprintf(" [%s:L%d]", relPath, frame.Line)
	return ret
}

func RemoveAll(str string, args ...rune) string {
	var chars = []rune(str)
	var newChars []rune
	for _, char := range chars {
		if slices.Contains(args, char) {
			continue
		}
		newChars = append(newChars, char)
	}
	return string(newChars)
}

var rootCache string

func GetProjectRoot() string {
	if rootCache != "" {
		return rootCache
	}
	_, file, _, _ := runtime.Caller(1)
	dir := path.Dir(file)
	if strings.HasSuffix(file, "main.go") {
		rootCache = dir
	}
	return dir
}
