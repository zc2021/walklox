package reporters

import (
	"fmt"
	"os"
)

func PrintErr(line int, msg string) {
	report(line, "", msg)
}

func report(line int, where, msg string) {
	fmt.Fprintf(os.Stderr,
		"[line %d] Error%s: %", line, where, msg)
}
