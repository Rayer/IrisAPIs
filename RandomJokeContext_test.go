package IrisAPIs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fetchRandomJoke(t *testing.T) {
	c, err := fetchRandomJoke()
	assert.Equal(t, err, nil)
	t.Logf("Fetched random joke : %+v", c)
}
