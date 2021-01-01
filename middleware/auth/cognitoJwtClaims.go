package auth

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type CognitoJwtClaims struct {
	EmailVerified bool   `json:"email_verified"`
	TokenUse      string `json:"token_use"`
	AuthTime      int64  `json:"auth_time"`
	UserName      string `json:"cognito:username"`
	GivenName     string `json:"given_name"`
	Email         string `json:"email"`
	jwt.StandardClaims
}

func (claims *CognitoJwtClaims) VerifyTokenUse(use string) bool {
	return strings.Compare(use, claims.TokenUse) == 0
}
