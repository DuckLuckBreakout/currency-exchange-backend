package jwt_token

import (
	"fmt"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/configer"

	"github.com/twinj/uuid"
)

type TokenType int

const (
	Access TokenType = iota
	Refresh
)

func CreateJwtToken(secrets *configer.SecretData) (*models.TokenDetails, error) {
	tokenDetails := &models.TokenDetails{}
	Uuid := uuid.NewV4().String()
	// Setting Token Details
	tokenDetails.AccessDetails.Expires = time.Now().Add(models.AccessTokenExpires).Unix()
	tokenDetails.AccessDetails.Uuid = Uuid

	tokenDetails.RefreshDetails.Expires = time.Now().Add(models.RefreshTokenExpires).Unix()
	tokenDetails.RefreshDetails.Uuid = Uuid

	var err error
	// Creating Access Token
	os.Setenv("ACCESS_SECRET", secrets.AccessSecret)
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["user_uuid"] = tokenDetails.AccessDetails.Uuid
	accessTokenClaims["expires"] = tokenDetails.AccessDetails.Expires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tokenDetails.Token.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating Refresh Token
	os.Setenv("REFRESH_SECRET", secrets.RefreshSecret)
	RefreshTokenClaims := jwt.MapClaims{}
	RefreshTokenClaims["user_uuid"] = tokenDetails.RefreshDetails.Uuid
	RefreshTokenClaims["expires"] = tokenDetails.RefreshDetails.Expires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshTokenClaims)
	tokenDetails.Token.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return tokenDetails, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request, tokenType TokenType, secrets *configer.SecretData) (*jwt.Token, error) {
	tokenString := ExtractToken(r)

	var tokenSecret string
	switch tokenType {
	case Access:
		tokenSecret = secrets.AccessSecret
	case Refresh:
		tokenSecret = secrets.RefreshSecret
	default:
		return nil, fmt.Errorf("Incorrect token type")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func TokenValid(r *http.Request, tokenType TokenType, secrets *configer.SecretData) error {
	token, err := VerifyToken(r, tokenType, secrets)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

func ExtractAccessTokenMetadata(r *http.Request, secrets *configer.SecretData) (*models.AccessTokenDetails, error) {
	token, err := VerifyToken(r, Access, secrets)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["user_uuid"].(string)
		if !ok {
			return nil, err
		}
		return &models.AccessTokenDetails{
			Uuid: accessUuid,
		}, nil
	}
	return nil, err
}

func ExtractRefreshTokenMetadata(r *http.Request, secrets *configer.SecretData) (*models.RefreshTokenDetails, error) {
	token, err := VerifyToken(r, Refresh, secrets)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["user_uuid"].(string)
		if !ok {
			return nil, err
		}
		return &models.RefreshTokenDetails{
			Uuid: accessUuid,
		}, nil
	}
	return nil, err
}
