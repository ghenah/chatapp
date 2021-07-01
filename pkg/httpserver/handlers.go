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
