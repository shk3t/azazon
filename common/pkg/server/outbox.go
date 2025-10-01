package server

import (
	"common/pkg/consts"
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

type TransactionalOutboxManager struct {
	dbConnPool     *pgxpool.Pool
	kafkaConnector *KafkaConnector
	observer       chan struct{}
	cancelCtx      context.CancelFunc
}

func NewTransactionalOutboxManager(
	dbConnPool *pgxpool.Pool,
	kafkaConnector *KafkaConnector,
) *TransactionalOutboxManager {
	m := &TransactionalOutboxManager{
		dbConnPool:     dbConnPool,
		kafkaConnector: kafkaConnector,
		observer:       make(chan struct{}, 1),
	}

	ticker := time.NewTicker(time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelCtx = cancel

	go func() {
		for {
			select {
			// db by timeout or queue
			case <-m.observer:
			case <-ticker.C:
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return m
}

const trablePostfix = "_outbox"

func getTableName(topic consts.TopicName) string {
	return string(topic) + trablePostfix
}

func (m *TransactionalOutboxManager) processUnprocessed() {
	ctx := context.Background()

	topics := []consts.TopicName{}

	for _, topic := range topics {
		table := getTableName(topic)
		messages := []kafka.Message{}

		tx, _ := m.dbConnPool.BeginTx(ctx, pgx.TxOptions{})
		defer tx.Rollback(ctx)

		// TODO: use transaction
		rows, err := m.dbConnPool.Query(
			ctx, `
			SELECT msg
			FROM $1
			WHERE processed = FALSE`,
			table,
		)
		if err != nil {
			tx.Rollback(ctx) // TODO: log all rollbacks
			continue
		}

		for rows.Next() {
			var encodedMsg []byte
			rows.Scan(&encodedMsg)

			var msg kafka.Message
			_ = json.Unmarshal(encodedMsg, &msg)

			messages = append(messages, msg)
		}

		err = m.kafkaConnector.Writers[topic].WriteMessages(ctx, messages...)
		if err != nil {
			tx.Rollback(ctx) // TODO: log all rollbacks
			continue
		}

		tx.Exec(ctx, `
			UPDATE $1
			SET processed = TRUE
			WHERE id IN $2`,
			table, // TODO: fetch ids
		)		


		tx.Commit(ctx)
	}

}

func (m *TransactionalOutboxManager) Enqueue(
	ctx context.Context,
	tx pgx.Tx,
	topic consts.TopicName,
	msg kafka.Message,
) {
	encodedMsg, _ := json.Marshal(msg)

	// TODO: add to migrations
	tx.Exec(
		ctx, `
        INSERT INTO $1 (processed, msg)
		VALUES ($2, $3)`,
		getTableName(topic),
		false, encodedMsg,
	)

	select {
	case m.observer <- struct{}{}:
	default:
	}
}

func (m *TransactionalOutboxManager) Close() {
	m.cancelCtx()
}