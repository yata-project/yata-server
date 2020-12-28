package jwk

import "fmt"

type JWKNotFoundError struct {
	Kid string
}

func (err *JWKNotFoundError) Error() string {
	return fmt.Sprintf("Key %s not found", err.Kid)
}
