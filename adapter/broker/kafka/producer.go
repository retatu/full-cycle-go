package kafka

import (
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/retatu/fullcycle-gateway/adapter/presenter"
)

type Producer struct {
	ConfigMap *ckafka.ConfigMap
	Presenter presenter.Presenter
}

func NewKafkaProducer(configMap *ckafka.ConfigMap, presenter presenter.Presenter) *Producer {
	return &Producer{ConfigMap: configMap, Presenter: presenter}
}

func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	producer, err := ckafka.NewProducer(p.ConfigMap)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = p.Presenter.Bind(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	presenterMsg, err := p.Presenter.Show()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          presenterMsg,
		Key:            key,
	}
	err = producer.Produce(message, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
