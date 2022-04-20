package spanner_swamp_test

import (
	"context"
	"math/rand"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sinmetal/spanner_swamp"
	"github.com/sinmetal/spanner_swamp/models"
)

func TestShardingSampleStore_Insert(t *testing.T) {
	ctx := context.Background()
	s := newShardingSampleStore(t)

	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	shardID := rand.Int63n(10)
	got, err := s.Insert(ctx, &models.ShardingSample{
		ShardingSampleID: id,
		ShardID:          shardID,
		CreatedAt:        spanner.CommitTimestamp,
		UpdatedAt:        spanner.CommitTimestamp,
	})
	if err != nil {
		t.Fatalf("id=%s, err=%s", id, err)
	}
	e := &models.ShardingSample{
		ShardingSampleID: id,
		ShardID:          shardID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if d := cmp.Diff(e, got, cmpopts.EquateApproxTime(10*time.Second)); d != "" {
		t.Errorf("ShardingSample exist diff\n%s", d)
	}
}

func newShardingSampleStore(t *testing.T) *spanner_swamp.ShardingSampleStore {
	ctx := context.Background()

	s, err := spanner_swamp.NewShardingSampleStore(ctx, spannerClient)
	if err != nil {
		t.Fatal(err)
	}
	return s
}
