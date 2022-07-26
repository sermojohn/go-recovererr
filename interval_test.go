package recovererr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxRetrier(t *testing.T) {
	mt := NewMaxTicker(5, 1)

	var counter int
	for i := 0; i < 100; i++ {
		t := <-mt.c
		if t != (time.Time{}) {
			counter++
		}

	}

	assert.Equal(t, 5, counter)
	assert.Equal(t, 0, len(mt.c))
}
