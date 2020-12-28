package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	"github.com/TheYeung1/yata-server/auth/jwk"
	"github.com/TheYeung1/yata-server/config"
)

type CognitoJwtAuthMiddleware struct {
	Cfg config.AwsCognitoUserPoolConfig
}

func (middleware CognitoJwtAuthMiddleware) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			log.Error("Request does not contain an Auth token")
			writeBadRequestHttpResponse(w, "Authorization Missing")
			return
		}
		if err := validateBearerToken(tokenString); err != nil {
			log.Error(err)
			writeBadRequestHttpResponse(w, "Bearer token missing")
			return
		}
		tokenString = trimBearerPrefix(tokenString)

		keys := jwk.AwsCognitoJWKSet{Config: middleware.Cfg}

		token, err := jwt.ParseWithClaims(tokenString, &CognitoJwtClaims{}, keys.GetValidationKey)
		if err != nil {
			log.Error("Error parsing JWT token: ", err)
			writeUnauthorizedHttpResponse(w, "Invalid JWT token")
			return
		}
		if err = validateClaims(token.Claims, middleware.Cfg); err != nil {
			log.Errorf("Invalid Claims: %+v", token.Claims, err)
			writeUnauthorizedHttpResponse(w, "Invalid Claims")
			return
		}

		// TODO: set the user ID to be used later?

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			log.Errorf("Invalid token: %+v", token)
			writeUnauthorizedHttpResponse(w, "Invalid Token")
		}
	})
}

func validateBearerToken(token string) error {
	if !strings.HasPrefix(token, "Bearer") {
		return errors.New("Authentication does not use bearer token")
	}
	return nil
}

func trimBearerPrefix(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}

func validateClaims(claims jwt.Claims, cfg config.AwsCognitoUserPoolConfig) error {
	cognitoClaims, ok := claims.(*CognitoJwtClaims)
	if !ok {
		return errors.New("Invalid claims")
	}
	if !cognitoClaims.VerifyAudience(cfg.AppClientID, true) {
		return errors.New("Invalid audience")
	}
	if !cognitoClaims.VerifyIssuer(getTokenIssuer(cfg), true) {
		return errors.New("Invalid issuer")
	}
	if !cognitoClaims.VerifyTokenUse("id") {
		return errors.New("Invalid token use")
	}
	return nil
}

func getTokenIssuer(cfg config.AwsCognitoUserPoolConfig) string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cfg.Region, cfg.UserPoolID)
}

func writeUnauthorizedHttpResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(msg))
}

func writeBadRequestHttpResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}
