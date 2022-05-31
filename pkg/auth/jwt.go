package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"NKSS-backend/pkg/handlers"

	"github.com/golang-jwt/jwt"
)

func Authenticator(next http.Handler, secret []byte) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization via JWT
		value := r.Header.Get("Authorization")
		token := strings.SplitN(value, "Bearer ", 2)
		if len(token) == 1 {
			handlers.RespondError(w, 401, "Token is absent")
			return
		}
		rollno, ok := validateJWT(token[1], secret)
		if !ok {
			handlers.RespondError(w, 400, "Token is invalid")
			return
		}

		ctx := context.WithValue(r.Context(), "rollno", rollno)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateJWT(role string, rollno string, hmacSecret []byte) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role":   role,
		"rollno": rollno,
	})

	tokenString, err := token.SignedString(hmacSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An unexpected error occurred: %s", err)
		return ""
	} else {
		return tokenString
	}

}

func validateJWT(tokenString string, hmacSecret []byte) (rollno string, ok bool) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			_, err := fmt.Fprintf(os.Stderr, "Unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}

		return hmacSecret, nil
	})

	if err != nil {
		fmt.Println(err)
		return "", false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["rollno"].(string), true
	} else {
		return "", false
	}

}
