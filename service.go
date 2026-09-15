package hws

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/shkvcache"
	"github.com/rs/xid"
)

var (
	ErrUrlNotFound           = errors.New("url not found")
	ErrFailedToGenerateAlias = errors.New("failed to generate unique alias")
	ErrAliasAlreadyExists    = errors.New("alias already exists")
)

type URLService interface {
	// Find returns a Url by its alias.
	// It returns [ErrUrlNotFound] if no url are found.
	Find(ctx context.Context, alias string) (string, error)
	// SaveRandom saves a random alias for a given url and returns it.
	// It returns [ErrFailedToGenerateAlias] if no unique alias could be generated.
	SaveRandom(ctx context.Context, url string) (string, error)
	// SaveCustom saves a custom alias for a given url. Returns the alias.
	// It returns [ErrAliasAlreadyExists] if the alias already exists.
	SaveCustom(ctx context.Context, url string, alias string) (string, error)
}

type urlService struct {
	cache shkvcache.Cache[string]
}

func NewURLService(cache shkvcache.Cache[string]) URLService {
	return &urlService{
		cache: cache,
	}
}

func (s *urlService) Find(ctx context.Context, alias string) (string, error) {
	op := "URLService.Get"
	url, ok := s.cache.Get(alias)
	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrUrlNotFound)
	}
	return url, nil
}

func (s *urlService) SaveRandom(ctx context.Context, urlStr string) (string, error) {
	op := "URLService.Save"
	alias, err := s.generateUniqueAlias(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	s.cache.Set(alias, urlStr, 0)
	return alias, nil
}

func (s *urlService) SaveCustom(ctx context.Context, urlStr string, alias string) (string, error) {
	op := "URLService.Save"
	_, ok := s.cache.Get(alias)
	if ok {
		return "", fmt.Errorf("%s: %w", op, ErrAliasAlreadyExists)
	}
	s.cache.Set(alias, urlStr, 0)
	return alias, nil
}

func (s *urlService) generateUniqueAlias(ctx context.Context) (string, error) {
	const op = "URLService.generateUniqueAlias"
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		alias := xid.New().String()[:8]
		_, ok := s.cache.Get(alias)
		if !ok {
			return alias, nil
		}
	}
	return "", fmt.Errorf("%s: %w", op, ErrFailedToGenerateAlias)
}
