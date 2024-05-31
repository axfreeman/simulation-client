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

// Simple authorization page. Allows the user to choose who to play as.
// Requests the server to authorize the selected user.
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
	// pick up the chosen user
	username := ctx.Params.ByName("username")
	if username == "" {
		utils.Trace(utils.Red, "The router did not pick up a valid player name")
		return
	}
	user, ok := models.Users[username]
	if !ok {
		utils.Trace(utils.Red, "The router did not pick up a valid player name")
		return
	}

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
	utils.Trace(utils.Gray, fmt.Sprintf("User %s will play\n", username))
	http.SetCookie(ctx.Writer, &http.Cookie{Name: "user", Value: username, Path: "/"})
	ctx.Request.URL.Path = `/user/dashboard`
	Router.HandleContext(ctx)
	ctx.Abort()
}

// func DisplayPlayerChoice(ctx *gin.Context) {
// 	ctx.HTML(http.StatusOK, "choose-player.html", gin.H{
// 		"Title":      "Choose player",
// 		"adminusers": models.AdminUserList,
// 		"users":      models.Users,
// 	})
// 	ctx.Abort()
// }

// Quit playing as the current user and start playing as another
// Doesn't mean leave the game.
//
//	 ... you can quit but you can never leave...
//
//		 Locally, set 'IsLoggedIn'
//		 Tell the server
//		 If the user cannot be found just return (error will already have been signalled)
//		 If the server complains, display an error.
func Quit(ctx *gin.Context) {
	utils.Trace(utils.Gray, "Quit was requested\n")
	userobject, ok := ctx.Get("userobject")
	utils.Trace(utils.Gray, fmt.Sprintf("Asked for a user object and OK status was %v\n", ok))
	utils.Trace(utils.Gray, fmt.Sprintf("The user object was %v\n", userobject))
	if !ok {
		DisplayErrorScreen(ctx, "Couldn't find a player. This is a programme error.\n")
	}
	user := userobject.(*models.User)
	utils.Trace(utils.Gray, fmt.Sprintf("The user name was %v\n", user.UserName))

	// Guest can never quit.
	if user.UserName != `guest` {
		// All others have to unlock.
		user.IsLocked = false
		_, err := api.ServerRequest(user.ApiKey, `admin/unlock/`+user.UserName) //TODO server should delete this user's simulations
		if err != nil {
			utils.DisplayError(ctx, fmt.Sprintf("User %s could not quit because the server objected.", user.UserName))
			ctx.Abort()
			return
		}

		http.SetCookie(ctx.Writer, &http.Cookie{Name: "user", Value: user.UserName, Path: "/", MaxAge: 0})
		utils.Trace(utils.Gray, fmt.Sprintf("%s has quit\n", user.UserName))
	}

	ctx.HTML(http.StatusOK, "choose-player.html", gin.H{
		"Title": "Choose player",
		"users": models.Users,
	})
	ctx.Abort()
	// log.Fatal("Diagnostic halt")
}
