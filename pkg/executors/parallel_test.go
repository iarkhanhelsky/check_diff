package executors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParallel(t *testing.T) {
	executor := NewParallel()
	results := make(chan int)

	done := []bool{false, false}
	executor.Add(func() error {
		done[0] = true
		return nil
	})
	executor.Add(func() error {
		done[1] = true
		return nil
	})

	err := executor.Run()
	close(results)
	assert.NoError(t, err)
	assert.Equal(t, []bool{true, true}, done)
}
