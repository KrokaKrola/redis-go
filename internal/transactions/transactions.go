package transactions

import (
	"fmt"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type Transactions struct {
	innerMap transactionsMap
	sync.RWMutex
}

type commandsList []*commands.Command

type transactionsMap map[string]commandsList

type Executor func(*commands.Command) resp.Value

func NewTransactions() *Transactions {
	return &Transactions{
		innerMap: make(transactionsMap),
	}
}

func (t *Transactions) IsActive(id string) bool {
	t.RLock()
	defer t.RUnlock()

	_, ok := t.innerMap[id]

	return ok
}

func (t *Transactions) Begin(id string) {
	t.Lock()
	defer t.Unlock()
	t.innerMap[id] = []*commands.Command{}
}

func (t *Transactions) Queue(id string, cmd *commands.Command) error {
	t.Lock()
	defer t.Unlock()

	list, ok := t.innerMap[id]
	if !ok {
		return fmt.Errorf("ERR unknown session id")
	}

	list = append(list, cmd)
	t.innerMap[id] = list

	return nil
}

func (t *Transactions) ExecuteAndDiscard(id string, executor Executor) *resp.Array {
	t.Lock()
	defer t.Unlock()

	arr := &resp.Array{}

	list, ok := t.innerMap[id]

	if !ok {
		return arr
	}

	for _, command := range list {
		arr.Elements = append(arr.Elements, executor(command))
	}

	delete(t.innerMap, id)

	return arr
}

func (t *Transactions) Discard(id string) {
	t.Lock()
	defer t.Unlock()

	delete(t.innerMap, id)
}
