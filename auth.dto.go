package main

import (
	_users "playlist/users"

	"github.com/golang-jwt/jwt/v5"
)

type CustomTokenClaim struct {
	User _users.User `json:"user"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Address string `json:"address" binding:"required"`
}