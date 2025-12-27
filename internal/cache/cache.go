package cache

import (
	"context"
	"sync"
	"time"

	"github.com/alkhanov95/api-gateway/internal/models"
	"github.com/alkhanov95/api-gateway/internal/repository"
)

type wrapUser struct {
	user      *models.User
	expiresAt time.Time
}

// Decorator wraps repo and adds in-memory cache with TTL
type Decorator struct {
	repo  repository.UserProvider
	mu    sync.RWMutex
	users map[string]wrapUser
	ttl   time.Duration // add method that works in background (goroutine) once in 30 secs -> deletes old data from cache
	// old data -> the data that expires at  < time.Now then we delete it
}

// New creates a new cache wrapper with given TTL
func New(repo repository.UserProvider, ttl time.Duration) *Decorator {
	return &Decorator{
		repo:  repo,
		users: make(map[string]wrapUser),
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
	d.set(u)
	return id, nil
}

func (d *Decorator) get(id string) *models.User {
	d.mu.RLock()
	defer d.mu.RUnlock()
	e, ok := d.users[id]
	if !ok {
		return nil // not in cache
	}
	// if TTL expired -> nothing to return
	if time.Now().After(e.expiresAt) {
		return nil
	}
	return e.user
}
func (d *Decorator) set(user *models.User) {
	if user == nil || user.ID == "" || d.ttl <= 0 {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	d.users[user.ID] = wrapUser{
		user:      user,
		expiresAt: time.Now().Add(d.ttl),
	}
}

func (d *Decorator) StartGC() {
	if d.ttl <= 0 {
		return
	}
	go func() {
		for {
			time.Sleep(30 * time.Second)
			now := time.Now()
			d.mu.Lock()
			for k, v := range d.users {
				if now.After(v.expiresAt) {
					delete(d.users, k)
				}
			}
			d.mu.Unlock()
		}
	}()
}

// GetUserByID -> check cache first. If missing/expired -> repo, then update cache
func (d *Decorator) GetUserByID(ctx context.Context, id string) (*models.User, error) {

	if user := d.get(id); user != nil {
		return user, nil
	}

	u, err := d.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	d.set(u)
	return u, nil
}

// List -> no cache, just proxy
func (d *Decorator) List(ctx context.Context) ([]models.User, error) {
	return d.repo.List(ctx)
}

func (d *Decorator) Update(ctx context.Context, u *models.User) error {
	if u == nil || u.ID == "" {
		return nil
	}
	// 1) updating db
	if err := d.repo.Update(ctx, u); err != nil {
		return err
	}
	// 2) updating cache
	d.set(u)
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
