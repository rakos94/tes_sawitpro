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
