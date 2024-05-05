// display.admin.go
// handlers for actions specific to the admin

package display

import (
	"capfront/api"
	"capfront/models"
	"capfront/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Display the admin dashboard
func AdminDashboard(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Data to delete a simulation ")
		return
	}

	if username != "admin" {
		utils.DisplayError(ctx, "Only the administrator can see the admin dashboard")
		return
	}
	ctx.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
		"Title":    "Admin Dashboard",
		"users":    models.Users,
		"username": username,
	})
}

// Resets the main database
// Only available to admin
func AdminReset(ctx *gin.Context) {
	username := utils.GUESTUSER

	if username != "admin" {
		log.Output(1, fmt.Sprintf("User %s tried to reset the database", username))

	}

	_, jsonErr := api.ServerRequest(username, "reset the database", `action/reset`)
	if jsonErr != nil {
		log.Output(1, "Reset failed")
	} else {
		log.Output(1, "COMPLETE RESET by admin")
	}

	AdminDashboard(ctx)
}
