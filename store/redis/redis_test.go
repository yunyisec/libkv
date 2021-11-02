package redis

import (
	"context"
	"testing"

	"github.com/yunyisec/libkv"
	"github.com/yunyisec/libkv/store"
	"github.com/yunyisec/libkv/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	client = "localhost:6379"
)

func makeRedisClient(t *testing.T) store.Store {
	kv, err := newRedis([]string{client}, "", 0)
	if err != nil {
		t.Fatalf("cannot create store: %v", err)
	}

	// NOTE: please turn on redis's notification
	// before you using watch/watchtree/lock related features
	kv.client.ConfigSet(context.Background(), "notify-keyspace-events", "KA")

	return kv
}

func TestRegister(t *testing.T) {
	Register()

	kv, err := libkv.NewStore(store.REDIS, []string{client}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	if _, ok := kv.(*Redis); !ok {
		t.Fatal("Error registering and initializing redis")
	}
}

func TestRedisStore(t *testing.T) {
	kv := makeRedisClient(t)
	lockTTL := makeRedisClient(t)
	kvTTL := makeRedisClient(t)

	testutils.RunTestCommon(t, kv)
	testutils.RunTestAtomic(t, kv)
	testutils.RunTestWatch(t, kv)
	testutils.RunTestLock(t, kv)
	testutils.RunTestLockTTL(t, kv, lockTTL)
	testutils.RunTestTTL(t, kv, kvTTL)
	testutils.RunCleanup(t, kv)
}
