package types

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropy     *ulid.MonotonicEntropy
	entropyOnce sync.Once
)

func NewID() string {
	entropyOnce.Do(func() {
		entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	})
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
