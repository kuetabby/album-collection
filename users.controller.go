package main

import (
	"fmt"
	"net/http"
	_shared "playlist/shared"
	_users "playlist/users"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context){
	var newUser _users.User

	if errBind := c.BindJSON(&newUser); errBind != nil {
		errMsg := errBind.Error()
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return
	} else {
		user, errUser  := _users.CreateUser(newUser, Conn, Ctx)

		if errUser != nil {
			errUserMsg := errUser.Error()
			c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
				BaseResponse: _shared.BaseResponse{
					Status: "ERROR",
					Message: &errUserMsg,
					Error: true,
				},
			})
			return
		}else {
			c.JSON(http.StatusOK, _shared.SingleResponse[_users.User]{
				Content: &user,
				BaseResponse: _shared.BaseResponse{
					Status: "SUCCESS",
					Error: false,
				},
			})
			return
		}
	}
}

func GetAllUsers(c *gin.Context) {
	users, errUsers := _users.GetAllUsers(Conn, Ctx)

	role := c.GetString("user_role")

	fmt.Printf("the role is %v\n", role)

	if(errUsers != nil) {
		errMsg := errUsers.Error()
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[_users.User]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return
	} else {	
		c.JSON(http.StatusOK, _users.GetAllUsersResponse[_users.CombinedUserWithAlbum]{
			Content: users,
			TotalRecord: len(users),
			BaseResponse: _shared.BaseResponse{
				Status: "SUCCESS",
				Error: false,
			},
		})
		return
	}
}

func DeleteUser(c *gin.Context) {
	address := c.Param("address")
	role := c.GetString("user_role")

	_,errUser := _users.RemoveUser(address, Conn, Ctx)

	if role != "admin" {
		errMsg := "only admin can delete user"
		c.JSON(http.StatusUnauthorized, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
	} else if errUser != nil {
		errMsg := errUser.Error()
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return
	} else {
		c.JSON(http.StatusOK, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "SUCCESS",
				Error: false,
			},
		})
	}
}

func GetOneUserByAddress(c *gin.Context) {
	address := c.Param("address")
	// id, errId := strconv.Atoi(params)

	// if(errId != nil){
	// 	errMsg := fmt.Sprintf("%v isn't a valid number", params)
	// 	c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
	// 		BaseResponse: _shared.BaseResponse{
	// 			Status: "ERROR",
	// 			Message: &errMsg,
	// 			Error: true,
	// 		},
	// 	})
	// 	return
	// } else {
		user,errUser := _users.GetUserByAddress(address, Conn, Ctx)

		if(errUser != nil) {
			errMsg := errUser.Error()
			c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
				BaseResponse: _shared.BaseResponse{
					Status: "ERROR",
					Message: &errMsg,
					Error: true,
				},
			})
			return
		} else {
			c.JSON(http.StatusOK, _users.GetUserResponse{
				Content: user,
				BaseResponse: _shared.BaseResponse{
					Status: "SUCCESS",
					Error: false,
				},
			})
			return
		}
	// }
}