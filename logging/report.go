//logging.report.go

// we are still experimenting with reportage.
/// this is a kind of generic place to put the experikments

package logging

import (
	"capfront/colour"
	"capfront/models"
	"fmt"
)

// creates a trace that logs events on the standard output
// and also creates a user-readable record called 'Trace'
// that provides information about the progress of the simulation

//TODO distinguish between what is logged to output and what is shown to the user

func Report(level int, simulation models.Simulation, message string) {
	fmt.Println(message) //placeholder
}

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
