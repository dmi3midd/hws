package hws_test

import (
	"context"
	"testing"

	"github.com/dmi3midd/hws"
	"github.com/dmi3midd/shkvcache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveRandomAndFind(t *testing.T) {
	ctx := context.Background()
	cache, _ := shkvcache.NewCache[string](ctx, &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 60,
		RunCleaner:      false,
	})
	defer cache.Close()
	service := hws.NewURLService(cache)

	testUrl := "https://example.com"

	// SaveRandom
	alias, err := service.SaveRandom(ctx, testUrl)
	require.NoError(t, err)
	assert.NotEmpty(t, alias)

	// Find
	url, err := service.Find(ctx, alias)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	assert.Equal(t, testUrl, url)
}

func TestSaveCustomAndFind(t *testing.T) {
	ctx := context.Background()
	cache, _ := shkvcache.NewCache[string](ctx, &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 60,
		RunCleaner:      false,
	})
	defer cache.Close()
	service := hws.NewURLService(cache)

	testUrl := "https://example.com"
	testAlias := "custom_alias"

	// SaveCustom
	alias, err := service.SaveCustom(ctx, testUrl, testAlias)
	require.NoError(t, err)
	assert.Equal(t, alias, testAlias)

	// Find
	url, err := service.Find(ctx, alias)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	assert.Equal(t, testUrl, url)
}

func TestFindNotFound(t *testing.T) {
	ctx := context.Background()
	cache, _ := shkvcache.NewCache[string](ctx, &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 60,
		RunCleaner:      false,
	})
	defer cache.Close()
	service := hws.NewURLService(cache)

	url, err := service.Find(ctx, "some_random_alias")
	require.Error(t, err)
	assert.Empty(t, url)
}

func TestSaveCustomWithExistingAlias(t *testing.T) {
	ctx := context.Background()
	cache, _ := shkvcache.NewCache[string](ctx, &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 60,
		RunCleaner:      false,
	})
	defer cache.Close()
	service := hws.NewURLService(cache)

	testUrl := "https://example.com"
	testAlias := "custom_alias"

	// SaveCustom 1
	alias, err := service.SaveCustom(ctx, testUrl, testAlias)
	require.NoError(t, err)
	assert.Equal(t, alias, testAlias)

	// SaveCustom 2
	alias, err = service.SaveCustom(ctx, testUrl, testAlias)
	require.Error(t, err)
	assert.Empty(t, alias)

	// Find
	url, err := service.Find(ctx, testAlias)
	require.NoError(t, err)
	assert.NotEmpty(t, url)
}
