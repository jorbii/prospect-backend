package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (s *Service) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = ValidateRegister(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.RegisterUser(req)
	if err != nil {
		if errors.Is(err, ErrUsernameExists) {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}

		if errors.Is(err, ErrEmailExists) {
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := RegisterResponse{
		Message:  "registration successful",
		Username: req.Username,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := s.LoginUser(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	token, err := s.tokenService.GenerateToken(user.ID.String())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := LoginResponse{
		AccessToken: token,
	}

	json.NewEncoder(w).Encode(response)
}
