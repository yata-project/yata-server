package jwk

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/TheYeung1/yata-server/config"
	"github.com/dgrijalva/jwt-go"
)

type AwsCognitoJWKSet struct {
	Config config.AwsCognitoUserPoolConfig
	keys   map[string]AwsCognitoJWK
}

type AwsCognitoJWK struct {
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	E   string `json:"e"`
	N   string `json:"n"`
	Use string `json:"use"`
}

// struct for deserializing the cognito well known jwk json
type cognitoJWKMap struct {
	Keys []AwsCognitoJWK `json:"keys"`
}

func (jwks *AwsCognitoJWKSet) GetValidationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, err := getKid(token)
	if err != nil {
		return nil, errors.New("key does not exist")
	}

	// TODO: cache the keys
	if err := jwks.Populate(); err != nil {
		return nil, err
	}

	key, ok := jwks.keys[kid]
	if !ok {
		return nil, &JWKNotFoundError{Kid: kid}
	}
	return key.ToSigningKey()
}

func (jwks *AwsCognitoJWKSet) Populate() error {
	// TODO: Do not use the default http client; it has no timeout set!
	resp, err := http.Get(jwks.Config.GetJWKEndpoint())
	if err != nil {
		return errors.New("could not get json web token")
	}

	b, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return errors.New("could not read json keys")
	}

	var cognitoJwks cognitoJWKMap
	err = json.Unmarshal(b, &cognitoJwks)
	if err != nil {
		return errors.New("could not unmarshal json keys")
	}

	jwks.keys = make(map[string]AwsCognitoJWK)
	for _, key := range cognitoJwks.Keys {
		jwks.keys[key.Kid] = key
	}
	return nil
}

func (key *AwsCognitoJWK) ToSigningKey() (interface{}, error) {
	decodedModulo, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	decodedExponent, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}
	if len(decodedExponent) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedExponent):], decodedExponent)
		decodedExponent = ndata
	}
	return &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(decodedModulo),
		E: int(binary.BigEndian.Uint32(decodedExponent[:])),
	}, nil
}

func getKid(token *jwt.Token) (string, error) {
	kid, ok := token.Header["kid"]
	if !ok {
		return "", errors.New("kid missing from token header")
	}
	kidStr, ok := kid.(string)
	if !ok {
		return "", errors.New("invalid type for kid")
	}
	return kidStr, nil
}
