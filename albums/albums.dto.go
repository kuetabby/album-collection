package albums

import shared "playlist/shared"

type Album struct {
    ID     int64  `json:"id"`
    Title  string `json:"title"`
    Artist string `json:"artist"`
    Price  float32 `json:"price"`
	User_Id int `json:"user_id"`
}

type GetAllAlbumsResponse struct {
	shared.BaseResponse
	Content []Album `json:"content"`
	TotalRecord int `json:"totalRecord"`
}

type GetAlbumResponse struct {
	shared.BaseResponse
	Content Album `json:"content"`
}