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
