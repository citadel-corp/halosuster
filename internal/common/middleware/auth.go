package middleware

import (
	"context"
	"net/http"

	"github.com/citadel-corp/halosuster/internal/common/jwt"
)

type ContextAuthKey struct{}

func AuthorizeITUser(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(tokenString) <= len("Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):]
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		subject, userType, err := jwt.VerifyAndGetSubject(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if userType != "IT" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextAuthKey{}, subject)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func AuthorizeITAndNurseUser(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(tokenString) <= len("Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):]
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		subject, _, err := jwt.VerifyAndGetSubject(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextAuthKey{}, subject)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
