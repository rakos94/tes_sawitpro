// This file contains types that are used in the repository layer.
package repository

type RegistrationInput struct {
	Name     string
	Phone    string
	Password string
}

type RegistrationOutput struct {
	Id int
}

type LoginInput struct {
	Phone string
}

type LoginOutput struct {
	Id       int
	Password string
	Token    string
}

type ProfileInput struct {
	Id int
}

type ProfileOutput struct {
	Name  string
	Phone string
}

type UpdateProfileInput struct {
	Id    int
	Name  string
	Phone string
}
