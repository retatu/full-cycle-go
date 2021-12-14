package transaction

import (
	"encoding/json"

	"github.com/retatu/fullcycle-gateway/usecase/process_transaction"
)

type KafkaPresenter struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

func NewTransactionKafkaPresenter() *KafkaPresenter {
	return &KafkaPresenter{}
}

func (t *KafkaPresenter) Bind(output interface{}) error {
	t.ID = output.(process_transaction.TransactionDtoOutput).ID
	t.Status = output.(process_transaction.TransactionDtoOutput).Status
	t.ErrorMessage = output.(process_transaction.TransactionDtoOutput).ErrorMessage
	return nil
}

func (t *KafkaPresenter) Show() ([]byte, error) {
	j, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return j, err
}
