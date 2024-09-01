package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshbarros/golang-ecommerce-api/service/auth"
	"github.com/joshbarros/golang-ecommerce-api/types"
)

// Mock implementation of UserStore with in-memory data
type mockUserStore struct {
	users map[string]*types.User
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserStore) CreateUser(user types.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	if m.users == nil {
		m.users = make(map[string]*types.User)
	}
	m.users[user.Email] = &user
	return nil
}

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("Should fail if the payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "josue",
			LastName:  "barros",
			Email:     "invalidemail",
			Password:  "87654321",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "josue",
			LastName:  "barros",
			Email:     "validemail@gmail.com",
			Password:  "87654321",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}

		if _, exists := userStore.users[payload.Email]; !exists {
			t.Errorf("Expected user to be created in the store")
		}
	})

	t.Run("Should fail if user with email already exists", func(t *testing.T) {
		userStore.users = map[string]*types.User{
			"validemail@gmail.com": {
				FirstName: "josue",
				LastName:  "barros",
				Email:     "validemail@gmail.com",
				Password:  "hashedpassword",
			},
		}

		payload := types.RegisterUserPayload{
			FirstName: "josue",
			LastName:  "barros",
			Email:     "validemail@gmail.com",
			Password:  "87654321",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Should fail login with incorrect email", func(t *testing.T) {
		payload := types.LoginUserPayload{
			Email:    "nonexistent@gmail.com",
			Password: "wrongpassword",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Should fail login with incorrect password", func(t *testing.T) {
		userStore.users = map[string]*types.User{
			"validemail@gmail.com": {
				FirstName: "josue",
				LastName:  "barros",
				Email:     "validemail@gmail.com",
				Password:  "hashedpassword", // Normally, this would be a hashed password
			},
		}

		payload := types.LoginUserPayload{
			Email:    "validemail@gmail.com",
			Password: "wrongpassword",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Should login successfully with correct credentials", func(t *testing.T) {
		hashedPassword, _ := auth.HashPassword("87654321")
		userStore.users = map[string]*types.User{
			"validemail@gmail.com": {
				FirstName: "josue",
				LastName:  "barros",
				Email:     "validemail@gmail.com",
				Password:  hashedPassword,
			},
		}

		payload := types.LoginUserPayload{
			Email:    "validemail@gmail.com",
			Password: "87654321",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal("Failed to decode JSON response")
		}

		if _, ok := response["token"]; !ok {
			t.Error("Expected a JWT token in the response")
		}
	})
}
