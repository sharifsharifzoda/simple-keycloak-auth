package otp

import (
	"bytes"
	"github.com/pquerna/otp/totp"
	"image/png"
)

func GenerateOtp() (string, bytes.Buffer, error) {
	var qrCode bytes.Buffer

	// генерация OTP
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "humopay",
		AccountName: "humopay",
		SecretSize:  15,
	})
	if err != nil {
		return "", qrCode, err
	}

	// создание QR - кода
	img, err := key.Image(200, 200)
	if err != nil {
		return "", qrCode, err
	}

	err = png.Encode(&qrCode, img)
	if err != nil {
		return "", qrCode, err
	}

	return key.Secret(), qrCode, nil
}

func ValidateOtp(passcode, secret string) bool {
	return totp.Validate(passcode, secret)
}
