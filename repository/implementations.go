package repository

import (
	"context"
)

func (r *Repository) Registration(ctx context.Context, input RegistrationInput) (output RegistrationOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO users (name, phone, password) VALUES ($1, $2, $3) RETURNING id", input.Name, input.Phone, input.Password).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) Login(ctx context.Context, input LoginInput) (output LoginOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "select id, password from users where phone = $1", input.Phone).Scan(&output.Id, &output.Password)
	if err != nil {
		return
	}
	return
}

func (r *Repository) AddNumLogin(ctx context.Context, id int) (err error) {
	_, err = r.Db.ExecContext(ctx, "update users set num_login = num_login + 1 where id = $1", id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) Profile(ctx context.Context, input ProfileInput) (output ProfileOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "select name, phone from users where id = $1", input.Id).Scan(&output.Name, &output.Phone)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateProfile(ctx context.Context, input UpdateProfileInput) (err error) {
	if input.Name != "" && input.Phone != "" {
		_, err = r.Db.ExecContext(ctx, "update users set name = $1, phone = $2 where id = $3", input.Name, input.Phone, input.Id)
		return
	} else if input.Name != "" {
		_, err = r.Db.ExecContext(ctx, "update users set name = $1 where id = $2", input.Name, input.Id)
		return
	} else if input.Phone != "" {
		_, err = r.Db.ExecContext(ctx, "update users set phone = $1 where id = $2", input.Phone, input.Id)
		return
	}
	return
}
