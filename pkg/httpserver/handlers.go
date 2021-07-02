package httpserver

import (
	"fmt"
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

	// Validate the user login details
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

	return writeResponse(c, ResponseAuthSuccess{
		User:        userInfo,
		AccessToken: accessToken,
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

	// Make sure the UserID belongs to the authenticated user (the owner of
	// the JWT)
	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*Claims)
	if (reqData.UserID != claims.UserID) || (reqData.Username != claims.Username) {
		fmt.Println("User ID:", reqData.UserID)
		fmt.Println("Username:", reqData.Username)
		fmt.Println("JWT User ID:", claims.UserID)
		fmt.Println("JWT Username:", claims.Username)
		fmt.Println("Claims:", claims)
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

	return writeResponse(c, ResponseSuccess{Success: true})
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
