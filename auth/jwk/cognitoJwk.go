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

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	"github.com/TheYeung1/yata-server/config"
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
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	kid, err := getKid(token)
	if err != nil {
		log.Error("Could not find key in header")
		return nil, errors.New("Key does not exist")
	}

	// TODO: cache the keys
	if err := jwks.Populate(); err != nil {
		log.Error(err)
		return nil, err
	}

	key, ok := jwks.keys[kid]
	if !ok {
		return nil, &JWKNotFoundError{Kid: kid}
	}
	return key.ToSigningKey()
}

func (jwks *AwsCognitoJWKSet) Populate() error {
	resp, err := http.Get(jwks.Config.GetJWKEndpoint())
	if err != nil {
		log.Error("Could not get the json keys")
		return errors.New("Could not get json web token")
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Error("Could not read the json keys response")
		return errors.New("Could not read json keys")
	}

	var cognitoJwks cognitoJWKMap
	err = json.Unmarshal(b, &cognitoJwks)
	if err != nil {
		log.Error("Could not unmarshal json keys")
		return errors.New("Could not unmarshal json keys")
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
		log.Error(err)
		return nil, err
	}
	decodedExponent, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		log.Error(err)
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
		return "", errors.New("Invalid type for kid")
	}
	return kidStr, nil
}
