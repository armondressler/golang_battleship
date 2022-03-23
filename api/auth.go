package api

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_battleship/player"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
)

type LoginBody struct {
	Playername string `json:"playername"`
	Password   string `json:"password"`
}

func hashPassword(password string, rehashCount int) (string, error) {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), rehashCount)
	return string(byteHash), err
}

func createToken(signingKey []byte, user string, expiresInSeconds int) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(expiresInSeconds)).Unix(),
		Subject:   user,
	}
	s, err := t.SignedString(signingKey)
	return s, err
}

func getPlayerFromContext(r *http.Request) (player.Player, error) {
	var playernameKey Playername = "playername"
	c := r.Context()
	return player.GetByName(c.Value(playernameKey).(string))
}

func Login(w http.ResponseWriter, r *http.Request, jwtsigningkey []byte) {
	var b LoginBody
	b.Playername = r.PostFormValue("playername")
	b.Password = r.PostFormValue("password")
	if len(b.Playername) == 0 && len(b.Password) == 0 {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&b)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
	}
	p, err := player.GetByName(b.Playername)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(b.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	t, err := createToken(jwtsigningkey, b.Playername, 60*60)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c := http.Cookie{
		Name:     JWT_COOKIE_NAME,
		Value:    t,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Expires:  time.Now().AddDate(1, 0, 0),
	}
	http.SetCookie(w, &c)
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	http.Redirect(w, r, "/dashboard.html", http.StatusSeeOther)
	//w.WriteHeader(http.StatusOK)
}

func (jwtm JWTMiddleware) CheckJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(JWT_COOKIE_NAME)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Debug(fmt.Sprintf("Checking token %s from cookie %s", c.Value, c.Name))
		t, err := jwt.ParseWithClaims(c.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtm.jwtSigningKey, nil
		})
		if err != nil {
			log.Warn("malformed token, ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !t.Valid {
			log.Warn("invalid token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims, ok := t.Claims.(*jwt.StandardClaims)
		if !ok {
			log.Warn("malformed claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var playernameKey Playername = "playername" //because context.WithValue told me so
		ctx := context.WithValue(r.Context(), playernameKey, claims.Subject)
		h.ServeHTTP(w, r.Clone(ctx))
	})
}
