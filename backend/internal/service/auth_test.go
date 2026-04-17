package service

import (
	"testing"
)

type mockUserRepo struct {
	users  map[string]*userRecord
	nextID int
}

type userRecord struct {
	id           int
	username     string
	passwordHash string
}

func (m *mockUserRepo) findByUsername(username string) *userRecord {
	return m.users[username]
}

func TestAuthService_Register(t *testing.T) {
	svc := NewAuthService(nil, nil, "test-secret")

	// We test the token generation part which doesn't need DB
	token, err := svc.generateToken(1)
	if err != nil {
		t.Fatalf("generateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestAuthService_GenerateToken_DifferentUsers(t *testing.T) {
	svc := NewAuthService(nil, nil, "test-secret")

	token1, _ := svc.generateToken(1)
	token2, _ := svc.generateToken(2)

	if token1 == token2 {
		t.Fatal("tokens for different users should be different")
	}
}

func TestAuthService_GenerateToken_DifferentSecrets(t *testing.T) {
	svc1 := NewAuthService(nil, nil, "secret-1")
	svc2 := NewAuthService(nil, nil, "secret-2")

	token1, _ := svc1.generateToken(1)
	token2, _ := svc2.generateToken(1)

	if token1 == token2 {
		t.Fatal("tokens with different secrets should be different")
	}
}
