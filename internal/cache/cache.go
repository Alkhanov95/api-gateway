package cache

import (
	"context"
	"sync"
	"time"

	"github.com/alkhanov95/api-gateway/internal/models"
	"github.com/alkhanov95/api-gateway/internal/repository"
)

// cacheEntry = user + when it expires
type cacheEntry struct {
	user      *models.User
	expiresAt time.Time
}

// Decorator wraps repo and adds in-memory cache with TTL
type Decorator struct {
	repo  repository.UserProvider
	mu    sync.Mutex
	users map[string]cacheEntry
	ttl   time.Duration
}

// New creates a new cache wrapper with given TTL
func New(repo repository.UserProvider, ttl time.Duration) *Decorator {
	return &Decorator{
		repo:  repo,
		users: make(map[string]cacheEntry),
		ttl:   ttl,
	}
}

// CreateUser -> repo, then put in cache with TTL
func (d *Decorator) CreateUser(ctx context.Context, u *models.User) (string, error) {
	id, err := d.repo.CreateUser(ctx, u)
	if err != nil {
		return "", err
	}
	if id == "" || u == nil {
		return id, nil
	}

	u.ID = id
	d.mu.Lock()
	d.users[id] = cacheEntry{
		user:      u,
		expiresAt: time.Now().Add(d.ttl),
	}
	d.mu.Unlock()

	return id, nil
}

// GetUserByID -> check cache first. If missing/expired -> repo, then update cache
func (d *Decorator) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	now := time.Now()

	d.mu.Lock()
	if e, ok := d.users[id]; ok {
		if now.Before(e.expiresAt) {
			u := e.user
			d.mu.Unlock()
			return u, nil
		}
		// expired -> delete
		delete(d.users, id)
	}
	d.mu.Unlock()

	u, err := d.repo.GetUserByID(ctx, id)
	if err != nil || u == nil {
		return u, err
	}

	d.mu.Lock()
	d.users[id] = cacheEntry{
		user:      u,
		expiresAt: now.Add(d.ttl),
	}
	d.mu.Unlock()

	return u, nil
}

// List -> no cache, just proxy
func (d *Decorator) List(ctx context.Context) ([]models.User, error) {
	return d.repo.List(ctx)
}

// Update -> repo, then refresh cache
func (d *Decorator) Update(ctx context.Context, u *models.User) error {
	if err := d.repo.Update(ctx, u); err != nil {
		return err
	}
	if u != nil && u.ID != "" {
		d.mu.Lock()
		d.users[u.ID] = cacheEntry{
			user:      u,
			expiresAt: time.Now().Add(d.ttl),
		}
		d.mu.Unlock()
	}
	return nil
}

// Delete -> repo, then remove from cache
func (d *Decorator) Delete(ctx context.Context, id string) error {
	if err := d.repo.Delete(ctx, id); err != nil {
		return err
	}
	d.mu.Lock()
	delete(d.users, id)
	d.mu.Unlock()
	return nil
}
