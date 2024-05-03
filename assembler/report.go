package assembler

import "fmt"

// ReportDiscard discards assembler errors.
func ReportDiscard(_ AssembleError) {}

// ReportPrint prints assembler errors to standard out.
func ReportPrint(err AssembleError) {
	_, _ = fmt.Printf("Assembler Error on line %d: %s", err.LineNo, err.Message)
}
