package gateways

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestOpen(t *testing.T) {
	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	result := Open(token, ctx, log)

	assert.Equal(t, token, result.getClient().Token)
	assert.Equal(t, ctx, result.getCtx())
	assert.Equal(t, log, result.getClient().Logger)
}

func TestMainCategoryCache(t *testing.T) {
	assert.Equal(t, true, len(mainCategoryCache) >= 7)
}

func TestCurenciesCache(t *testing.T) {
	assert.Equal(t, true, len(curenciesCache) >= 148)
}

func TestClose(t *testing.T) {
	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	executor := conn.GetCurecies()

	assert.Equal(t, false, executor.isClose())
	executor.Close()
	assert.Equal(t, true, executor.isClose())
}
