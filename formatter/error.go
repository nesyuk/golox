package printer

import "fmt"

func ReportError(line int, where string, message string) string {
	return fmt.Sprintf("[line %d] Error%v: %v", line, where, message)
}
