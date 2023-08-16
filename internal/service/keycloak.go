package service

import (
	"context"
	"falconapi/internal/model"
	"falconapi/pkg/logging"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/pkg/errors"
	"strings"
)

type IdentityManager struct {
	BaseUrl             string
	Realm               string
	RestApiClientId     string
	RestApiClientSecret string
	Logger              *logging.Logger
}

func NewIdentityManager(cfg *model.KeyCloak, logger *logging.Logger) *IdentityManager {
	return &IdentityManager{
		BaseUrl:             cfg.BaseUrl,
		Realm:               cfg.Realm,
		RestApiClientId:     cfg.RestApi.ClientId,
		RestApiClientSecret: cfg.RestApi.ClientSecret,
		Logger:              logger,
	}
}

func (im *IdentityManager) loginClient(ctx context.Context) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.BaseUrl)

	token, err := client.LoginClient(ctx, im.RestApiClientId, im.RestApiClientSecret, im.Realm)
	if err != nil {
		im.Logger.Errorf("unable to login the rest client: %v", err)
		return nil, errors.Wrap(err, "unable to login the rest client")
	}

	return token, nil
}

func (im *IdentityManager) loginUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.BaseUrl)

	token, err := client.Login(ctx, im.RestApiClientId, im.RestApiClientSecret, im.Realm, username, password)
	if err != nil {
		im.Logger.Errorf("unable to login the user: %v", err)
		return nil, errors.Wrap(err, "unable to login the user")
	}

	return token, err
}

func (im *IdentityManager) CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error) {

	token, err := im.loginClient(ctx)
	if err != nil {
		im.Logger.Errorln(err)
		return nil, err
	}

	client := gocloak.NewClient(im.BaseUrl)

	userId, err := client.CreateUser(ctx, token.AccessToken, im.Realm, user)
	if err != nil {
		im.Logger.Errorf("unable to create the user: %v", err)
		return nil, errors.Wrap(err, "unable to create the user")
	}

	err = client.SetPassword(ctx, token.AccessToken, userId, im.Realm, password, false)
	if err != nil {
		im.Logger.Errorf("unable to set the password for the user: %v", err)
		return nil, errors.Wrap(err, "unable to set the password for the user")
	}

	var roleNameLowerCase = strings.ToLower(role)
	roleKeycloak, err := client.GetRealmRole(ctx, token.AccessToken, im.Realm, roleNameLowerCase)
	if err != nil {
		im.Logger.Errorf("unable to get role by name: '%v'", roleNameLowerCase)
		return nil, errors.Wrap(err, fmt.Sprintf("unable to get role by name: '%v'", roleNameLowerCase))
	}
	err = client.AddRealmRoleToUser(ctx, token.AccessToken, im.Realm, userId, []gocloak.Role{
		*roleKeycloak,
	})
	if err != nil {
		im.Logger.Errorf("unable to add a Realm role to user: %v", err)
		return nil, errors.Wrap(err, "unable to add a Realm role to user")
	}

	userKeycloak, err := client.GetUserByID(ctx, token.AccessToken, im.Realm, userId)
	if err != nil {
		im.Logger.Errorf("unable to get recently created user: %v", err)
		return nil, errors.Wrap(err, "unable to get recently created user")
	}

	return userKeycloak, nil
}

func (im *IdentityManager) LoginUser(ctx context.Context, username, password string) (*gocloak.User, *gocloak.JWT, error) {

	client := gocloak.NewClient(im.BaseUrl)

	userToken, err := im.loginUser(ctx, username, password)
	if err != nil {
		im.Logger.Errorln(err)
		return nil, nil, err
	}

	clientToken, err := im.loginClient(ctx)
	if err != nil {
		im.Logger.Errorln(err)
		return nil, nil, err
	}

	userInfo, err := client.GetUserInfo(ctx, userToken.AccessToken, im.Realm)
	if err != nil {
		im.Logger.Errorln(err)
		return nil, nil, err
	}

	userByID, err := client.GetUserByID(ctx, clientToken.AccessToken, im.Realm, *userInfo.Sub)
	if err != nil {
		im.Logger.Errorln(err)
		return nil, nil, err
	}

	return userByID, userToken, nil
}

func (im *IdentityManager) UpdateUserAttributes(ctx context.Context, userID, key, value string) error {
	client := gocloak.NewClient(im.BaseUrl)

	clientToken, err := im.loginClient(ctx)
	if err != nil {
		im.Logger.Errorln(err)
		return err
	}

	userByID, err := client.GetUserByID(ctx, clientToken.AccessToken, im.Realm, userID)
	if err != nil {
		im.Logger.Errorln(err)
		return err
	}

	attribute := make(map[string][]string)
	attribute[key] = []string{value}

	userByID.Attributes = &attribute

	err = client.UpdateUser(ctx, clientToken.AccessToken, im.Realm, *userByID)
	if err != nil {
		im.Logger.Errorln(err)
		return err
	}

	return nil
}

func (im *IdentityManager) RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {

	client := gocloak.NewClient(im.BaseUrl)

	rptResult, err := client.RetrospectToken(ctx, accessToken, im.RestApiClientId, im.RestApiClientSecret, im.Realm)
	if err != nil {
		im.Logger.Errorf("unable to retrospect token: %v", err)
		return nil, errors.Wrap(err, "unable to retrospect token")
	}
	return rptResult, nil
}
