package IrisAPIs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fetchRandomJoke(t *testing.T) {
	c, err := fetchRandomJoke(context.TODO())
	assert.Equal(t, err, nil)
	t.Logf("Fetched random joke : %+v", c)
}
