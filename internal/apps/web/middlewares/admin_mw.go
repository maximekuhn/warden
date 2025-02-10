package middlewares

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AdminMiddleware is a middleware that protects admin routes
type AdminMiddleware struct {
	username             string
	bcryptHashedPassword string
}

func NewAdminMiddleware(
	username string,
	bcryptHashedPassword string,
) *AdminMiddleware {
	return &AdminMiddleware{
		username:             username,
		bcryptHashedPassword: bcryptHashedPassword,
	}
}

func (mw *AdminMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, pass, ok := r.BasicAuth()
		passwordMatch := bcrypt.CompareHashAndPassword([]byte(mw.bcryptHashedPassword), []byte(pass)) == nil
		if !ok || username != mw.username || !passwordMatch {
			// highly advanced method to prevent brute force attacks
			time.Sleep(3 * time.Second)

			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
