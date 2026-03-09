package examples

import (
	"context"
	"fmt"
	"math/rand/v2"
	"media-equipment-tracker/pkg/txmanager"
)

type CreateOrderUseCase struct {
	UserRepo  *UserRepository
	OrderRepo *OrderRepository
	Txm       txmanager.TxManager
}

func (s *CreateOrderUseCase) CreateOrder(ctx context.Context, name string, sum int, userID int) error {
	user, err := s.UserRepo.Get(ctx, userID)
	if err != nil {
		return err
	}

	if user.Money < sum {
		return fmt.Errorf("not enough money")
	}

	user.Money -= sum

	order := &Order{
		ID:     int(rand.Int32()),
		Name:   name,
		Sum:    sum,
		UserID: userID,
	}

	err = s.Txm.WithinTx(ctx, func(ctx context.Context) error {
		err = s.UserRepo.Update(ctx, user)
		if err != nil {
			return err
		}

		err = s.OrderRepo.Create(ctx, order)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
