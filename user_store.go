package spanner_swamp

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/sinmetal/spanner_swamp/models"
)

type UserStore struct {
	sc *spanner.Client
}

func NewUserStore(ctx context.Context, sc *spanner.Client) (*UserStore, error) {
	return &UserStore{
		sc: sc,
	}, nil
}

func (s *UserStore) Insert(ctx context.Context, v *models.User) (ret *models.User, err error) {
	mu, err := spanner.InsertStruct(v.Table(), v)
	if err != nil {
		return nil, err
	}
	commitTimestamp, err := s.sc.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		if err := tx.BufferWrite([]*spanner.Mutation{mu}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	v.CreatedAt = commitTimestamp
	v.UpdatedAt = commitTimestamp
	return v, nil
}
