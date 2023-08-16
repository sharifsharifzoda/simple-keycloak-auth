package service

import (
	"context"
	"falconapi/internal/model"
	"falconapi/pkg/logging"
	"github.com/Nerzal/gocloak/v13"
)

type KeyCloakService interface {
	CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error)
	LoginUser(ctx context.Context, username, password string) (*gocloak.User, *gocloak.JWT, error)
	UpdateUserAttributes(ctx context.Context, userID, key, value string) error
	RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error)
}

type Service struct {
	KeyCloakService
}

func NewService(cfg *model.Config, logger *logging.Logger) *Service {
	return &Service{
		KeyCloakService: NewIdentityManager(&cfg.KeyCloak, logger),
	}
}
