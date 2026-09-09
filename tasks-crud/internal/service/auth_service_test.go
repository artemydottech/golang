package service

import (
	"testing"
	"time"

	"tasks-crud/internal/config"
	"tasks-crud/internal/domain"
	"tasks-crud/internal/repository"
)

func newTestAuthService(t *testing.T) *AuthService {
	t.Helper()

	cfg := &config.Config{
		JWTSecret:  "test-secret",
		JWTExpiry:  time.Hour,
		BcryptCost: 4,
	}

	return NewAuthService(repository.NewInMemoryUserRepository(), cfg)
}

func validRegistration() domain.CreateUserRequest {
	return domain.CreateUserRequest{
		Username: "tester",
		Email:    "tester@example.com",
		Password: "password123",
	}
}

func TestSeededAdminCanLogIn(t *testing.T) {
	svc := newTestAuthService(t)

	user, err := svc.Login(domain.LoginRequest{Email: "admin@example.com", Password: "admin123"})
	if err != nil {
		t.Fatalf("seeded admin could not log in: %v", err)
	}
	if user.Username != "admin" {
		t.Fatalf("got user %q, want admin", user.Username)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	svc := newTestAuthService(t)

	if _, err := svc.Login(domain.LoginRequest{Email: "admin@example.com", Password: "nope"}); err == nil {
		t.Fatal("login succeeded with a wrong password")
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	cases := map[string]domain.CreateUserRequest{
		"empty username": {Username: "", Email: "a@b.co", Password: "password123"},
		"short username": {Username: "ab", Email: "a@b.co", Password: "password123"},
		"empty email":    {Username: "tester", Email: "", Password: "password123"},
		"broken email":   {Username: "tester", Email: "tester.example.com", Password: "password123"},
		"short password": {Username: "tester", Email: "a@b.co", Password: "short"},
	}

	for name, req := range cases {
		t.Run(name, func(t *testing.T) {
			svc := newTestAuthService(t)
			if _, err := svc.Register(req); err == nil {
				t.Fatalf("Register accepted %s", name)
			}
		})
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	svc := newTestAuthService(t)

	if _, err := svc.Register(validRegistration()); err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	if _, err := svc.Register(validRegistration()); err == nil {
		t.Fatal("Register accepted a duplicate email")
	}
}

func TestRegisteredUserCanLogInAndValidateToken(t *testing.T) {
	svc := newTestAuthService(t)

	user, err := svc.Register(validRegistration())
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	token, err := svc.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.UserID != user.ID || claims.Email != user.Email {
		t.Fatalf("claims %+v do not describe user %+v", claims, user)
	}
}

func TestValidateTokenRejectsForeignSecret(t *testing.T) {
	issuer := newTestAuthService(t)
	user, err := issuer.Register(validRegistration())
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	token, err := issuer.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	verifier := NewAuthService(repository.NewInMemoryUserRepository(), &config.Config{
		JWTSecret:  "another-secret",
		JWTExpiry:  time.Hour,
		BcryptCost: 4,
	})

	if _, err := verifier.ValidateToken(token); err == nil {
		t.Fatal("token signed with a different secret was accepted")
	}
}

func TestValidateTokenRejectsExpired(t *testing.T) {
	svc := NewAuthService(repository.NewInMemoryUserRepository(), &config.Config{
		JWTSecret:  "test-secret",
		JWTExpiry:  -time.Minute,
		BcryptCost: 4,
	})

	token, err := svc.GenerateToken(&domain.User{ID: 1, Username: "a", Email: "a@b.co"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if _, err := svc.ValidateToken(token); err == nil {
		t.Fatal("expired token was accepted")
	}
}
