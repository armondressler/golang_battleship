package api

import (
	"encoding/json"
	"golang_battleship/player"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type LoginBody struct {
	Playername     string `json:"Playername"`
	PasswordBCrypt string `json:"PasswordBCrypt"`
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

func Login(w http.ResponseWriter, r *http.Request, jwtsigningkey []byte) {
	var b LoginBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	p, err := player.GetByName(b.Playername)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if p.PasswordBCrypt != b.PasswordBCrypt {
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
	log.Info(c)
	http.SetCookie(w, &c)
	w.WriteHeader(http.StatusAccepted)

}
