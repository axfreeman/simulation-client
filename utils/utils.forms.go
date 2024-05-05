// utils.forms.go
package utils

import (
	"capfront/colour"
	"capfront/logging"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Function to handle errors that need to be displayed to the user.
// Wrapped here to save space and also because error handling may
// be changed later.

func DisplayError(ctx *gin.Context, message string) {
	log.Output(1, message)
	ctx.HTML(http.StatusBadRequest, "errors.html", gin.H{
		"message": message,
	})
}

// Display the login form with a message and an offer to register.
func DisplayLogin(ctx *gin.Context, message string) {
	if pc, file, no, ok := runtime.Caller(1); ok {
		logging.Trace(colour.Cyan, fmt.Sprintf("Login form invoked by %s#%d with pointer %d", file, no, pc))
	}
	ctx.HTML(
		http.StatusOK,
		"login.html",
		gin.H{"message": message})
}

// Shortcut for DisplayLogin with a generic message
func InitialLogin(ctx *gin.Context) {
	DisplayLogin(ctx, "Please log in")
}
