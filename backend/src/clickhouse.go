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
	var conn driver.Conn

	for i := 0; i < 5; i++ {
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"clickhouse:9000"},
			Auth: clickhouse.Auth{
				Database: "sentinel",
				Username: "sentinel",
				Password: "password",
			},
			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
		})

		if err != nil {
			log.Printf("Error connecting to ClickHouse: %v. Retrying in 3 seconds...", err)
			time.Sleep(3 * time.Second)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = conn.Ping(ctx); err != nil {
			log.Printf("Error pinging ClickHouse: %v. Retrying in 3 seconds...", err)
			time.Sleep(3 * time.Second)
			continue
		}

		chConn = conn
		fmt.Println("Successfully connected to ClickHouse!")
		return
	}

	log.Fatalf("Could not connect to ClickHouse after several retries: %v", err)
}
