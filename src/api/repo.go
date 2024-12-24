package api

import (
	"sync"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Biography string    `json:"bio"`
}

type InMemoryStore struct {
	mu    sync.RWMutex
	items map[uuid.UUID]User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		items: make(map[uuid.UUID]User),
	}
}

func (db *InMemoryStore) GetUsers() []User {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]User, 0)

	for _, user := range db.items {
		users = append(users, user)
	}

	return users
}

func (db *InMemoryStore) GetUserById(id uuid.UUID) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	key := uuid.UUID(id)

	user, ok := db.items[key]
	if !ok {
		return user, ErrorNotFound
	}

	return user, nil
}

func (db *InMemoryStore) CreateUser(user User) {
	db.mu.Lock()
	defer db.mu.Unlock()

	id := uuid.UUID(uuid.New())

	db.items[id] = User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Biography: user.Biography,
		ID:        id,
	}
}

func (db *InMemoryStore) UpdateUser(user User, id uuid.UUID) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	key := uuid.UUID(id)
	_, ok := db.items[key]
	if !ok {
		return ErrorNotFound
	}

	db.items[key] = User{
		ID:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Biography: user.Biography,
	}

	return nil
}

func (db *InMemoryStore) DeleteUser(id uuid.UUID) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	key := uuid.UUID(id)

	_, ok := db.items[key]
	if !ok {
		return ErrorNotFound
	}

	delete(db.items, key)
	return nil
}
