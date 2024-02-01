package handler

import (
	"crypto/rsa"
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const SALT = "salt"

var RSA *rsa.PrivateKey

// Handle login user.
// (POST /auth/login)
func (s *Server) Login(ctx echo.Context) error {
	var req generated.LoginRequest
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "could not bind request body")
	}

	// check phone
	out, err := s.Repository.Login(ctx.Request().Context(), repository.LoginInput{
		Phone: req.Phone,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(out.Password), []byte(req.Password+SALT))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "wrong password"})
	}

	// set token claim
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id": out.Id,
	})

	// get private key
	privateKey, err := os.ReadFile("rsakey.pem")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	// generate token
	token, err := t.SignedString(key)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	// add num login
	err = s.Repository.AddNumLogin(ctx.Request().Context(), out.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	resp := generated.LoginResponse{
		Id:    out.Id,
		Token: token,
	}
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
	err = checkErrorValidation(ctx, validateRegistrationRequest(req))
	if err != nil {
		return err
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
	id, err := parseTokenFromSignedCtx(ctx)
	if err != nil {
		return err
	}

	out, err := s.Repository.Profile(ctx.Request().Context(), repository.ProfileInput{
		Id: id,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	resp := generated.GetProfileResponse{
		Name:  out.Name,
		Phone: out.Phone,
	}
	return ctx.JSON(http.StatusOK, resp)
}

// Handle update profile user.
// (PATCH /profile)
func (s *Server) UpdateProfile(ctx echo.Context) error {
	id, err := parseTokenFromSignedCtx(ctx)
	if err != nil {
		return err
	}

	var req generated.UpdateProfileRequest
	err = ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "could not bind request body")
	}

	// validation
	err = checkErrorValidation(ctx, validateUpdateProfile(req))
	if err != nil {
		return err
	}

	updateProfileInput := repository.UpdateProfileInput{
		Id: id,
	}
	if req.Name != nil {
		updateProfileInput.Name = *req.Name
	}
	if req.Phone != nil {
		updateProfileInput.Phone = *req.Phone
	}
	err = s.Repository.UpdateProfile(ctx.Request().Context(), updateProfileInput)
	if err != nil {
		if strings.Contains(err.Error(), "unique_phone") {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "phone already exist"})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	var resp generated.Success
	resp.Message = "update profile success"
	return ctx.JSON(http.StatusOK, resp)
}

func checkErrorValidation(ctx echo.Context, errorValidations map[string]string) error {
	if len(errorValidations) > 0 {
		validations := []generated.Validations{}
		for k, v := range errorValidations {
			validations = append(validations, generated.Validations{Field: k, Error: v})
		}
		return ctx.JSON(http.StatusBadRequest, generated.ErrorValidationResponse{Errors: validations})
	}

	return nil
}
