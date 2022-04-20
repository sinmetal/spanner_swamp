package spanner_swamp

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/sinmetal/spanner_swamp/models"
)

type ShardingSampleStore struct {
	sc *spanner.Client
}

func NewShardingSampleStore(ctx context.Context, sc *spanner.Client) (*ShardingSampleStore, error) {
	return &ShardingSampleStore{
		sc: sc,
	}, nil
}

func (s *ShardingSampleStore) Insert(ctx context.Context, v *models.ShardingSample) (*models.ShardingSample, error) {
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
