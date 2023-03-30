package gateways

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
)

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
