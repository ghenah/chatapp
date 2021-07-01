package httpserver

import (
	"net/http"

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
	accessToken, err := generateUserSession(userInfo.ID)
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
