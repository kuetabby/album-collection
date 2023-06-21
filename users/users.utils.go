package users

import (
	"fmt"
	_albums "playlist/albums"
	"strings"
)

func validateUser(user User) (bool, error) {
	var emptyProperties []string

	switch 0 {
	case len(user.Address):
		emptyProperties = append(emptyProperties, "Address")
	default:
	}

	if(len(emptyProperties) > 0){
		errMsg := fmt.Sprintf("%s properties cannot be empty", strings.Join(emptyProperties, ", "))

		return true, fmt.Errorf(errMsg)
	}

	return false, nil
}


// Function to check if an album object has default values
func isDefaultAlbum(album _albums.Album) bool {
    return album.ID == 0 && album.Title == "" && album.Artist == "" && album.Price == 0 && album.User_Id == 0
}