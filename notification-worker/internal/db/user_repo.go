package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxUserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *PgxUserRepository {
	return &PgxUserRepository{Pool: pool}
}

func (repo *PgxUserRepository) GetUserEmailsByIds(ctx context.Context, userIds []string) ([]string, error) {
	ids := make([]interface{}, len(userIds))
	for i, id := range userIds {
		ids[i] = id
	}

	const getEmailsSQL = `
        SELECT email FROM users WHERE id = ANY($1);
    `

	rows, err := repo.Pool.Query(ctx, getEmailsSQL, ids)
	if err != nil {
		return nil, fmt.Errorf("error querying user emails: %w", err)
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("error scanning email: %w", err)
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return emails, nil
}

func (repo *PgxUserRepository) GetUserPhonesByIds(ctx context.Context, userIds []string) ([]string, error) {
	ids := make([]interface{}, len(userIds))
	for i, id := range userIds {
		ids[i] = id
	}

	const getPhoneNumberSQL = `
        SELECT phone_number FROM users WHERE id = ANY($1);
    `

	rows, err := repo.Pool.Query(ctx, getPhoneNumberSQL, ids)
	if err != nil {
		return nil, fmt.Errorf("error querying user phoneNumber: %w", err)
	}
	defer rows.Close()

	var phone_numbers []string
	for rows.Next() {
		var phone_number string
		if err := rows.Scan(&phone_number); err != nil {
			return nil, fmt.Errorf("error scanning phone numbers: %w", err)
		}
		phone_numbers = append(phone_numbers, phone_number)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return phone_numbers, nil
}
