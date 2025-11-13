package sentinel

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var chConn driver.Conn

func InitClickHouse() {
	var err error
	chConn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{"clickhouse:9000"},
		Auth: clickhouse.Auth{
			Database: "sentinel",
		},
		Debug: true,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
	})

	if err != nil {
		log.Fatalf("Error connecting to ClickHouse: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := chConn.Ping(ctx); err != nil {
		log.Fatalf("Error pinging ClickHouse: %v", err)
	}

	fmt.Println("Successfully connected to ClickHouse!")
}
