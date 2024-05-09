//logging.report.go

// we are still experimenting with reportage.
/// this is a kind of generic place to put the experikments

package logging

import (
	"capfront/colour"
	"fmt"
)

// switch detailed tracing on or off
//
//	if TraceLevel is true, produce trace diagnostics
var TraceLevel bool = true

// function to provide details to help trace where things happened
func Trace(startColour string, message string) {
	if !TraceLevel {
		return
	}
	fmt.Print(startColour + message + colour.Reset)
}
