package pgrepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"ads/internal/user"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func (r *AuthPostgres) CreateUserDb(user user.UserDb) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (*user.UserDb, error) {
	var user user.UserDb
	query := fmt.Sprintf("SELECT id, name, username FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return &user, err
}

func (r *AuthPostgres) CheckUserDb(id int) (*user.UserDb, error) {
	var user user.UserDb
	query := fmt.Sprintf("SELECT id, name, username FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthPostgres) UpdateUserDb(username string, id int) (*user.UserDb, error) {
	var user user.UserDb
	query := fmt.Sprintf("UPDATE %s SET username=$1 WHERE id=$2", usersTable)
	_, err := r.db.Exec(query, username, id)
	if err != nil {
		return nil, err
	}
	query = fmt.Sprintf("SELECT id, name, username FROM %s WHERE id=$1", usersTable)
	err = r.db.Get(&user, query, id)
	return &user, err
}

func (r *AuthPostgres) DeleteUserDb(id int) (*user.UserDb, error) {
	var user user.UserDb
	queryGet := fmt.Sprintf("SELECT id, name, username FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, queryGet, id) 
	if err != nil {
		return nil, err
	}
	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE id=$1", usersTable)
	_, err = r.db.Exec(queryDelete, id)
	return &user, err
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
