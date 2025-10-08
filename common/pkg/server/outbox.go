package server

import (
	"common/pkg/consts"
	"context"
	"encoding/json"
	"fmt"
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

const tablePostfix = "_outbox"

func getTableName(topic consts.TopicName) string {
	return string(topic) + tablePostfix
}

func (m *TransactionalOutboxManager) processUnprocessed() {
	ctx := context.Background()

	for topic := range m.trackedTopics {
		table := getTableName(topic)
		messages := []kafka.Message{}
		msgIds := []int{}

		rows, err := m.dbConnPool.Query(
			ctx,
			fmt.Sprintf(`
				SELECT id, msg
				FROM %s
				WHERE processed = FALSE`,
				table,
			),
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

		m.dbConnPool.Exec(ctx,
			fmt.Sprintf(`
				UPDATE %s
				SET processed = TRUE, msg = (E'')
				WHERE id = ANY ($1)`,
				table,
			),
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
		ctx,
		fmt.Sprintf(`
			INSERT INTO %s (processed, msg)
			VALUES (FALSE, $1)`,
			getTableName(topic),
		),
		encodedMsg,
	)

	select {
	case m.observer <- struct{}{}:
	default:
	}
}

func (m *TransactionalOutboxManager) Close() {
	m.cancelCtx()
}