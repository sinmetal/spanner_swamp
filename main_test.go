package spanner_swamp_test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/google/uuid"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

var spannerClient *spanner.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	projectID := "unittest"
	instanceName := fmt.Sprintf("tins-%s", strings.ReplaceAll(uuid.New().String(), "-", ""))
	dbName := "test"
	dbFullName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceName, dbName)

	// ddls, err := readDDLs("./ddl/example.sql", "./ddl/sharding-sample.sql")
	ddls, err := readDDLs("./ddl/sharding-sample.sql")
	if err != nil {
		panic(err)
	}
	for _, ddl := range ddls {
		fmt.Println(ddl)
		fmt.Println("-----------------------------------")
	}

	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		panic(err)
	}
	ope1, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceName,
	})
	if err != nil {
		panic(err)
	}
	_, err = ope1.Wait(ctx)
	if err != nil {
		panic(err)
	}

	dbAdmin, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		panic(err)
	}
	ope2, err := dbAdmin.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", projectID, instanceName),
		CreateStatement: fmt.Sprintf("CREATE DATABASE %s", dbName),
		ExtraStatements: ddls,
	})
	if err != nil {
		panic(err)
	}
	_, err = ope2.Wait(ctx)
	if err != nil {
		panic(err)
	}

	config := spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened: 1,
			MaxOpened: 10,
		},
	}
	sc, err := spanner.NewClientWithConfig(ctx, dbFullName, config)
	if err != nil {
		panic(err)
	}
	spannerClient = sc
	defer func() {
		spannerClient.Close()
	}()

	m.Run()
}

// readDDLs is DDLが書かれたファイルを読み込んで、CreateDatabaseのExtraStatementsに渡せるよう状態にする
// DDLのコメントは無視してくれないので、頭が `#` の場合、その行は無視する
func readDDLs(filePaths ...string) ([]string, error) {
	var ret []string
	for _, fp := range filePaths {
		f, err := os.Open(fp)
		if err != nil {
			return nil, fmt.Errorf("failed read %s: %w", fp, err)
		}

		var sb strings.Builder
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			txt := scanner.Text()
			v := strings.TrimSpace(txt)
			if strings.HasPrefix(v, "#") {
				continue
			}
			sb.WriteString(txt)
			sb.WriteString("\n")
		}

		if err = scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed scan %s: %w", fp, err)
		}
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
		l := strings.Split(sb.String(), ";")

		for _, v := range l {
			if len(strings.TrimSpace(v)) < 1 {
				continue
			}
			ret = append(ret, v)
		}
	}
	return ret, nil
}
