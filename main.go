package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aanzolaavila/splitwise.go"
	"github.com/dcerbino/splitwiseGoApi/gateways"
)

func main() {
	argsWithoutProg := os.Args[1:]
	logger := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)
	ctx := context.Background()

	conn := gateways.Open(argsWithoutProg[0], ctx, logger)

	params := splitwise.ExpensesParams{}
	params[splitwise.ExpensesDatedAfter] = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
	params[splitwise.ExpensesGroupId] = argsWithoutProg[1]
	params[splitwise.ExpensesLimit] = 200

	ce := conn.GetExpenses(params)

	cont := 0

	for e := range ce.GetChan() {
		if !e.Payment {
			fmt.Printf("record: %s , %s, %s ,%v\n", e.Date, e.Category.Name, e.Description, e.Cost)
		}
		cont++
	}

	fmt.Printf("number of Expenses: %d\n", cont)

	cont = 0
	ce2 := conn.GetMainCategories()
	for e := range ce2.GetChan() {
		fmt.Println("Category1: ", e)
		cont++
	}

	fmt.Printf("number of Category: %d\n", cont)

	cont = 0
	ce2 = conn.GetMainCategories()

	for e := range ce2.GetChan() {
		fmt.Println("Category2: ", e)
		cont++
	}

	fmt.Printf("number of Category: %d\n", cont)

	cont = 0
	ce3 := conn.GetCurecies()
	for e := range ce3.GetChan() {
		fmt.Println("Currency: ", e)
		cont++
	}

	fmt.Printf("number of Curecies: %d\n", cont)

	cont = 0
	ce4 := conn.GetFriends()
	for e := range ce4.GetChan() {
		fmt.Println("Friend: ", e)
		cont++
	}

	fmt.Printf("number of Friends: %d\n", cont)

	cont = 0
	ce5 := conn.GetFriends()
	for e := range ce5.GetChan() {
		fmt.Println("Group: ", e)
		cont++
	}

	fmt.Printf("number of friends: %d\n", cont)

	param := splitwise.NotificationsParams{}
	param[splitwise.NotificationsLimit] = 1000

	cont = 0
	ce6 := conn.GetNotifications(param)
	for e := range ce6.GetChan() {
		fmt.Println("Group: ", e)
		cont++
	}

	fmt.Printf("number of notifications: %d\n", cont)
}
