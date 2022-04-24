package spanner_swamp_test

import (
	"context"
	"os"
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

func TestUserStore_Insert(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	s := newUserStore(t)

	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	const mailAddress = "example@example.com"
	const firstName = "Taro"
	const LastName = "Spanner"
	got, err := s.Insert(ctx, &models.User{
		UserID:      id,
		MailAddress: mailAddress,
		FirstName:   firstName,
		LastName:    LastName,
		CreatedAt:   spanner.CommitTimestamp,
		UpdatedAt:   spanner.CommitTimestamp,
	})
	if err != nil {
		t.Fatalf("id=%s, err=%s", id, err)
	}
	e := &models.User{
		UserID:      id,
		MailAddress: mailAddress,
		FirstName:   firstName,
		LastName:    LastName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if d := cmp.Diff(e, got, cmpopts.EquateApproxTime(10*time.Second)); d != "" {
		t.Errorf("User exist diff\n%s", d)
	}
}

func newUserStore(t *testing.T) *spanner_swamp.UserStore {
	ctx := context.Background()

	s, err := spanner_swamp.NewUserStore(ctx, spannerClient)
	if err != nil {
		t.Fatal(err)
	}
	return s
}
