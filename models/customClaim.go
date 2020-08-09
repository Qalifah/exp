package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type CustomClaim struct {
	UserID 	uint
	Name	string
	Email	string
	*jwt.StandardClaims
}