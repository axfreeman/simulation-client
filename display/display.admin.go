// display.admin.go
// handlers for actions specific to the admin

package display

import (
	"capfront/models"
	"capfront/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Display the admin dashboard
func AdminDashboard(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
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
// TODO not done yet
func AdminReset(ctx *gin.Context) {
}
