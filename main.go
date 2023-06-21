package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Ctx context.Context
	Cancel context.CancelFunc
	// Conn *pgx.Conn
	Conn *pgxpool.Pool
	ErrConn error
	JwtKey = []byte("vCU!pnmJJg5B%hqTxymj")
)


func main() {
	Ctx, Cancel = context.WithTimeout(context.Background(), 15*time.Minute)

	defer Cancel()

	// Ctx = context.Background()

	urlDb := "postgres://myuser:mypassword@localhost:5432/albums"
	Conn, ErrConn = createConnectionPool(urlDb, 64)

	if ErrConn != nil {
        log.Fatal(ErrConn)
    }

	pingErr := Conn.Ping(Ctx)

	if(pingErr != nil){
		log.Fatal(pingErr)
	}
	fmt.Println("connected")

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	router.POST("/login", login)

	router.GET("/albums", GetAllAlbums)
	router.GET("/users", GetAllUsers)
	
	router.GET("/albums/artist/:name", GetOneAlbumByName)
	router.GET("/albums/:id", GetOneAlbumById)

	router.GET("/users/:address", GetOneUserByAddress)

	auth := router.Group("/auth")

	auth.Use(authMiddleware()) 
	{
		auth.POST("/users",  CreateUser)
		auth.DELETE("/users", DeleteUser)

		auth.POST("/albums", PostCreateAlbum)

		auth.POST("/albums/:id", PostUpdateAlbum)
		auth.DELETE("/albums/:id", DeleteAlbum)
	}

	errRouter := router.Run(":5000")
	if errRouter != nil {
		panic("[Error] failed to start Gin server due to: " + errRouter.Error())
	}
}

func createConnectionPool(connectionString string, maxConnections int32)(*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	config.MaxConns = maxConnections
	config.MaxConnIdleTime = time.Minute * 5 // Set maximum idle time for a connection

	pool, err := pgxpool.NewWithConfig(Ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return pool, nil
}