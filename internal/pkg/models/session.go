package models

import (
	"time"

	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/sanitizer"
)

var (
	AccessTokenExpires  = time.Minute * 10
	RefreshTokenExpires = time.Hour * 24 * 7
)

type TokenDetails struct {
	Token          Token               `json:"token"`
	AccessDetails  AccessTokenDetails  `json:"accessDetails"`
	RefreshDetails RefreshTokenDetails `json:"refreshDetails"`
}

func (s *TokenDetails) Sanitize() {
	s.Token.Sanitize()
	s.AccessDetails.Sanitize()
	s.RefreshDetails.Sanitize()
}

type AccessTokenDetails struct {
	Uuid    string `json:"uuid" valid:"required, type(string)"`
	Expires int64  `json:"expires" valid:"type(int64)"`
}

func (s *AccessTokenDetails) Sanitize() {
	snt := sanitizer.NewSanitizer()
	s.Uuid = snt.Sanitize(s.Uuid)
}

type RefreshTokenDetails struct {
	Uuid    string `json:"uuid" valid:"required, type(string)"`
	Expires int64  `json:"expires" valid:"type(int64)"`
}

func (s *RefreshTokenDetails) Sanitize() {
	snt := sanitizer.NewSanitizer()
	s.Uuid = snt.Sanitize(s.Uuid)
}

type Token struct {
	AccessToken  string `json:"accessToken" valid:"required, type(string)"`
	RefreshToken string `json:"refreshToken" valid:"required, type(string)"`
}

func (s *Token) Sanitize() {
	snt := sanitizer.NewSanitizer()
	s.RefreshToken = snt.Sanitize(s.RefreshToken)
	s.AccessToken = snt.Sanitize(s.AccessToken)
}
