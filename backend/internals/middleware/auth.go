package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/thegreatestgiant/Charity-Tracker/internals/auth"
	"github.com/thegreatestgiant/Charity-Tracker/internals/handlers"
)

func Authenticate(w http.ResponseWriter, r *http.Request, cfg *handlers.App) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	claims, err := auth.Verifyer(cookie.Value, cfg.JWT)
	if err != nil {
		log.Fatal("Couldn't get claims", err)
	}

	ctx := context.WithValue(r.Context(), "users_id", claims.Subject)
	r.WithContext(ctx)
}
