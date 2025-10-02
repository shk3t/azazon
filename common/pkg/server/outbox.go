package server

import (
	"common/pkg/consts"
	"context"
	"encoding/json"
	"log"
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
	logger         *log.Logger
	trackedTopics  map[consts.TopicName]struct{}
}

func NewTransactionalOutboxManager(
	dbConnPool *pgxpool.Pool,
	kafkaConnector *KafkaConnector,
	logger *log.Logger,
) *TransactionalOutboxManager {
	m := &TransactionalOutboxManager{
		dbConnPool:     dbConnPool,
		kafkaConnector: kafkaConnector,
		observer:       make(chan struct{}, 1),
		logger:         logger,
		trackedTopics:  make(map[consts.TopicName]struct{}),
	}

	ticker := time.NewTicker(time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelCtx = cancel

	go func() {
		for {
			select {
			case <-m.observer:
				m.processUnprocessed()
			case <-ticker.C:
				m.processUnprocessed()
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

	for topic := range m.trackedTopics {
		table := getTableName(topic)
		messages := []kafka.Message{}
		msgIds := []int{}

		rows, err := m.dbConnPool.Query(
			ctx, `
			SELECT id, msg
			FROM $1
			WHERE processed = FALSE`,
			table,
		)
		if err != nil {
			m.logger.Println(err)
			continue
		}

		for rows.Next() {
			var id int
			var encoded []byte
			rows.Scan(&id, &encoded)

			var msg kafka.Message
			_ = json.Unmarshal(encoded, &msg)

			messages = append(messages, msg)
			msgIds = append(msgIds, id)
		}

		err = m.kafkaConnector.Writers[topic].WriteMessages(ctx, messages...)
		if err != nil {
			m.logger.Println(err)
			continue
		}

		m.dbConnPool.Exec(ctx, `
			UPDATE $1
			SET processed = TRUE, msg = (E''),
			WHERE id IN $2`,
			table,
			msgIds,
		)
	}

}

func (m *TransactionalOutboxManager) Enqueue(
	ctx context.Context,
	tx pgx.Tx,
	topic consts.TopicName,
	msg kafka.Message,
) {
	m.trackedTopics[topic] = struct{}{}
	encodedMsg, _ := json.Marshal(msg)

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