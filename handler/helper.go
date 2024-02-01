package handler

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func validateRegistrationRequest(req generated.RegistrationRequest) map[string]string {
	// validation name
	errorValidations := map[string]string{}
	if req.Name == "" {
		errorValidations["name"] = "required"
	} else if len(req.Name) < 3 {
		errorValidations["name"] = "min length 3"
	} else if len(req.Name) > 60 {
		errorValidations["name"] = "max length 60"
	}

	// validate phone
	if req.Phone == "" {
		errorValidations["phone"] = "required"
	} else if len(req.Phone) < 10 {
		errorValidations["phone"] = "min length 10"
	} else if len(req.Phone) > 13 {
		errorValidations["phone"] = "max length 13"
	} else {
		if req.Phone[0:3] != "+62" {
			errorValidations["phone"] = "must start with +62"
		}
	}

	// validate password
	if req.Password == "" {
		errorValidations["password"] = "required"
	} else if len(req.Password) < 6 {
		errorValidations["password"] = "min length 6"
	} else if len(req.Password) > 64 {
		errorValidations["password"] = "max length 64"
	} else {
		hasUpper := false
		hasLower := false
		hasNumber := false
		hasSpecial := false

		for _, char := range req.Password {
			if unicode.IsUpper(char) {
				hasUpper = true
			} else if unicode.IsLower(char) {
				hasLower = true
			} else if unicode.IsNumber(char) {
				hasNumber = true
			} else if !unicode.IsLetter(char) {
				hasSpecial = true
			}
		}
		if strings.Contains(req.Password, " ") {
			errorValidations["password"] = "must not contain spaces"
		}
		if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			errorValidations["password"] = "must containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric)"
		}
	}

	return errorValidations
}

func validateUpdateProfile(req generated.UpdateProfileRequest) map[string]string {
	// validation name
	errorValidations := map[string]string{}
	if req.Name != nil {
		name := *req.Name
		if len(name) < 3 {
			errorValidations["name"] = "min length 3"
		} else if len(name) > 60 {
			errorValidations["name"] = "max length 60"
		}
	}

	// validate phone
	if req.Phone != nil {
		phone := *req.Phone
		if len(phone) < 10 {
			errorValidations["phone"] = "min length 10"
		} else if len(phone) > 13 {
			errorValidations["phone"] = "max length 13"
		} else {
			if phone[0:3] != "+62" {
				errorValidations["phone"] = "must start with +62"
			}
		}
	}

	return errorValidations
}

func parseTokenFromSignedCtx(ctx echo.Context) (int, error) {
	reqToken := ctx.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return 0, fmt.Errorf("error parsing token")
	}

	tokenString := splitToken[1]

	publicKey, err := os.ReadFile("public.pem")
	if err != nil {
		return 0, fmt.Errorf("error reading public key file: %v\n", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return 0, fmt.Errorf("error parsing RSA public key: %v\n", err)
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return 0, fmt.Errorf("error parsing token: %v", err)
	}

	claim := parsedToken.Claims.(jwt.MapClaims)
	id := claim["id"].(float64)

	return int(id), nil
}
