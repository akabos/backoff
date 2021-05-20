package backoff

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBackoff(t *testing.T) {
	b := Backoff{
		Initial: time.Second,
		Final:   time.Second * 5,
		Scale:   1.75,
	}
	var r []time.Duration
	for i := 0; i < 5; i++ {
		r = append(r, b.Next())
	}
	b.Reset()
	r = append(r, b.Next())
	require.Equal(t, []time.Duration{
		time.Millisecond * 1000,
		time.Millisecond * 1750,
		time.Millisecond * 3062 + time.Microsecond * 500,
		time.Millisecond * 5000,
		time.Millisecond * 5000,
		time.Millisecond * 1000,
	}, r)
}

func TestBackoffDefaults(t *testing.T) {
	b := Backoff{}
	var r []time.Duration
	for i := 0; i < 8; i++ {
		r = append(r, b.Next())
	}
	b.Reset()
	r = append(r, b.Next())
	require.Equal(t, []time.Duration{
		time.Second * 1,
		time.Millisecond * 1500,
		time.Millisecond * 2250,
		time.Millisecond * 3375,
		time.Millisecond * 5062 + time.Microsecond * 500,
		time.Millisecond * 7593 + time.Microsecond * 750,
		time.Second * 10,
		time.Second * 10,
		time.Second * 1,
	}, r)
}

func TestSleep(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		b := Backoff{Initial: time.Millisecond}
		err := b.Sleep(context.Background())
		require.NoError(t, err)
	})
	t.Run("cancelled", func(t *testing.T) {
		b := Backoff{}
		ctx, cancel := context.WithCancel(context.Background())
		go cancel()
		err := b.Sleep(ctx)
		require.Error(t, err)
		require.True(t, errors.Is(err, context.Canceled))
	})
}