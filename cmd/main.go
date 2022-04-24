package main

import (
	"context"
	"flag"
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"github.com/sinmetal/spanner_swamp"
	"github.com/sinmetal/spanner_swamp/models"
)

func main() {
	ctx := context.Background()

	maxUserIDVar := flag.String("max-user-id", "0", "")

	flag.Parse()

	maxUserID, err := strconv.Atoi(*maxUserIDVar)
	if err != nil {
		panic(err)
	}
	fmt.Printf("max-user-id:%d\n", maxUserID)

	config := spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened: 1,
			MaxOpened: 10,
		},
	}

	sc, err := spanner.NewClientWithConfig(ctx, "projects/gcpug-public-spanner/instances/merpay-sponsored-instance/databases/sinmetal", config)
	if err != nil {
		panic(err)
	}
	defer sc.Close()

	shardingSampleStore, err := spanner_swamp.NewShardingSampleStore(ctx, sc)
	if err != nil {
		panic(err)
	}
	userStore, err := spanner_swamp.NewUserStore(ctx, sc)
	if err != nil {
		panic(err)
	}

	for i := 1; i < maxUserID; i++ {
		userID := fmt.Sprintf("u%09d", i)
		_, err = userStore.Insert(ctx, &models.User{
			UserID:      userID,
			MailAddress: fmt.Sprintf("example-%s@example.com", userID),
			FirstName:   "Taro",
			LastName:    "Spanner",
			CreatedAt:   spanner.CommitTimestamp,
			UpdatedAt:   spanner.CommitTimestamp,
		})
		fmt.Printf("User:%s:%s\n", time.Now(), userID)
	}

	for {
		id := strings.ReplaceAll(uuid.New().String(), "-", "")
		shardID := int64(crc32.ChecksumIEEE([]byte(id)) % 10)
		_, err := shardingSampleStore.Insert(ctx, &models.ShardingSample{
			ShardingSampleID: id,
			ShardID:          shardID,
			CreatedAt:        spanner.CommitTimestamp,
			UpdatedAt:        spanner.CommitTimestamp,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("ShardingSample:%s:%s\n", time.Now(), id)

		time.Sleep(time.Duration(1+rand.Int31n(100)) * time.Millisecond)
	}
}
