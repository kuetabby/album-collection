package users

import (
	albums "playlist/albums"
	shared "playlist/shared"
	"time"
)

type User struct {
    ID     int64 `json:"id"`
	Address string `json:"address"`
	Username string `json:"username"`
	Role string `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type CombinedUserWithAlbum struct {
	User
	Albums []albums.Album `json:"albums"`
}

type GetAllUsersResponse[T any] struct {
	shared.BaseResponse
	Content []T `json:"content"`
	TotalRecord int `json:"totalRecord"`
}

type GetUserResponse struct {
	shared.BaseResponse
	Content CombinedUserWithAlbum `json:"content"`
}