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
	// pick up the cookie information that was set by SynchWithServer middleware
	username := ctx.Params.ByName("username")
	if username == "" {
		utils.Trace(utils.Red, " ERROR: The router did not pick up a valid player name")
	}
	user := models.Users[username]

	// lock this user at the server.
	// It's just possible someone else gets in first, so abort if this doesn't work.
	_, err := api.ServerRequest(user.ApiKey, `admin/lock/`+username)
	if err != nil {
		utils.DisplayError(ctx, fmt.Sprintf("Could not play as %s. Maybe somebody else got in first. Try again and tell me if the error persists", username))
		ctx.Abort()
		return
	}

	// lock this user at the client
	user.IsLocked = true

	// Set cookie with no MaxAge and no Expiry
	// NOTE it is claimed this will be deleted when browser closes, but it isn't.
	// see, eg https://stackoverflow.com/questions/10617954/chrome-doesnt-delete-session-cookies/10772420#10772420
	utils.Trace(utils.Gray, fmt.Sprintf("user %s will play\n", username))
	http.SetCookie(ctx.Writer, &http.Cookie{Name: "user", Value: username, Path: "/"})
	ctx.Redirect(http.StatusPermanentRedirect, `/user/dashboard`)
	// ctx.Request.URL.Path = `/`
	// Router.HandleContext(ctx)
}

func Lock(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "choose-player.html", gin.H{
		"Title":      "Choose player",
		"adminusers": models.AdminUserList,
		"users":      models.Users,
	})
}

// Quit playing as the current user.
//
//	Locally, set 'IsLoggedIn'
//	Tell the server
//	If the user cannot be found just return (error will already have been signalled)
//	If the server complains, display an error.
func Quit(ctx *gin.Context) {
	utils.Trace(utils.Gray, "Quit was requested\n")
	userobject, ok := ctx.Get("userobject")
	if !ok {
		return
	}
	user := userobject.(*models.User)
	user.IsLocked = false
	_, err := api.ServerRequest(user.ApiKey, `admin/unlock/`+user.UserName) //TODO server should delete this user's simulations
	if err != nil {
		utils.DisplayError(ctx, fmt.Sprintf("User %s could not quit because the server objected.", user.UserName))
		ctx.Abort()
		return
	}

	// Delete any cookie stil hanging around
	http.SetCookie(ctx.Writer, &http.Cookie{Name: "user", Value: user.UserName, Path: "/", MaxAge: 0})
	utils.Trace(utils.Gray, fmt.Sprintf("%s has quit\n", user.UserName))
	ctx.Request.URL.Path = `/`
	Router.HandleContext(ctx)
}
