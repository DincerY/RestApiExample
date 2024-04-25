package app

import (
	"RestApiExample/models"
	u "RestApiExample/utils"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login"}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		var tokenPart string

		if tokenHeader == "" {
			cookies := r.Cookies()
			var cookie http.Cookie
			if len(cookies) != 0 {
				cookie = *cookies[0]
				tokenPart = cookie.Value
			} else {
				response = u.Message(false, "Token gönderilmelidir")
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				u.Respond(w, response)
				return
			}
		} else {
			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				response := u.Message(false, "Hatalı yada geçersiz token")
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				u.Respond(w, response)
				return
			}
			tokenPart = splitted[1]
		}

		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			response = u.Message(false, "Token hatalı! Cookie mevcutsa heeader tarafında auhtorization kısımını silmeniz yeterli")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid {
			response = u.Message(false, "Token geçersiz!")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		fmt.Sprintf("Kullanıcı %v", tk.Username)
		ctx := context.WithValue(r.Context(), "token", tk)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}
