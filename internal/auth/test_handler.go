package auth

import (
	"fmt"
	"net/http"
)

func TestAuthHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey)

	if userID == nil {
		http.Error(w, "user id not found", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "authenticated user:", userID)
}
