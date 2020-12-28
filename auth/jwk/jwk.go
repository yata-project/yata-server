package jwk

type JWKSet interface {
	// Returns the JsonWebKey given the key ID. nil is returned if the key is not found
	GetKey(string) JWK
}

type JWK interface {
	ToSigningKey() (interface{}, error)
}
