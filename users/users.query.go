package users

var ExecuteInsertUser = `INSERT INTO "User" (address, username) VALUES ($1, $2)`

var QueryAllUsers = `SELECT
		"User".id,
		"User".address,
		"User".username,
		"User".createdAt,
		"User".role,
		COALESCE(json_agg(album_info), json_agg('[]'::json)) AS albums
	FROM
		"User"
	LEFT JOIN (
		SELECT
			"Album".id,
			"Album".title,
			"Album".artist,
			"Album".price,
			"Album".user_id
		FROM
			"Album"
		) AS album_info ON "User".id = album_info.user_id
	GROUP BY
		"User".id, "User".address, "User".username, "User".createdAt
	ORDER BY
		"User".id ASC;
`

var QueryUserByAddress = `SELECT
		"User".id,
		"User".address,
		"User".username,
		"User".createdAt,
		"User".role,
	json_agg(json_build_object(
		'id', "Album".id,
		'title', "Album".title,
		'artist', "Album".artist,
		'price', "Album".price,
		'user_id', "Album".user_id
	)) AS album
	FROM
		"User"
	LEFT JOIN
		"Album" ON "User".id = "Album".user_id
	WHERE
		"User".address = $1
	GROUP BY
		"User".id;
`

var QueryUserByAddressWithoutAlbum = `SELECT * FROM "User" WHERE address = $1`

var QueryRemoveUser = `DELETE FROM "User" where address = $1`