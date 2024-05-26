//utils.report.go

// we are still experimenting with reportage.
/// this is a kind of generic place to put the experikments

package utils

import (
	"fmt"
)

// switch detailed tracing on or off
//
//	if TraceLevel is true, produce trace diagnostics
var TraceLevel bool = true

// Provide status information on simulation and authorization progress.
//
// TODO ONLY for development - replace by Logger calls at some point
func Trace(startColour string, message string) {
	if !TraceLevel {
		return
	}
	fmt.Print(startColour + message + Reset)
}
