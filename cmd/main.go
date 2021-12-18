package main

import (
	"database/sql"
	"encoding/json"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/mattn/go-sqlite3"
	"github.com/retatu/fullcycle-gateway/adapter/broker/kafka"
	"github.com/retatu/fullcycle-gateway/adapter/factory"
	"github.com/retatu/fullcycle-gateway/adapter/presenter/transaction"
	"github.com/retatu/fullcycle-gateway/usecase/process_transaction"
)

func main() {
	//db
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	//repository
	repositoryFactory := factory.NewRepositoryDatabaseFactory(db)
	repository := repositoryFactory.CreateTransactionRepository()
	//confiMapProducer
	configMapProducer := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:172.17.0.1",
	}
	kafkaPresenter := transaction.NewTransactionKafkaPresenter()
	//producer
	producer := kafka.NewKafkaProducer(configMapProducer, kafkaPresenter)
	//confiMapConsumer
	configMapConsumer := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:172.17.0.1",
		"client.id":         "goapp",
		"group.id":          "goapp",
	}
	//topic
	topics := []string{"transactions"}
	//consumer
	var msgChan = make(chan *ckafka.Message)
	consumer := kafka.NewConsumer(configMapConsumer, topics)
	go consumer.Consume(msgChan)

	//usecase
	usecase := process_transaction.NewProcessTransaction(repository, producer, "transactions_result")
	for msg := range msgChan {
		var input process_transaction.TransactionDtoInput
		json.Unmarshal(msg.Value, &input)
		usecase.Execute(input)
	}
}
