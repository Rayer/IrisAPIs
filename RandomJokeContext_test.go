package IrisAPIs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fetchRandomJoke(t *testing.T) {
	c, err := fetchRandomJoke()
	assert.Equal(t, err, nil)
	fmt.Printf("Fetched random joke : %+v", c)
}
