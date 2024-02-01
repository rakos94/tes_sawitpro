// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	Registration(ctx context.Context, input RegistrationInput) (output RegistrationOutput, err error)
	Login(ctx context.Context, input LoginInput) (output LoginOutput, err error)
	AddNumLogin(ctx context.Context, id int) (err error)
	Profile(ctx context.Context, input ProfileInput) (output ProfileOutput, err error)
	UpdateProfile(ctx context.Context, input UpdateProfileInput) (err error)
}
