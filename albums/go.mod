module playlist/albums

go 1.20

replace playlist/shared => ../shared

require (
	github.com/jackc/pgx/v5 v5.4.0
	playlist/shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.0 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/text v0.10.0 // indirect
)
