package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/retatu/fullcycle-gateway/adapter/broker/kafka"
	"github.com/retatu/fullcycle-gateway/adapter/factory"
	"github.com/retatu/fullcycle-gateway/adapter/presenter/transaction"
	"github.com/retatu/fullcycle-gateway/usecase/process_transaction"
)

func main() {
	fmt.Println("Running")
	//db
	db, err := sql.Open("mysql", os.Getenv("MYSQL_USERNAME")+":"+os.Getenv("MYSQL_PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+os.Getenv("MYSQL_DATABASE"))
	if err != nil {
		log.Fatal(err)
	}
	//repository
	repositoryFactory := factory.NewRepositoryDatabaseFactory(db)
	repository := repositoryFactory.CreateTransactionRepository()
	//confiMapProducer
	configMapProducer := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("SECURITY_PROTOCOL"),
		"sasl.mechanisms":   os.Getenv("SASL_MECHANISMS"),
		"sasl.username":     os.Getenv("SASL_USERNAME"),
		"sasl.password":     os.Getenv("SASL_PASSWORD"),
	}
	kafkaPresenter := transaction.NewTransactionKafkaPresenter()
	//producer
	producer := kafka.NewKafkaProducer(configMapProducer, kafkaPresenter)
	//confiMapConsumer
	configMapConsumer := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("SECURITY_PROTOCOL"),
		"sasl.mechanisms":   os.Getenv("SASL_MECHANISMS"),
		"sasl.username":     os.Getenv("SASL_USERNAME"),
		"sasl.password":     os.Getenv("SASL_PASSWORD"),
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
