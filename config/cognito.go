package config

import "fmt"

type AwsCognitoUserPoolConfig struct {
	AppClientID string
	Region      string
	UserPoolID  string
}

func (cfg AwsCognitoUserPoolConfig) GetJWKEndpoint() string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", cfg.Region, cfg.UserPoolID)
}
