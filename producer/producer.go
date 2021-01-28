package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer interface {
	Publish(ctx context.Context, key, value []byte) error
	Close() error
}

type kafkaProducer struct {
	*kafka.Writer
	stop bool
}

func CreateKafkaProducer(brokerList []string, topic string) Producer {
	connectionTest(brokerList)
	var k = &kafkaProducer{
		Writer: &kafka.Writer{
			Addr: kafka.TCP(brokerList...),
			Topic: topic,
			Balancer: &kafka.LeastBytes{},
			RequiredAcks: 1,
			BatchSize: 10,
			BatchTimeout: time.Second*5,
		},
	}
	return k
}

func (k *kafkaProducer) Publish(ctx context.Context, key, value []byte) error {
	return k.WriteMessages(ctx, kafka.Message{Key: key, Value: value})
}

func (k *kafkaProducer) Close() error {
	if k.stop == false {
		k.stop = true
		return k.Writer.Close()
	}
	return nil
}

func connectionTest(brokerList []string) {
	for _, addr := range brokerList {
		conn, err := kafka.Dial("tcp", addr)
		if err != nil {
			panic(err.Error())
		}
		conn.Close()
	}
}

