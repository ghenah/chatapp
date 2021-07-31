package httpserver

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/ghenah/chatapp/pkg/idatastore"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// userRegister
// @Summary Register a new user.
// @Description Register a new user by providing a password as well as a
// @Description unique username and email address.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RequestUserRegister true "Body must contain a username, an email, and a password"
// @Success 200 {object} ResponseSuccess
// @Failure 400
// @Failure 500
// @Router /auth/signup [post]
func userRegister(c echo.Context) (err error) {
	reqData := &RequestUserRegister{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = reqData.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), 12)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = ds.CreateUser(reqData.Username, reqData.Email, string(hashedPassword))
	if err != nil {
		if err == idatastore.ErrorDuplicateEntry {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return writeResponse(c, ResponseSuccess{Success: true})
}

// userAuthenticate
// @Summary Authenticate a user.
// @Description Check user's login credentials and provide an access token
// @Description if the registration was successful.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RequestUserAuthenticate true "Body must contain a username and a password"
// @Success 200 {object} ResponseAuthSuccess
// @Failure 400
// @Failure 500
// @Router /auth/signin [post]
func userAuthencticate(c echo.Context) (err error) {
	reqData := &RequestUserAuthenticate{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = reqData.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Verify the user login details
	userPassword, err := ds.GetUserPassword(reqData.Username)
	if err != nil {
		if err == idatastore.ErrorUserNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = bcrypt.CompareHashAndPassword(userPassword, []byte(reqData.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid login credentials")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// Generate the user "session"
	userInfo, err := ds.GetUser(reqData.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	accessToken, err := generateUserSession(userInfo.ID, userInfo.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// Generate a refresh token
	refreshToken, err := generateRefreshToken(userInfo.ID, userInfo.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.Path = "/refresh-token"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return writeResponse(c, ResponseAuthSuccess{
		User:        userInfo,
		AccessToken: accessToken,
	})
}

// getAuthenticatedUserInfo
// @Summary Get authenticated user info.
// @Description Get the up-to-date information on an authenticated user.
// @Tags user
// @Produce json
// @Success 200 {object} ResponseAuthSuccess
// @Failure 400
// @Failure 500
// @Router /api/v1/users/profile [get]
func getAuthenticatedUserInfo(c echo.Context) (err error) {
	reqData := &RequestUserAuthenticate{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get the user info from the access token.
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)

	// Generate the user "session"
	userInfo, err := ds.GetUser(claims.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	// Make sure the user ID in the access token corresponds to the ID of the
	// user fetched from the database (as users are fetched by the username).
	if claims.UserID != userInfo.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID and username don't match")
	}

	return writeResponse(c, ResponseAuthorizedUserInfo{
		User: userInfo,
	})
}

// userSearch
// @Summary List registered users.
// @Description Provides a list of registered users (usernames only).
// @Tags user
// @Produce json
// @Success 200 {object} ResponseUserSearch
// @Failure 500
// @Router /api/v1/users/search/ [get]
func userSearch(c echo.Context) (err error) {
	usersList, err := ds.GetAllUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return writeResponse(c, ResponseUserSearch{
		UsersList: usersList,
	})
}

// userFriendAdd
// @Summary Add a new friend to a user's friend list.
// @Description Adds a new friend to the list of friends of a user.
// @Tags user
// @Accept json
// @Param body body RequestAddUserToList true "Body must contain a user ID and a friend's ID."
// @Produce json
// @Success 200 {object} ResponseSuccess
// @Failure 500
// @Router /api/v1/users/friends [post]
func userFriendAdd(c echo.Context) (err error) {
	reqData := &RequestAddUserToList{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = reqData.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// A user cannot add themselves to the friends list.
	if reqData.UserID == reqData.FriendID {
		return writeResponse(c, ResponseSuccess{Success: true})
	}

	// Make sure the UserID belongs to the authenticated user (the owner of
	// the JWT)
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)
	if reqData.UserID != claims.UserID {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect user details")
	}

	err = ds.AddFriend(reqData.UserID, reqData.FriendID)
	if err == idatastore.ErrorUserInIgnoreList {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return writeResponse(c, ResponseSuccess{Success: true})
}

// userFriendRemove
// @Summary Remove a friend from the user's friends list.
// @Description Remove a friend to the list of friends of a user.
// @Tags user
// @Accept json
// @Param body body RequestAddUserToList true "Body must contain a user ID and a friend's ID."
// @Produce json
// @Success 200 {object} ResponseSuccess
// @Failure 500
// @Router /api/v1/users/friends [delete]
func userFriendRemove(c echo.Context) (err error) {
	reqData := &RequestAddUserToList{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = reqData.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Make sure the UserID belongs to the authenticated user (the owner of
	// the JWT)
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)
	if reqData.UserID != claims.UserID {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect user details")
	}

	err = ds.RemoveFriend(reqData.UserID, reqData.FriendID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return writeResponse(c, ResponseSuccess{Success: true})
}

// userIgnoredAdd
// @Summary Add a new user to a user's ignore list.
// @Description Adds a new user to the list of friends of a user. If the user
// @Description was in the friends list, they are removed from it.
// @Tags user
// @Accept json
// @Param body body RequestAddUserToList true "Body must contain a user ID and an ignored user's ID."
// @Produce json
// @Success 200 {object} ResponseSuccess
// @Failure 500
// @Router /api/v1/users/ignored [post]
func userIgnoredAdd(c echo.Context) (err error) {
	reqData := &RequestAddUserToList{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = reqData.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// A user cannot add themselves to the ignore list.
	if reqData.UserID == reqData.FriendID {
		return writeResponse(c, ResponseSuccess{Success: true})
	}

	// Make sure the UserID belongs to the authenticated user (the owner of
	// the JWT)
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)
	if reqData.UserID != claims.UserID {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect user details")
	}

	err = ds.AddIgnored(reqData.UserID, reqData.FriendID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return writeResponse(c, ResponseSuccess{Success: true})
}

// userIgnoredRemove
// @Summary Remove a user from the user's ignored list.
// @Description Remove a user from the ignored list of a user.
// @Tags user
// @Accept json
// @Param body body RequestAddUserToList true "Body must contain a user ID and a friend's ID."
// @Produce json
// @Success 200 {object} ResponseSuccess
// @Failure 500
// @Router /api/v1/users/ignored [delete]
func userIgnoredRemove(c echo.Context) (err error) {
	reqData := &RequestAddUserToList{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = reqData.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Make sure the UserID belongs to the authenticated user (the owner of
	// the JWT)
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)
	if reqData.UserID != claims.UserID {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect user details")
	}

	err = ds.RemoveIgnored(reqData.UserID, reqData.FriendID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return writeResponse(c, ResponseSuccess{Success: true})
}

// userUpdatePassword
// @Summary Update the user password.
// @Description Update the user password. Both the old and the new passwords
// @Description must the supplied.
// @Tags user
// @Accept json
// @Param body body RequestUserUpdatePassword true "Body must contain a user ID, the old password and the new password"
// @Produce json
// @Success 200 {object} ResponseSuccess
// @Failure 500
// @Router /api/v1/users/update/password [put]
func userUpdatePassword(c echo.Context) (err error) {
	reqData := &RequestUserUpdatePassword{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = reqData.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Make sure the UserID belongs to the authenticated user (the owner of
	// the JWT)
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)
	if (reqData.UserID != claims.UserID) || (reqData.Username != claims.Username) {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect user details")
	}

	// Validate the user password
	userPassword, err := ds.GetUserPassword(reqData.Username)
	if err == idatastore.ErrorUserNotFound {
		return echo.NewHTTPError(http.StatusBadRequest)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = bcrypt.CompareHashAndPassword(userPassword, []byte(reqData.OldPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid login credentials")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// Update the password
	newPasswordHashed, err := bcrypt.GenerateFromPassword([]byte(reqData.NewPassword), 12)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = ds.UpdateUserPassword(reqData.UserID, string(newPasswordHashed))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return writeResponse(c, ResponseSuccess{Success: true})
}

// userUpdateUsername
// @Summary Update the username.
// @Description Update the username. The old password must also be supplied.
// @Tags user
// @Accept json
// @Param body body RequestUserUpdateUsername true "Body must contain a user ID, the old password and the new username"
// @Produce json
// @Success 200 {object} ResponseSuccess
// @Failure 500
// @Router /api/v1/users/update/username [put]
func userUpdateUsername(c echo.Context) (err error) {
	reqData := &RequestUserUpdateUsername{}
	if err = c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = reqData.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Make sure the UserID belongs to the authenticated user (the owner of
	// the JWT)
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)
	if reqData.UserID != claims.UserID {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect user details")
	}

	// Validate the user password
	userPassword, err := ds.GetUserPassword(reqData.Username)
	if err == idatastore.ErrorUserNotFound {
		return echo.NewHTTPError(http.StatusBadRequest)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = bcrypt.CompareHashAndPassword(userPassword, []byte(reqData.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid login credentials")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// Update the username
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = ds.UpdateUsername(reqData.UserID, reqData.NewUsername)
	if err == idatastore.ErrorDuplicateEntry {
		return echo.NewHTTPError(http.StatusBadRequest, "username taken")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return writeResponse(c, ResponseSuccess{Success: true})
}

// chatRoomSearch
// @Summary List currently active chat rooms.
// @Description Provides a list of currently active chat rooms. The rooms that
// @Description have their owners present in the user's ignore list are
// @Descripiton not included.
// @Tags chat
// @Produce json
// @Success 200 {object} ResponseChatRoomSearch
// @Failure 500
// @Router /api/v1/chat/rooms/search/ [get]
func chatRoomSearch(c echo.Context) (err error) {
	// Take the user details out of the access token
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)

	// The user ID is required for the app to filter out ingored users' chats.
	chatRoomsList, err := ca.GetAllChatRooms(claims.UserID, claims.Username, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return writeResponse(c, ResponseChatRoomSearch{
		ChatRoomsList: chatRoomsList,
	})
}

// chatGetWSTicket
// @Summary Get a WebSocket authentication ticket.
// @Description Obtain a WebSocket authentication ticket, which must be
// @Description provided during the initial handshake.
// @Tags chat
// @Produce json
// @Success 200 {object} ResponseWSTicket
// @Failure 500
// @Router /api/v1/chat/ticket [get]
func chatGetWSTicket(c echo.Context) (err error) {
	// Take the user details out of the access token
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)

	ticket, err := generateUserSession(claims.UserID, claims.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return writeResponse(c, ResponseWSTicket{
		WSTicket: ticket,
	})
}

// refreshAccessToken
// @Summary Refresh the access token.
// @Description Refresh the access token using the refresh token from the
// @Description Http-Only cookie.
// @Tags refresh-token
// @Produce json
// @Success 200 {object} ResponseATRefreshSuccess
// @Failure 500
// @Router /refresh-token [get]
func refreshAccessToken(c echo.Context) (err error) {
	// Take the user details out of the access token
	// u := c.Get("user").(*jwt.Token)
	// claims := u.Claims.(*Claims)

	// Generate a new access token
	userInfo, err := ds.GetUser("Jack")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	accessToken, err := generateUserSession(userInfo.ID, userInfo.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	// Generate a new refresh token as well
	refreshToken, err := generateRefreshToken(userInfo.ID, userInfo.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.Path = "/refresh-token"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return writeResponse(c, ResponseATRefreshSuccess{
		AccessToken: accessToken,
	})
}

// chatConnectionInit performs the necessary chat initialization steps
// when a user upgrades their connection to WebSocket.
func chatConnectionInit(c echo.Context) (err error) {
	// Take the user details out of the access token
	ticket := c.Get("ticket").(*jwt.Token)
	claims := ticket.Claims.(*Claims)

	err = serveConnection(c.Response(), c.Request(), ca.InMsgQueue, claims.UserID, claims.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

// writeResponse writes the response in the format specified in the
// Accept header; the default format is "application/json".
// Supported formats:
// - JSON
func writeResponse(c echo.Context, data interface{}) error {
	switch c.Request().Header.Get("Accept") {
	case "application/json":
		return c.JSON(http.StatusOK, data)
	default:
		return c.JSON(http.StatusOK, data)
	}
}
