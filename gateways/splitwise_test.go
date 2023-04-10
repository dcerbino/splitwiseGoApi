package gateways

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
)

const getFriends200Response = `
{
  "friends": [
    {
      "id": 15,
      "first_name": "Ada",
      "last_name": "Lovelace",
      "email": "ada@example.com",
      "registration_status": "confirmed",
      "picture": {
        "small": "string",
        "medium": "string",
        "large": "string"
      },
      "groups": [
        {
          "group_id": 571,
          "balance": [
            {
              "currency_code": "USD",
              "amount": "414.5"
            }
          ]
        }
      ],
      "balance": [
        {
          "currency_code": "USD",
          "amount": "414.5"
        }
      ],
      "updated_at": "2019-08-24T14:15:22Z"
    },
	{
		"id": 16,
		"first_name": "Pepe",
		"last_name": "ponce",
		"email": "pepe@example.com",
		"registration_status": "confirmed",
		"picture": {
		  "small": "string",
		  "medium": "string",
		  "large": "string"
		},
		"groups": [
		  {
			"group_id": 571,
			"balance": [
			  {
				"currency_code": "USD",
				"amount": "329.5"
			  }
			]
		  }
		],
		"balance": [
		  {
			"currency_code": "USD",
			"amount": "329.5"
		  }
		],
		"updated_at": "2019-08-23T14:15:22Z"
	  }
  ]
}
`

type httpClientStub struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (c httpClientStub) Do(r *http.Request) (*http.Response, error) {
	if c.DoFunc == nil {
		panic("mocked function is nil")
	}

	return c.DoFunc(r)
}

type logger interface {
	Printf(string, ...interface{})
}

type testLogger struct {
	buf    bytes.Buffer
	logger logger
	once   sync.Once
	T      *testing.T
}

func (l testLogger) Printf(s string, args ...interface{}) {
	l.once.Do(func() {
		tname := l.T.Name()
		prefix := fmt.Sprintf("%s:: ", tname)
		l.logger = log.New(io.Writer(&l.buf), prefix, log.LstdFlags)

		l.T.Cleanup(func() {
			if l.T.Failed() {
				fmt.Print(l.buf.String())
			}
		})
	})

	l.logger.Printf(s, args...)
}

func TestOpen(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	result := Open(token, ctx, log)

	assert.Equal(token, result.getClient().Token)
	assert.Equal(ctx, result.getCtx())
	assert.Equal(log, result.getClient().Logger)
}

func TestMainCategoryCache(t *testing.T) {
	assert.Equal(t, true, len(mainCategoryCache) >= 7)
}

func TestCurenciesCache(t *testing.T) {
	assert.Equal(t, true, len(curenciesCache) >= 148)
}

func TestClose(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	executor := conn.GetCurecies()

	assert.Equal(false, executor.isClose())
	executor.Close()
	assert.Equal(true, executor.isClose())
}

func TestGetCategory(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	category, err := conn.GetMainCategory(resources.Identifier(1))
	assert.Equal(nil, err)
	assert.Equal("Utilities", category.Name)

}

func TestGetCategoryNotFound(t *testing.T) {
	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	_, err := conn.GetMainCategory(resources.Identifier(0))
	assert.EqualErrorf(t, err, (&ElementNotFound{}).Error(), "Error should be: %v, got: %v", (&ElementNotFound{}).Error(), err)

}

func TestGetCategoryies(t *testing.T) {
	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	executor := conn.GetMainCategories()
	cont := 0
	want := 7

	for range executor.GetChan() {
		cont++
	}

	assert.GreaterOrEqual(t, cont, want, "Get categories should return unless %d arg and got it %d", want, cont)
}

func TestGetCurrencies(t *testing.T) {
	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	executor := conn.GetCurecies()
	cont := 0
	want := 148

	for range executor.GetChan() {
		cont++
	}

	assert.GreaterOrEqual(t, cont, want, "Get currencies should return unless %d currencies and got it %d", want, cont)
}

func TestGetCurrency(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	currencyCode := "USD"
	currencyUnit := "$"

	currency, err := conn.GetCurency(currencyCode)

	assert.NoError(err, "%s should be present as a currency code", currencyCode)
	assert.Equal(currencyUnit, currency.Unit, "%s currency unit should be %s but got %s", currencyCode, currencyUnit, currency.Unit)
}

func TestGetCurrencyNotFund(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	currencyCode := "US"

	_, err := conn.GetCurency(currencyCode)

	assert.EqualErrorf(err, (&ElementNotFound{}).Error(), "Error should be: %v, got: %v", (&ElementNotFound{}).Error(), err)
}

func TestGetFriends(t *testing.T) {
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(getFriends200Response))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}

	type responseStruct struct {
		Friends []resources.Friend
	}
	wantedRespounce := responseStruct{}

	err := json.Unmarshal([]byte(getFriends200Response), &wantedRespounce)

	if err != nil {
		panic(err)
	}

	conn := getTestConnection(t, doFunc)

	executor := conn.GetFriends()

	cont := 0

	for range executor.GetChan() {
		cont++
	}

	assert.Equal(t, len(wantedRespounce.Friends), cont)
}

func getTestConnection(t *testing.T, doFunc func(r *http.Request) (*http.Response, error)) SwConnection {
	client := Open("testtoken", context.Background(), log.New(os.Stdout, "Test Splitwise LOG: ", log.Lshortfile))

	bareclient, ok := client.(*swConnectionStruct)

	if !ok {
		panic("unable to convert interface to istance")
	}

	cliststub := httpClientStub{
		DoFunc: doFunc,
	}

	bareclient.client.HttpClient = cliststub
	bareclient.client.Logger = testLogger{
		T: t,
	}

	return bareclient
}
