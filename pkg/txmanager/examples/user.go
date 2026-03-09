package examples

import (
	"context"
	"media-equipment-tracker/pkg/txmanager/pgtx"

	"github.com/Masterminds/squirrel"
)

type User struct {
	ID   int
	Name string
	Age  int

	Money int
}

type UserRepository struct {
	DB *pgtx.DB
}

func (r *UserRepository) Get(ctx context.Context, id int) (*User, error) {
	var user User

	q, err := r.DB.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	query, args, err := squirrel.Select("id", "name", "age", "money").
		PlaceholderFormat(squirrel.Dollar).
		From("users").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	err = q.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Age, &user.Money)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	q, err := r.DB.GetConn(ctx)
	if err != nil {
		return err
	}

	query, args, err := squirrel.Update("users").
		PlaceholderFormat(squirrel.Dollar).
		Set("name", user.Name).
		Set("age", user.Age).
		Set("money", user.Money).
		Where(squirrel.Eq{"id": user.ID}).ToSql()
	if err != nil {
		return err
	}

	_, err = q.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
