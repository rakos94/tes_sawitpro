package handler

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const SALT = "salt"

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

// Handle login user.
// (POST /auth/login)
func (s *Server) Login(ctx echo.Context) error {
	var resp generated.LoginResponse
	return ctx.JSON(http.StatusOK, resp)
}

// Handle registration user.
// (POST /auth/registration)
func (s *Server) Registration(ctx echo.Context) error {
	var req generated.RegistrationRequest
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "could not bind request body")
	}

	// validation
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

	if len(errorValidations) > 0 {
		validations := []generated.Validations{}
		for k, v := range errorValidations {
			validations = append(validations, generated.Validations{Field: k, Error: v})
		}
		return ctx.JSON(http.StatusBadRequest, generated.ErrorValidationResponse{Errors: validations})
	}

	// hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password+SALT), bcrypt.DefaultCost)

	res, err := s.Repository.Registration(ctx.Request().Context(), repository.RegistrationInput{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: string(hash),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	var resp generated.RegistrationResponse
	resp.Id = res.Id
	return ctx.JSON(http.StatusOK, resp)
}

// Handle get profile user.
// (GET /profile)
func (s *Server) Profile(ctx echo.Context) error {
	var resp generated.GetProfileResponse
	return ctx.JSON(http.StatusOK, resp)
}

// Handle update profile user.
// (PATCH /profile)
func (s *Server) UpdateProfile(ctx echo.Context) error {
	var resp generated.Success
	return ctx.JSON(http.StatusOK, resp)
}
