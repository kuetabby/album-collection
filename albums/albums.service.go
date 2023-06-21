package albums

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateAlbum(album Album, userId int, conn *pgxpool.Pool, ctx context.Context)(Album, error) {
	var createdAlbum Album

	if( len(album.Title) == 0 || len(album.Artist) == 0 || album.Price == 0) {
		isErrValidate, errValidate := validateAlbum(album)

		if(isErrValidate && errValidate != nil){
			return createdAlbum, errValidate
		}
	}

	result,err := conn.Exec(ctx, `INSERT INTO "Album" (title, artist, price, user_id) VALUES ($1, $2, $3, $4)`, album.Title, album.Artist, album.Price, userId)
	if err != nil {
        return createdAlbum, fmt.Errorf("addAlbum: %v", err)
    }

	rowsAffected := result.RowsAffected()

	fmt.Printf("rows affected %d\n", rowsAffected)

	if (rowsAffected > 0) {
		row := conn.QueryRow(ctx, `SELECT id,title,artist,price,user_id FROM "Album" WHERE title = $1`, album.Title)
		if errRow := row.Scan(&createdAlbum.ID, &createdAlbum.Title, &createdAlbum.Artist, &createdAlbum.Price, &createdAlbum.User_Id); errRow != nil {
			if errRow == pgx.ErrNoRows {
				return createdAlbum, fmt.Errorf("get album %v: no such album", album.Title)
			}
			return createdAlbum, fmt.Errorf("failed to get the album: %v", err)
		}
	}

	return createdAlbum, nil
}

func GetAlbums(conn *pgxpool.Pool, ctx context.Context)([]Album, error) {
	var albums []Album

	rows,err := conn.Query(ctx, `SELECT id, title, artist,price FROM "Album" ORDER BY id ASC`)
	if(err != nil){
		return albums, fmt.Errorf("albums : %v", err)
	}

	defer rows.Close()

	hasRows := rows.Next()

	for hasRows {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return albums, fmt.Errorf("unable to get next album: %v", err)
		}

		fmt.Printf("album %v", album)
		hasRows = rows.Next() // Update the hasRows variable

		albums = append(albums, album)
	}

	// Handle the case when there are no rows
	if len(albums) == 0 {
		return []Album{}, nil
	}

	if rows.Err(); err != nil {
		return albums, fmt.Errorf(`albums rows error: %v`, err)
	}

	return albums, nil
}

func GetAlbumsById(id int, conn *pgxpool.Pool, ctx context.Context)(Album, error) {
	var album Album

	row := conn.QueryRow(ctx, `SELECT * FROM "Album" where id = $1`, id)
	if err := row.Scan(&album.ID, &album.User_Id, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == pgx.ErrNoRows {
			return album, fmt.Errorf("albumsById %d: no such album", id)
		}
		return album, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return album, nil
}

func GetAlbumsByArtist(name string, conn *pgxpool.Pool, ctx context.Context)([]Album, error){
	var albums []Album

	rows,err := conn.Query(ctx, `SELECT * FROM "Album" WHERE artist = $1 ORDER BY id ASC`, name)
    if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }

	defer rows.Close()

	for rows.Next(){
		var album Album

		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func UpdateAlbum(album Album, conn *pgxpool.Pool, ctx context.Context)(*Album, error) {
	var updatedAlbum Album

	if( len(album.Title) == 0 || len(album.Artist) == 0 || album.Price == 0) {
		isErrValidate, errValidate := validateAlbum(album)
		if(isErrValidate && errValidate != nil){
			return &updatedAlbum, errValidate
		}
	}

	result,err := conn.Exec(ctx, `UPDATE "Album" SET title = $1, artist = $2, price = $3 WHERE id = $4`, album.Title, album.Artist, album.Price, album.ID)
	if(err != nil) {
		return nil, fmt.Errorf("failed to update album: %v", err)
	}

	rowsAffected := result.RowsAffected()

	if(rowsAffected > 0){
		err = conn.QueryRow(ctx, `SELECT id, title, artist, price FROM "Album" WHERE title = $1`, album.Title).Scan(&updatedAlbum.ID, &updatedAlbum.Title, &updatedAlbum.Artist, &updatedAlbum.Price)
		if(err != nil){
			return nil, fmt.Errorf("failed to retrieve updated album: %v", err)
		}
	}

	return &updatedAlbum, nil
}

func RemoveAlbum(id int, conn *pgxpool.Pool, ctx context.Context)(string, error) {
	if(id == 0){
		errMsg := fmt.Sprintf("%d can't be empty", id)
		return "", fmt.Errorf(errMsg)
	}

	result,err := conn.Exec(ctx, `DELETE FROM "Album" WHERE id = $1`, id)
	if(err != nil) {
		return "", fmt.Errorf("deleteAlbum %d: %v", id, err)
	}

	rowsAffected := result.RowsAffected()

	if(rowsAffected == 0){
		return "", fmt.Errorf("failed to delete album %d: no rows affected", rowsAffected)
	}

	return "SUCCESS", nil
}

func validateAlbum(album Album) (bool, error) {
	var emptyProperties []string

	switch 0 {
	case len(album.Title):
		emptyProperties = append(emptyProperties, "Title")
	case len(album.Artist):
		emptyProperties = append(emptyProperties, "Artist")
	case int(album.Price):
		emptyProperties = append(emptyProperties, "Price")
	default:
	}

	if(len(emptyProperties) > 0){
		errMsg := fmt.Sprintf("%s properties cannot be empty", strings.Join(emptyProperties, ", "))

		return true, fmt.Errorf(errMsg)
	}

	return false, nil
}
