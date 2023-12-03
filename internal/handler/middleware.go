package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("You hit the: " + r.RequestURI)
		// Do stuff here
		tokenString, err := r.Cookie("Authorization")
		if err != nil {
			log.Println("Can't get the cookie")
			w.Header().Set("userID", "guest")
		} else if tokenString.Value == "" {
			log.Println("User logged out")
			w.Header().Set("userID", "guest")
		} else {

			token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				log.Println(err)
				w.Header().Set("userID", "guest")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if float64(time.Now().Unix()) > claims["exp"].(float64) {
					log.Println("Token expired")
					w.Header().Set("userID", "guest")
				} else {
					w.Header().Set("userID", claims["sub"].(string))
				}

			} else {
				log.Println(err)
				w.Header().Set("userID", "guest")
			}
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
