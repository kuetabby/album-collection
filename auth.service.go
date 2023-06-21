package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	_shared "playlist/shared"
	_users "playlist/users"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

func authentication(body LoginRequest)(string,error) {
	var userData _users.User

	row := Conn.QueryRow(Ctx, `SELECT * from "User" where address = $1`, body.Address)
	if err := row.Scan(&userData.ID, &userData.Address, &userData.Username, &userData.CreatedAt, &userData.Role); err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("invalid address %v", err)
		}

		return "", fmt.Errorf(`error login %v`, err)
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &CustomTokenClaim{
		User: userData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, errTokenString := token.SignedString(JwtKey)

	return tokenString, errTokenString
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
        } else {
			token, err := jwt.ParseWithClaims(tokenString, &CustomTokenClaim{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid method")
				}

				return JwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					errSignature := "invalid token signature"
					c.JSON(http.StatusUnauthorized, _shared.SingleResponse[any]{
						BaseResponse: _shared.BaseResponse{
							Status: "ERROR",
							Message: &errSignature,
							Error: true,
						},
					})
					c.Abort()
					return
				}

				errExpired := "token expired"
				c.JSON(http.StatusUnauthorized, _shared.SingleResponse[any]{
						BaseResponse: _shared.BaseResponse{
							Status: "ERROR",
							Message: &errExpired,
							Error: true,
						},
					})
				c.Abort()
				return
			}

			if !token.Valid {
				errToken := "invalid token"
				c.JSON(http.StatusUnauthorized, _shared.SingleResponse[any]{
						BaseResponse: _shared.BaseResponse{
							Status: "ERROR",
							Message: &errToken,
							Error: true,
						},
					})
				return
			}

			// Set the user information from the token claims to the context
			if claims, ok := token.Claims.(*CustomTokenClaim); ok && token.Valid {
				c.Set("user_role", claims.User.Role)
				c.Set("user_id", claims.User.ID)
				fmt.Printf("User %v", claims.User.ID)
			}

			c.Next()
		}
    }
}


func decodeLogin(c *gin.Context)(*_users.User, error) {
	tokenString := c.GetHeader("Authorization")

	token, err := jwt.ParseWithClaims(tokenString, &CustomTokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	
	if claims, ok := token.Claims.(*CustomTokenClaim); ok && token.Valid {
		return &claims.User, nil
	} else {
		return nil, err
	}
}