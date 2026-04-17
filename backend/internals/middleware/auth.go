package middleware

// func Authenticate(next http.Handler, cfg *handlers.App) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("session_id")
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
//
// 		claims, err := auth.Verifyer(cookie.Value, cfg.JWT)
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			log.Fatal("Couldn't get claims", err)
// 		}
//
// 		uuid, err := claims.GetSubject()
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			log.Fatal("Bad uuid")
// 		}
//
// 		ctx := context.WithValue(r.Context(), "users_id", uuid)
// 		r.WithContext(ctx)
//
// 		next.ServeHTTP(w, r)
// 	})
// }
