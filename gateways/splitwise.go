package gateways

import (
	"context"
	"fmt"
	"log"

	"github.com/aanzolaavila/splitwise.go"
	"github.com/aanzolaavila/splitwise.go/resources"
)

type splitwiseResouces interface {
	resources.Expense |
		resources.MainCategory |
		resources.Currency |
		resources.Friend |
		resources.Group |
		resources.Notification
}

type swConnectionStruct struct {
	ctx    context.Context
	client splitwise.Client
}

type SwConnection interface {
	GetMainCategories() CommandExecutor[resources.MainCategory]
	GetMainCategory(id resources.Identifier) (resources.MainCategory, error)
	GetCurecies() CommandExecutor[resources.Currency]
	GetCurency(code string) (resources.Currency, error)
	GetFriends() CommandExecutor[resources.Friend]
	GetFriend(id int) (resources.Friend, error)
	GetGroups() CommandExecutor[resources.Group]
	GetGroup(id int) (resources.Group, error)
	GetNotifications(params splitwise.NotificationsParams) CommandExecutor[resources.Notification]
	GetExpense(id int) (resources.Expense, error)
	GetExpenses(params splitwise.ExpensesParams) CommandExecutor[resources.Expense]
	getClient() splitwise.Client
}

type commandExecutorStruct[T splitwiseResouces] struct {
	swConnectionStruct
	ch    chan T
	close bool
}

type CommandExecutor[T splitwiseResouces] interface {
	isClose() bool
	Close()
	GetChan() <-chan T
}

type ElementNotFound struct{}

func (m *ElementNotFound) Error() string {
	return "Element Not Found"
}

func getTokenClient(token string) splitwise.Client {
	return splitwise.Client{
		Token: token,
	}
}

var mainCategoryCache map[resources.Identifier]resources.MainCategory = make(map[resources.Identifier]resources.MainCategory)
var curenciesCache map[string]resources.Currency = make(map[string]resources.Currency)

func Open(token string, ctx context.Context, log *log.Logger) SwConnection {

	conn := swConnectionStruct{}

	conn.client = getTokenClient(token)
	conn.client.Logger = log
	conn.ctx = ctx
	return conn
}

func (ce *commandExecutorStruct[T]) isClose() bool {
	return ce.close
}

func (ce *commandExecutorStruct[T]) Close() {
	ce.close = true
	for range ce.ch {
	}
}

func (ce *commandExecutorStruct[T]) GetChan() <-chan T {
	return ce.ch
}

func simpleExecutor[T splitwiseResouces](conn swConnectionStruct, method func(ctx context.Context) ([]T, error)) CommandExecutor[T] {
	ch := make(chan T)
	ce := commandExecutorStruct[T]{}
	ce.ch = ch
	ce.swConnectionStruct = conn
	ce.close = false

	go func(ch chan<- T) {
		defer ce.cleanCe()
		defer recoverClosedChannel()
		entities, err := method(ce.ctx)

		if err != nil {
			ce.getClient().Logger.Printf(err.Error())
			return
		}

		for _, e := range entities {
			if ce.isClose() {
				return
			}
			ch <- e
		}

	}(ch)
	return &ce

}

func (conn swConnectionStruct) GetMainCategories() CommandExecutor[resources.MainCategory] {
	if len(mainCategoryCache) == 0 {
		client := conn.getClient()
		ce := simpleExecutor(conn, client.GetCategories)
		for v := range ce.GetChan() {
			mainCategoryCache[resources.Identifier(v.ID)] = v
		}
		return conn.GetMainCategories()
	}
	return conn.retriveCachedMainCategor()
}

func (conn swConnectionStruct) GetMainCategory(id resources.Identifier) (resources.MainCategory, error) {
	if len(mainCategoryCache) == 0 {
		tmp := conn.GetMainCategories()
		for range tmp.GetChan() {
		}
	}
	result, ok := mainCategoryCache[id]
	if !ok {
		return resources.MainCategory{}, &ElementNotFound{}
	}

	return result, nil
}

func (conn swConnectionStruct) retriveCachedMainCategor() CommandExecutor[resources.MainCategory] {
	ch := make(chan resources.MainCategory)
	ce := commandExecutorStruct[resources.MainCategory]{}
	ce.ch = ch
	ce.swConnectionStruct = conn
	ce.close = false

	go func(ch chan<- resources.MainCategory) {
		defer ce.cleanCe()
		defer recoverClosedChannel()
		for _, v := range mainCategoryCache {
			ch <- v
		}
	}(ch)
	return &ce
}

func (conn swConnectionStruct) GetCurecies() CommandExecutor[resources.Currency] {
	if len(curenciesCache) == 0 {
		client := conn.getClient()
		ce := simpleExecutor(conn, client.GetCurrencies)
		for v := range ce.GetChan() {
			curenciesCache[v.CurrencyCode] = v
		}
		return conn.GetCurecies()
	}
	return conn.retriveCachedCurrency()
}

func (conn swConnectionStruct) GetCurency(code string) (resources.Currency, error) {
	if len(curenciesCache) == 0 {
		tmp := conn.GetCurecies()
		for range tmp.GetChan() {
		}
	}
	result, ok := curenciesCache[code]
	if !ok {
		return resources.Currency{}, &ElementNotFound{}
	}

	return result, nil
}

func (conn swConnectionStruct) retriveCachedCurrency() CommandExecutor[resources.Currency] {
	ch := make(chan resources.Currency)
	ce := commandExecutorStruct[resources.Currency]{}
	ce.ch = ch
	ce.swConnectionStruct = conn
	ce.close = false

	go func(ch chan<- resources.Currency) {
		defer ce.cleanCe()
		defer recoverClosedChannel()
		for _, v := range curenciesCache {
			ch <- v
		}
	}(ch)
	return &ce

}

func (conn swConnectionStruct) GetFriends() CommandExecutor[resources.Friend] {
	client := conn.getClient()
	return simpleExecutor(conn, client.GetFriends)
}

func (conn swConnectionStruct) GetFriend(id int) (resources.Friend, error) {
	client := conn.getClient()
	return client.GetFriend(conn.ctx, id)
}

func (conn swConnectionStruct) GetGroups() CommandExecutor[resources.Group] {
	client := conn.getClient()

	return simpleExecutor(conn, client.GetGroups)
}

func (conn swConnectionStruct) GetGroup(id int) (resources.Group, error) {
	client := conn.getClient()

	return client.GetGroup(conn.ctx, id)
}

func (conn swConnectionStruct) GetNotifications(params splitwise.NotificationsParams) CommandExecutor[resources.Notification] {
	ch := make(chan resources.Notification)
	ce := commandExecutorStruct[resources.Notification]{}
	ce.ch = ch
	ce.swConnectionStruct = conn
	ce.close = false

	go func(ch chan<- resources.Notification) {
		defer ce.cleanCe()
		defer recoverClosedChannel()
		client := conn.getClient()

		notifications, err := client.GetNotifications(ce.ctx, params)
		if err != nil {
			ce.client.Logger.Printf(err.Error())
			return
		}

		for _, e := range notifications {
			if ce.isClose() {
				return
			}
			ch <- e
		}

	}(ch)

	return &ce
}

func (conn swConnectionStruct) GetExpense(id int) (resources.Expense, error) {
	client := conn.getClient()

	return client.GetExpense(conn.ctx, id)
}

func (conn swConnectionStruct) GetExpenses(params splitwise.ExpensesParams) CommandExecutor[resources.Expense] {
	ch := make(chan resources.Expense)
	ce := commandExecutorStruct[resources.Expense]{}

	ce.ch = ch
	ce.swConnectionStruct = conn
	ce.close = false

	go func(ch chan<- resources.Expense) {
		defer ce.cleanCe()
		defer recoverClosedChannel()
		client := conn.getClient()

		var (
			cont int
		)
		for !ce.close {
			expenses, err := client.GetExpenses(ce.ctx, params)
			if err != nil {
				ce.client.Logger.Printf(err.Error())
				break
			}

			if len(expenses) == 0 {
				break
			}

			for _, e := range expenses {
				if ce.isClose() {
					return
				}
				ch <- e
				cont++
			}
			incOffset(params, cont)
		}
		fmt.Println("Exit background function")
	}(ch)

	return &ce
}

func (conn swConnectionStruct) getClient() splitwise.Client {
	return conn.client
}

func incOffset(params splitwise.ExpensesParams, inc int) {
	params[splitwise.ExpensesOffset] = inc
}

func recoverClosedChannel() {
	// recover from panic caused by writing to a closed channel
	if r := recover(); r != nil {
		err := fmt.Errorf("%v", r)
		fmt.Printf("write: error writing on channel: %v\n", err)
		return
	}
}

func (ce *commandExecutorStruct[T]) cleanCe() {
	ce.close = true
	close(ce.ch)
	fmt.Println("channel closed succesfuly")
}
