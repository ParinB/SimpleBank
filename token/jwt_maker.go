package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const minSecretKeySize = 32

//JWTMaker is a Json Web Token Maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker created a new JWTMaker
func NewJWTMaker(secretKey string) (Maker , error)  {
	if len(secretKey) < minSecretKeySize {
		return nil , fmt.Errorf("invalid key size : must be at least %d characters",minSecretKeySize)
	}
	return  &JWTMaker{secretKey},nil
}
// CreateToken creates a new token for a specific username and duration
func ( Maker *JWTMaker) CreateToken(username string,duration time.Duration) (string,error) {
	payload , err := NewPayload(username,duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,payload)
	return jwtToken.SignedString([] byte(Maker.secretKey))

}
func (Maker *JWTMaker)  VerifyToken(token string) (*Payload,error) {
	// callback function used to submit key for verification
	//in key method we can find its signing algo  via token.Method field
	keyFunc := func(token *jwt.Token) (interface{},error) {
		//its type is a SigningMethod which is just an interface
		//Convert it to signing method hmac  because we are using  SigningMethodES256
		//which  is instance of method hmac struct
		_,ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			// means algo or token does not match with out signing algo
			return nil ,ErrInvalidToken
		}
		return []byte(Maker.secretKey),nil
	}
	jwtToken,err := jwt.ParseWithClaims(token,&Payload{},keyFunc)
	if err != nil {
		verr,ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner,ErrExpiredToken) {
			return nil,ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload,ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil,ErrInvalidToken
	}
	return payload,nil
}