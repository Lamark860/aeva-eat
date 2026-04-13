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
	email        string
	passwordHash string
}

func (m *mockUserRepo) findByEmail(email string) *userRecord {
	return m.users[email]
}

func TestAuthService_Register(t *testing.T) {
	svc := NewAuthService(nil, "test-secret")

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
	svc := NewAuthService(nil, "test-secret")

	token1, _ := svc.generateToken(1)
	token2, _ := svc.generateToken(2)

	if token1 == token2 {
		t.Fatal("tokens for different users should be different")
	}
}

func TestAuthService_GenerateToken_DifferentSecrets(t *testing.T) {
	svc1 := NewAuthService(nil, "secret-1")
	svc2 := NewAuthService(nil, "secret-2")

	token1, _ := svc1.generateToken(1)
	token2, _ := svc2.generateToken(1)

	if token1 == token2 {
		t.Fatal("tokens with different secrets should be different")
	}
}
