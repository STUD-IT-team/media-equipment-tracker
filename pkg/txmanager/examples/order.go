package examples

import (
	"context"
	"media-equipment-tracker/pkg/txmanager/pgtx"

	"github.com/Masterminds/squirrel"
)

type Order struct {
	ID     int
	Name   string
	Sum    int
	UserID int
}

type OrderRepository struct {
	DB *pgtx.DB
}

func (r *OrderRepository) Get(ctx context.Context, id int) (*Order, error) {
	var order Order

	q, err := r.DB.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	query, args, err := squirrel.Select("id", "name", "sum", "user_id").
		PlaceholderFormat(squirrel.Dollar).
		From("orders").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}

	err = q.QueryRow(ctx, query, args...).Scan(&order.ID, &order.Name, &order.Sum, &order.UserID)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) Create(ctx context.Context, order *Order) error {
	q, err := r.DB.GetConn(ctx)
	if err != nil {
		return err
	}

	query, args, err := squirrel.Insert("orders").
		PlaceholderFormat(squirrel.Dollar).
		Columns("id", "name", "sum", "user_id").
		Values(order.ID, order.Name, order.Sum, order.UserID).ToSql()
	if err != nil {
		return err
	}

	_, err = q.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
