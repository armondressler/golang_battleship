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
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
)

type battleshipContextKey string

type jwtBlacklist map[string]int64

var JWTBlacklist = jwtBlacklist{}

type LoginBody struct {
	Playername string `json:"playername"`
	Password   string `json:"password"`
}

func (j *jwtBlacklist) Blacklist(jwtID string, expiry int64) {
	i := *j //not sure why this can't be done in one line
	i[jwtID] = expiry
}

func (j *jwtBlacklist) isBlacklisted(jwtID string) bool {
	i := *j
	if _, ok := i[jwtID]; !ok {
		return false
	}
	return true
}

func (j *jwtBlacklist) PurgeExpiredTokens() {
	now := time.Now().Unix()
	for id, expiry := range *j {
		if expiry < now {
			delete(*j, id)
		}
	}
}

func hashPassword(password string, rehashCount int) (string, error) {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), rehashCount)
	return string(byteHash), err
}

func createToken(signingKey []byte, user string, expiresInSeconds int) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	jwtID := uuid.New()
	t.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(expiresInSeconds)).Unix(),
		Id:        jwtID.String(),
		Subject:   user,
	}
	s, err := t.SignedString(signingKey)
	return s, err
}

func getPlayerFromContext(r *http.Request) (player.Player, error) {
	var playernameKey battleshipContextKey = "jwtPlayername"
	c := r.Context()
	return player.GetByName(c.Value(playernameKey).(string))
}

func getJwtExpiryFromContext(r *http.Request) (int64, error) {
	var expiryKey battleshipContextKey = "jwtExpiry"
	c := r.Context()
	v, ok := c.Value(expiryKey).(int64)
	if !ok {
		return 0, fmt.Errorf("jwt payload key jwtExpiry missing in context")
	}
	return v, nil
}

func getJwtIDFromContext(r *http.Request) (string, error) {
	var jwtIDKey battleshipContextKey = "jwtID"
	c := r.Context()
	v, ok := c.Value(jwtIDKey).(string)
	if !ok {
		return "", fmt.Errorf("jwt payload key jwtID missing in context")
	}
	return v, nil
}

func Logout(w http.ResponseWriter, r *http.Request, jwtID string, jwtExpiry int64) {
	JWTBlacklist.Blacklist(jwtID, jwtExpiry)
	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
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
}

func (jwtm JWTMiddleware) CheckJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(jwtm.jwtCookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
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
		if JWTBlacklist.isBlacklisted(claims.Id) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), battleshipContextKey("jwtPlayername"), claims.Subject)
		ctx = context.WithValue(ctx, battleshipContextKey("jwtID"), claims.Id)
		ctx = context.WithValue(ctx, battleshipContextKey("jwtExpiry"), claims.ExpiresAt)
		h.ServeHTTP(w, r.Clone(ctx))
	})
}
