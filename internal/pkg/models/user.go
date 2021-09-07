package models

import (
	"time"

	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/sanitizer"
)

var CodeExpires = time.Minute * 10

type SignupUserRequest struct {
	Username string `json:"username" valid:"stringlength(1|32)"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"password"`
}

func (s *SignupUserRequest) Sanitize() {
	snt := sanitizer.NewSanitizer()
	s.Email = snt.Sanitize(s.Email)
	s.Password = snt.Sanitize(s.Password)
	s.Username = snt.Sanitize(s.Username)
}

type LoginUserRequest struct {
	EmailOrUsername string `json:"emailOrUsername" valid:"email_or_username"`
	Password        string `json:"password" valid:"password"`
}

func (s *LoginUserRequest) Sanitize() {
	snt := sanitizer.NewSanitizer()
	s.EmailOrUsername = snt.Sanitize(s.EmailOrUsername)
	s.Password = snt.Sanitize(s.Password)
}

type UserData struct {
	Id       uint64     `json:"-"`
	Email    string     `json:"email" valid:"email"`
	Username string     `json:"username" valid:"stringlength(1|32)"`
	Avatar   UserAvatar `json:"avatar" valid:"required"`
	Password []byte     `json:"-"`
}

func (s *UserData) Sanitize() {
	snt := sanitizer.NewSanitizer()
	s.Username = snt.Sanitize(s.Username)
	s.Email = snt.Sanitize(s.Email)
	s.Avatar.Sanitize()
}

type UserAvatar struct {
	Url string `json:"url" valid:"url"`
}

func (s *UserAvatar) Sanitize() {
	snt := sanitizer.NewSanitizer()
	s.Url = snt.Sanitize(s.Url)
}
