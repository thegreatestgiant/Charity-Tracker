package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/Charity-Tracker/internal/auth"
)

func AuthGuard(next http.Handler, jwt []byte) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := auth.Verifyer(cookie.Value, jwt)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Fatal("Couldn't get claims", err)
		}

		uuidStr, err := claims.GetSubject()
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Bad uuid")
			return
		}

		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			log.Printf("Couldn't get uuid: %v ", err)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user_id", uuid))

		next.ServeHTTP(w, r)
	})
}
