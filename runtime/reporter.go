package runtime

import "fmt"

type Reporter interface {
	Error(format string, a ...any)
	Print(format string)
}

type StdoutReporter struct {
}

func (r *StdoutReporter) Error(format string, a ...any) {
	fmt.Printf(format, a...)
}

func (r *StdoutReporter) Print(s string) {
	fmt.Print(s)
}
