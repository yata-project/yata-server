package server

import (
	"net/http"

	"github.com/TheYeung1/yata-server/middleware/auth"
	log "github.com/sirupsen/logrus"
)

type WhoAmIOutput struct {
	EmailVerified bool
	TokenUse      string
	AuthTime      int64
	UserName      string
	GivenName     string
	Email         string

	// Standard JWT claims.
	Audience  string
	ExpiresAt int64
	Id        string
	IssuedAt  int64
	Issuer    string
	NotBefore int64
	Subject   string
}

func (s *Server) WhoAmI(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.CognitoClaims(r.Context())
	if !ok || claims == nil {
		log.Error("failed to get user ID from request context")
		renderInternalServerError(w)
		return
	}

	out := WhoAmIOutput{
		EmailVerified: claims.EmailVerified,
		TokenUse:      claims.TokenUse,
		AuthTime:      claims.AuthTime,
		UserName:      claims.UserName,
		GivenName:     claims.GivenName,
		Email:         claims.Email,

		Audience:  claims.Audience,
		ExpiresAt: claims.ExpiresAt,
		Id:        claims.Id,
		IssuedAt:  claims.IssuedAt,
		Issuer:    claims.Issuer,
		NotBefore: claims.NotBefore,
		Subject:   claims.Subject,
	}
	log.WithField("output", out).Debug("user determined")
	renderJSON(w, http.StatusOK, out)
}
