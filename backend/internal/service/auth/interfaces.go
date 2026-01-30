package auth

import "eshbuket/internal/transport/http/dto"

type IAuthService interface {
	Authenticate(login, password string) bool
	CreateSession(req dto.LoginRequest) string
}
