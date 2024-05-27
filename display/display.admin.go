// display.admin.go
// handlers for actions specific to the admin

package display

import (
	"capfront/api"
	"capfront/models"
	"capfront/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Display the admin dashboard
func AdminDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
		"Title": "Admin Dashboard",
		"users": models.AdminUserList,
	})
}

// Resets the main database
// Only available to admin
// TODO not done yet
func AdminReset(ctx *gin.Context) {
}

// Authorization function.
// Requests the server to authorize user to play.
// The client should already know this user's ApiKey.
// This is not a token system. It simply averts conflicts by ensuring
// that only one client can play as a given user at the same time.
//
//	Asks the server to lock the user
//
//	If the lock fails respond with an error and display on the client.
//
//	If the lock succeeds, set a cookie and register the lock on the User object.
func SelectUser(ctx *gin.Context) {
	u := ctx.Param("username")

	utils.Trace(utils.Yellow, fmt.Sprintf("user %s will play\n", u))
	// lock this user
	_, err := api.ServerRequest(models.Users[u].ApiKey, `admin/lock/`+u)
	if err != nil {
		utils.DisplayError(ctx, fmt.Sprintf("Could not play as user %s. It's just possible somebody else got in first", u))
		ctx.Abort()
		return
	}
	models.Users[u].IsLocked = true
	http.SetCookie(ctx.Writer, &http.Cookie{Name: "user", Value: u, Path: "/"})
	// TODO a more sensible redirect?
	ctx.Request.URL.Path = `/`
	Router.HandleContext(ctx)
}

func Lock(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "choose-player.html", gin.H{
		"Title": "Choose player",
		"users": models.AdminUserList,
	})
}

// Quit playing as the current user.
//
//	Locally, set 'IsLoggedIn'
//	Tell the server
//	If the user cannot be found just return (error will already have been signalled)
//	If the server complains, display an error.
func Quit(ctx *gin.Context) {
	utils.Trace(utils.Purple, "Quit was requested\n")
	userobject, ok := ctx.Get("userobject")
	if !ok {
		return
	}
	user := userobject.(*models.User)
	user.IsLocked = false
	_, err := api.ServerRequest(user.ApiKey, `admin/unlock/`+user.UserName)
	if err != nil {
		utils.DisplayError(ctx, fmt.Sprintf("User %s could not quit because the server objected.", user.UserName))
		ctx.Abort()
		return
	}

	utils.Trace(utils.Purple, fmt.Sprintf(" %s has quit\n", user.UserName))
	// TODO a more sensible redirect?
	ctx.Request.URL.Path = `/`
	Router.HandleContext(ctx)
}
