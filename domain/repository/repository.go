package repository

type TransacationRepository interface {
	Insert(id string, accountId string, amount float64, status string, errorMessage string) error
}
