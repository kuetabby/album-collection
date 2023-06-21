package users

import (
	"context"
	"fmt"

	_albums "playlist/albums"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(user User, conn *pgxpool.Pool, ctx context.Context)(User, error) {
	var createdUser User

	if( len(user.Address) == 0) {
		isErrValidate, errValidate := validateUser(user)

		if(isErrValidate && errValidate != nil){
			return createdUser, errValidate
		}
	}

	username := user.Username

	if len(user.Username) == 0 {
		username = "unnamed"
	}

	result,err := conn.Exec(ctx, ExecuteInsertUser, user.Address, username)
	if err != nil {
        return createdUser, fmt.Errorf("addUser: %v", err)
    }

	rowsAffected := result.RowsAffected()

	fmt.Printf("rows affected %d\n", rowsAffected)

	if (rowsAffected > 0) {
		row := conn.QueryRow(ctx, QueryUserByAddressWithoutAlbum, user.Address)
		if errRow := row.Scan(&createdUser.ID, &createdUser.Address, &createdUser.Username, &createdUser.CreatedAt, &createdUser.Role); errRow != nil {
			if errRow == pgx.ErrNoRows {
				return createdUser, fmt.Errorf("user not found : %v", user.Address)
			}
			return createdUser, fmt.Errorf("failed to get the user: %v", err)
		}
	}

	return createdUser, nil
}

func GetAllUsers(conn *pgxpool.Pool, ctx context.Context)([]CombinedUserWithAlbum, error) {
	var users []CombinedUserWithAlbum

	rows,err := conn.Query(ctx, QueryAllUsers)
	if(err != nil){
		return users, fmt.Errorf("users : %v", err)
	}

	defer rows.Close()

	hasRows := rows.Next()

	for hasRows {
		var user CombinedUserWithAlbum
		if err := rows.Scan(&user.ID, &user.Address, &user.Username, &user.CreatedAt, &user.Role, &user.Albums); err != nil {
			return users, fmt.Errorf("unable to get next user: %v", err)
		}

		hasRows = rows.Next() // Update the hasRows variable

		if len(user.Albums) == 1 && isDefaultAlbum(user.Albums[0]) {
			user.Albums = []_albums.Album{}
		}

		users = append(users, user)
	}

	// Handle the case when there are no rows
	if len(users) == 0 {
		return []CombinedUserWithAlbum{}, nil
	}

	if rows.Err(); err != nil {
		return users, fmt.Errorf(`users rows error: %v`, err)
	}

	return users, nil
}

func GetUserByAddress(address string, conn *pgxpool.Pool, ctx context.Context)(CombinedUserWithAlbum, error) {
	var user CombinedUserWithAlbum

	row := conn.QueryRow(ctx, QueryUserByAddress, address)
	if err := row.Scan(&user.ID, &user.Address, &user.Username, &user.CreatedAt, &user.Role, &user.Albums); err != nil {
		if err == pgx.ErrNoRows {
			return user, fmt.Errorf("userById %v: no such user", address)
		}
		return user, fmt.Errorf("userById %v: %v", address, err)
	}
	
	if isDefaultAlbum(user.Albums[0]){
		user.Albums = []_albums.Album{}
	}

	return user, nil
}

func RemoveUser(address string, conn *pgxpool.Pool, ctx context.Context)(string, error) {
	if len(address) == 0 {
		errMsg := fmt.Sprintf("%v can't be empty", address)
		return "", fmt.Errorf(errMsg)
	}

	result,err := conn.Exec(ctx, QueryRemoveUser, address)
	
	if(err != nil) {
		return "", fmt.Errorf("delete user %v: %v", address, err)
	}

	rowsAffected := result.RowsAffected()

	if(rowsAffected == 0) {
		return "", fmt.Errorf("failed to delete user %d: no rows affected", rowsAffected)
	}

	return "SUCCESS", nil
}