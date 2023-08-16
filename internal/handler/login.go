package handler

import (
	"falconapi/internal/model"
	"falconapi/pkg/otp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) LoginUser(c *gin.Context) {
	var user *model.User

	if err := c.BindJSON(&user); err != nil {
		h.Logger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "please send valid data"})
		return
	}

	/*
		client := h.Client

		token, err := client.Login(c.Request.Context(), h.Cfg.KeyCloak.RestApi.ClientId, h.Cfg.KeyCloak.RestApi.ClientSecret, h.Cfg.KeyCloak.Realm, user.Username, user.Password)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		clientToken, err := client.LoginClient(c.Request.Context(), h.Cfg.KeyCloak.RestApi.ClientId, h.Cfg.KeyCloak.RestApi.ClientSecret, h.Cfg.KeyCloak.Realm)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		userInfo, err := client.GetUserInfo(c.Request.Context(), token.AccessToken, h.Cfg.KeyCloak.Realm)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		fmt.Println("userinfo: ", userInfo)

		userByID, err := client.GetUserByID(c.Request.Context(), clientToken.AccessToken, h.Cfg.KeyCloak.Realm, *userInfo.Sub)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

	*/

	keycloakUser, _, err := h.Service.LoginUser(c.Request.Context(), user.Username, user.Password)
	if err != nil {
		h.Logger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "internal server error"})
		return
	}

	if keycloakUser.Attributes == nil {
		c.JSON(http.StatusOK, gin.H{"message": "please generate OTP", "otp_generated": false, "user_id": *keycloakUser.ID})
	} else if keycloakUser.Attributes != nil {
		attributes := *keycloakUser.Attributes
		_, ok := attributes["otp_secret"]
		if !ok {
			c.JSON(http.StatusOK, gin.H{"message": "please generate OTP", "otp_generated": false, "user_id": *keycloakUser.ID})
		} else if ok {
			c.JSON(http.StatusOK, gin.H{"message": "please validate OTP", "otp_generated": true, "user_id": *keycloakUser.ID})
		}
	}

}

func (h *Handler) GenerateOtp(c *gin.Context) {
	var payload *model.OtpInput

	if err := c.BindJSON(&payload); err != nil {
		h.Logger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "internal server error"})
		return
	}

	/*
		// Генерация OTP
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "humopay",
			AccountName: "humopay",
			SecretSize:  15,
		})
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "internal server error"})
			return
		}

		// Генерация QR - кода
		image, err := key.Image(200, 200)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "internal server error"})
			return
		}
		var qrCode bytes.Buffer

		err = png.Encode(&qrCode, image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "internal server error"})
			return
		}

	*/
	otpSecret, buf, err := otp.GenerateOtp()
	if err != nil {
		h.Logger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "internal server error"})
		return
	}

	/*
		// Добавление секретного ключа OTP в KeyCloak
		client := h.Client

		clientToken, err := client.LoginClient(c.Request.Context(), h.Cfg.KeyCloak.RestApi.ClientId, h.Cfg.KeyCloak.RestApi.ClientSecret, h.Cfg.KeyCloak.Realm)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		userByID, err := client.GetUserByID(c.Request.Context(), clientToken.AccessToken, h.Cfg.KeyCloak.Realm, payload.UserId)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		attribute := make(map[string][]string)
		attribute["otp_secret"] = []string{key.Secret()}

		userByID.Attributes = &attribute

		err = client.UpdateUser(c.Request.Context(), clientToken.AccessToken, h.Cfg.KeyCloak.Realm, *userByID)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

	*/
	err = h.Service.UpdateUserAttributes(c.Request.Context(), payload.UserId, "otp_secret", otpSecret)
	if err != nil {
		h.Logger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "internal server error"})
		return
	}

	/*
		c.Header("Content-Type", "image/png")
		c.Header("Content-Length", strconv.Itoa(len(qrCode.Bytes())))

		c.Status(http.StatusOK)
		_, _ = c.Writer.Write(qrCode.Bytes())

	*/

	c.Header("Content-Type", "image/png")
	c.Header("Content-Length", strconv.Itoa(len(buf.Bytes())))

	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(buf.Bytes())

}

func (h *Handler) ValidateOtp(c *gin.Context) {
	var payload *model.OtpInput

	err := c.BindJSON(&payload)
	if err != nil {
		h.Logger.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "internal server error"})
		return
	}

	/*
		client := h.Client

		token, err := client.Login(c.Request.Context(), h.Cfg.KeyCloak.RestApi.ClientId, h.Cfg.KeyCloak.RestApi.ClientSecret, h.Cfg.KeyCloak.Realm, payload.Username, payload.Password)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		clientToken, err := client.LoginClient(c.Request.Context(), h.Cfg.KeyCloak.RestApi.ClientId, h.Cfg.KeyCloak.RestApi.ClientSecret, h.Cfg.KeyCloak.Realm)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		userInfo, err := client.GetUserInfo(c.Request.Context(), token.AccessToken, h.Cfg.KeyCloak.Realm)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		userByID, err := client.GetUserByID(c.Request.Context(), clientToken.AccessToken, h.Cfg.KeyCloak.Realm, *userInfo.Sub)
		if err != nil {
			h.Logger.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		if *userByID.Attributes == nil {
			h.Logger.Println("empty attributes")
			c.JSON(http.StatusForbidden, "unknown error")
			return
		}

		attributes := *userByID.Attributes
		secretKey := attributes["otp_secret"][0]

		validate := totp.Validate(payload.OtpToken, secretKey)
		if !validate {
			h.Logger.Println("invalid otp token provided")
			c.JSON(http.StatusBadRequest, "invalid otp token provided")
			return
		}

	*/

	keycloakUser, userToken, err := h.Service.LoginUser(c.Request.Context(), payload.Username, payload.Password)
	if err != nil {
		h.Logger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "internal server error"})
		return
	}

	h.Logger.Println(keycloakUser)

	c.JSON(http.StatusOK, gin.H{"access_token": userToken.AccessToken})
}
