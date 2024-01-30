package repository

import "context"

func (r *Repository) Registration(ctx context.Context, input RegistrationInput) (output RegistrationOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO users (name, phone, password) VALUES ($1, $2, $3) RETURNING id", input.Name, input.Phone, input.Password).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}
