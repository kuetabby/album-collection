package main

import (
	"net/http"

	_shared "playlist/shared"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	var loginData LoginRequest

	if err := c.BindJSON(&loginData); err != nil {
		errMsg := err.Error()
        c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
        return
    } else {
		accessToken, errToken := authentication(loginData)

		if errToken != nil {
			errMsg := errToken.Error()
			c.JSON(http.StatusUnauthorized, _shared.SingleResponse[any]{
				BaseResponse: _shared.BaseResponse{
					Status: "ERROR",
					Message: &errMsg,
					Error: true,
				},
			})
			return
		} else {
			c.JSON(http.StatusOK, _shared.SingleResponse[string]{
				Content: &accessToken,
				BaseResponse: _shared.BaseResponse{
					Status: "SUCCESS",
					Error: false,
				},
			})
			return
		}
	}
}