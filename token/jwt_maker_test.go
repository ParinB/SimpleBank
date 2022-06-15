package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/parin/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T)  {
	maker,err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username,duration)
	fmt.Println("Err >>>>>>>>",err)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload,err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username ,payload.Username)
	require.WithinDuration(t, issuedAt,payload.IssuedAt,time.Second)
	require.WithinDuration(t, expiredAt,payload.ExpiredAt,time.Second)


}
func TestExpiredJWTToken(t *testing.T) {
	maker ,err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token,err := maker.CreateToken(util.RandomOwner(),-time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload,err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err , ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidtJWTTokenAlgNone(t *testing.T)  {
	payload,err := NewPayload(util.RandomOwner(),time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone,payload)
	token,err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker , err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload,err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err ,ErrInvalidToken.Error())
	require.Nil(t, payload)
}
