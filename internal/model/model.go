package model

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserRole struct {
	Username string
	Role     string
	UserID   string
}

type OtpInput struct {
	UserId   string `form:"user_id" json:"user_id"`
	OtpToken string `form:"otp_token" json:"otp_token"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type Config struct {
	Listen   ListenCfg `json:"Listen"`
	KeyCloak KeyCloak  `json:"KeyCloak"`
}

// ListenCfg конфигурации для соединения с сервером
type ListenCfg struct {
	ListenIP   string `json:"ListenIP"`
	ListenPort string `json:"ListenPort"`
}

// KeyCloak конфигурации для соединения с Keycloak
type KeyCloak struct {
	Realm               string  `json:"Realm"`
	BaseUrl             string  `json:"BaseUrl"`
	RestApi             RestApi `json:"RestApi"`
	RealmRS256PublicKey string  `json:"RealmRS256PublicKey"`
}

type RestApi struct {
	ClientId     string `json:"ClientId"`
	ClientSecret string `json:"ClientSecret"`
}
