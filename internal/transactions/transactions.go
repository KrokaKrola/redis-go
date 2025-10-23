package transactions

import (
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
)

type Transactions struct {
	innerMap transactionsMap
	sync.RWMutex
}

type commandsList []*commands.Command

type transactionsMap map[string]commandsList

func NewTransactions() *Transactions {
	return &Transactions{
		innerMap: make(transactionsMap),
	}
}

func (t *Transactions) CleanupTransactionById(id string) {
	t.Lock()
	defer t.Unlock()

	delete(t.innerMap, id)
}

func (t *Transactions) GetTransactionById(id string) (commandsList, bool) {
	t.Lock()
	defer t.Unlock()

	list, ok := t.innerMap[id]

	return list, ok
}

func (t *Transactions) UpdateTransactionById(id string, list commandsList) {
	t.Lock()
	defer t.Unlock()

	_, ok := t.innerMap[id]
	if ok {
		t.innerMap[id] = list
	}
}

func (t *Transactions) NewTransactionsListById(id string) {
	t.Lock()
	defer t.Unlock()
	t.innerMap[id] = []*commands.Command{}
}
