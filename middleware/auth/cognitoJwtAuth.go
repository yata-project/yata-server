package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/TheYeung1/yata-server/auth/jwk"
	"github.com/TheYeung1/yata-server/config"
	"github.com/TheYeung1/yata-server/server/request"
	"github.com/dgrijalva/jwt-go"
)

type CognitoJwtAuthMiddleware struct {
	Cfg config.AwsCognitoUserPoolConfig
}

func (middleware CognitoJwtAuthMiddleware) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := request.Logger(r.Context())
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			log.Error("request does not contain an Auth token")
			writeBadRequestHttpResponse(w, "Authorization Missing")
			return
		}
		if err := validateBearerToken(tokenString); err != nil {
			log.WithError(err).Error("bearer token not valid")
			writeBadRequestHttpResponse(w, "Bearer token missing")
			return
		}
		tokenString = trimBearerPrefix(tokenString)

		keys := jwk.AwsCognitoJWKSet{Config: middleware.Cfg}

		token, err := jwt.ParseWithClaims(tokenString, &CognitoJwtClaims{}, keys.GetValidationKey)
		if err != nil {
			log.WithError(err).Error("failed to parse JWT token")
			writeUnauthorizedHttpResponse(w, "Invalid JWT token")
			return
		}
		cognitoClaims, ok := token.Claims.(*CognitoJwtClaims)
		if !ok {
			log.WithField("claims", token.Claims).Error("claims could not be casted to Cognito type")
			writeInternalErrorResponse(w, "Sorry! Something went wrong")
			return
		}

		if err = validateClaims(cognitoClaims, middleware.Cfg); err != nil {
			log.WithError(err).WithField("claims", token.Claims).Error("claims are not valid")
			writeUnauthorizedHttpResponse(w, "Invalid Claims")
			return
		}

		log.WithField("claims", cognitoClaims).Debug("claims validated")

		r = r.WithContext(request.WithUserID(r.Context(), cognitoClaims.Subject))

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			log.WithField("token", token).Error("token is not valid")
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

func validateClaims(claims *CognitoJwtClaims, cfg config.AwsCognitoUserPoolConfig) error {
	if !claims.VerifyAudience(cfg.AppClientID, true) {
		return errors.New("invalid audience")
	}
	if !claims.VerifyIssuer(getTokenIssuer(cfg), true) {
		return errors.New("invalid issuer")
	}
	if !claims.VerifyTokenUse("id") {
		return errors.New("invalid token use")
	}
	return nil
}

func getTokenIssuer(cfg config.AwsCognitoUserPoolConfig) string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cfg.Region, cfg.UserPoolID)
}

func writeUnauthorizedHttpResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(msg))
}

func writeBadRequestHttpResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(msg))
}

func writeInternalErrorResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(msg))
}
