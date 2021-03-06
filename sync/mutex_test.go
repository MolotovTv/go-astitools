package astisync_test

import (
	"testing"
	"time"

	"github.com/molotovtv/go-astitools/sync"
	"github.com/stretchr/testify/assert"
)

func TestRWMutex_IsDeadlocked(t *testing.T) {
	var m = astisync.NewRWMutex("test", false)
	d, _ := m.IsDeadlocked(time.Millisecond)
	assert.False(t, d)
	m.Lock()
	d, c := m.IsDeadlocked(time.Millisecond)
	assert.True(t, d)
	assert.Contains(t, c, "/sync/mutex_test.go:15")
}
