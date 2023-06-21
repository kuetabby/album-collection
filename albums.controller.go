package main

import (
	"fmt"
	"net/http"
	_albums "playlist/albums"
	_shared "playlist/shared"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllAlbums(c *gin.Context) {
	albums, errAlbums := _albums.GetAlbums(Conn, Ctx)

	if(errAlbums != nil) {
		errMsg := errAlbums.Error()
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[_albums.Album]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return
	} else {	
		fmt.Printf("length of the albums is %d", len(albums))
		c.JSON(http.StatusOK, _albums.GetAllAlbumsResponse{
			Content: albums,
			TotalRecord: len(albums),
			BaseResponse: _shared.BaseResponse{
				Status: "SUCCESS",
				Error: false,
			},
		})
		return
	}
}

func GetOneAlbumById(c *gin.Context) {
	params := c.Param("id")
	id, errId := strconv.Atoi(params)

	if(errId != nil){
		errMsg := fmt.Sprintf("%v isn't a valid number", params)
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return
	} else {
		album,errAlbum := _albums.GetAlbumsById(id, Conn, Ctx)

		if(errAlbum != nil) {
			errMsg := errAlbum.Error()
			c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
				BaseResponse: _shared.BaseResponse{
					Status: "ERROR",
					Message: &errMsg,
					Error: true,
				},
			})
			return
		} else {
			c.JSON(http.StatusOK, _albums.GetAlbumResponse{
				Content: album,
				BaseResponse: _shared.BaseResponse{
					Status: "SUCCESS",
					Error: false,
				},
			})
			return
		}
	}
}

func GetOneAlbumByName(c *gin.Context) {
	name := c.Param("name")

	if(len(name) == 0) {
		errMsg := fmt.Sprintln("name can't be empty")
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return
	} else {
		albums,errAlbums := _albums.GetAlbumsByArtist(name, Conn, Ctx)

		if(errAlbums != nil) {
			errMsg := errAlbums.Error()
			c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
				BaseResponse: _shared.BaseResponse{
					Status: "ERROR",
					Message: &errMsg,
					Error: true,
				},
			})
			return
		} else {
			c.JSON(http.StatusOK, _albums.GetAllAlbumsResponse{
				Content: albums,
				TotalRecord: len(albums),
				BaseResponse: _shared.BaseResponse{
					Status: "SUCCESS",
					Error: false,
				},
			})
			return
		}
	}
}

func PostCreateAlbum(c *gin.Context) {
	user_id := c.GetInt64("user_id")

	var newAlbum _albums.Album

	if err:= c.BindJSON(&newAlbum); err != nil {
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
		// fmt.Println(newAlbum)
	 album, errAlbum := _albums.CreateAlbum(newAlbum, int(user_id), Conn, Ctx)
	 
	 if(errAlbum != nil){
		errMsg := errAlbum.Error()
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return
	 } else {
		c.JSON(http.StatusOK, _shared.SingleResponse[_albums.Album]{
			Content: &album,
			BaseResponse: _shared.BaseResponse{
				Status: "SUCCESS",
				Error: false,
			},
		})
		return
	 }
	}
}

func PostUpdateAlbum(c *gin.Context){
	params := c.Param("id")
	id, errId := strconv.Atoi(params)

	if errId != nil {
		// errMsg := errId.Error()
		errMsg := fmt.Sprintf("%v isn't a valid number", params)
		c.JSON(http.StatusNotFound, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return 
	} 

	var updateAlbum _albums.Album
	
	if isVerified := verifyAlbum(id, c); !isVerified  {
	 if errBind := c.BindJSON(&updateAlbum); errBind != nil {
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
			album,errAlbum := _albums.UpdateAlbum(_albums.Album{
				ID: int64(id),
				Title: updateAlbum.Title,
				Artist: updateAlbum.Artist,
				Price: updateAlbum.Price,	
			}, Conn, Ctx)

			if errAlbum != nil {
				errMsg := errAlbum.Error()
				c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
					BaseResponse: _shared.BaseResponse{
						Status: "ERROR",
						Message: &errMsg,
						Error: true,
					},
				})
				return
			} 

			c.JSON(http.StatusOK, _shared.SingleResponse[_albums.Album]{
					Content: album,
					BaseResponse: _shared.BaseResponse{
						Status: "SUCCESS",
						Error: false,
					},
				})
			return
		}
	}
}

func DeleteAlbum(c *gin.Context){
	params := c.Param("id")
	id, errId := strconv.Atoi(params)

	if errId != nil {
		// errMsg := errId.Error()
		errMsg := fmt.Sprintf("%v isn't a valid number", params)
		c.JSON(http.StatusNotFound, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return 
	} 

	if isVerified := verifyAlbum(id, c); !isVerified  {
		_,errAlbum := _albums.RemoveAlbum(id, Conn, Ctx)

		if errAlbum != nil {
			errDelete := errAlbum.Error()
			c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
				BaseResponse: _shared.BaseResponse{
					Status: "ERROR",
					Message: &errDelete,
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
}

func verifyAlbum(id int, c *gin.Context)(bool) {
	if getUser, errGetUser := decodeLogin(c); errGetUser != nil {
		errMsg := errGetUser.Error()
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status:  "ERROR",
				Message: &errMsg,
				Error:   true,
			},
		})
		return true
	} else if getAlbum,errGetAlbum := _albums.GetAlbumsById(id, Conn, Ctx); errGetAlbum != nil {
		errMsg := errGetAlbum.Error()
		c.JSON(http.StatusBadRequest, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return true
	} else if getAlbum.User_Id != int(getUser.ID) {
		errMsg := "You are not authorized to update this album"
		c.JSON(http.StatusUnauthorized, _shared.SingleResponse[any]{
			BaseResponse: _shared.BaseResponse{
				Status: "ERROR",
				Message: &errMsg,
				Error: true,
			},
		})
		return true
	}  else {
		return false
	}
}