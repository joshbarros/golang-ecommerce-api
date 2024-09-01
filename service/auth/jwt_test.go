package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/joshbarros/golang-ecommerce-api/config"
	"github.com/joshbarros/golang-ecommerce-api/types"
)

// Mock implementation of the UserStore interface
type mockUserStore struct {
	users map[int]*types.User
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserStore) CreateUser(user types.User) error {
	if _, exists := m.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	m.users[user.ID] = &user
	return nil
}

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")
	config.Envs.JWTSecret = "secret"          // Set the secret in the config
	config.Envs.JWTExpirationInSeconds = 3600 // 1 hour for testing

	token, err := CreateJWT(secret, 1)

	if err != nil {
		t.Errorf("Error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be not empty")
	}

	parsedToken, err := validateToken(token)
	if err != nil {
		t.Errorf("Error validating JWT: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Expected token to be valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		t.Error("Expected valid claims")
	}

	userID, err := strconv.Atoi(claims["userID"].(string))
	if err != nil || userID != 1 {
		t.Errorf("Expected userID to be 1, got %d", userID)
	}
}

func TestWithJWTAuth(t *testing.T) {
	secret := []byte("secret")
	config.Envs.JWTSecret = string(secret)

	mockStore := &mockUserStore{
		users: map[int]*types.User{
			1: {ID: 1, Email: "user@example.com", FirstName: "John", LastName: "Doe"},
		},
	}

	// Create a valid JWT for testing
	token, _ := CreateJWT(secret, 1)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"Valid Token", token, http.StatusOK},
		{"Invalid Token", "invalidtoken", http.StatusForbidden},
		{"Missing Token", "", http.StatusForbidden},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/test", nil)
		if tt.token != "" {
			req.Header.Set("Authorization", tt.token)
		}

		rr := httptest.NewRecorder()
		handler := WithJWTAuth(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}, mockStore)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
		}
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserKey, 1)
	userID := GetUserIDFromContext(ctx)

	if userID != 1 {
		t.Errorf("Expected userID to be 1, got %d", userID)
	}

	// Test with no user ID in context
	emptyCtx := context.Background()
	userID = GetUserIDFromContext(emptyCtx)

	if userID != -1 {
		t.Errorf("Expected userID to be -1, got %d", userID)
	}
}

func TestGetTokenFromRequest(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
		queryParam    string
		expectedToken string
	}{
		{"Token in Authorization Header", "Bearer testtoken", "", "Bearer testtoken"},
		{"Token in Query Parameter", "", "testtoken", "testtoken"},
		{"No Token Provided", "", "", ""},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/test", nil)
		if tt.authorization != "" {
			req.Header.Set("Authorization", tt.authorization)
		}
		if tt.queryParam != "" {
			q := req.URL.Query()
			q.Add("token", tt.queryParam)
			req.URL.RawQuery = q.Encode()
		}

		token := getTokenFromRequest(req)

		if token != tt.expectedToken {
			t.Errorf("Expected token to be %s, got %s", tt.expectedToken, token)
		}
	}
}
