package api

import (
	"golang_battleship/cmd"
	"golang_battleship/game"
	"golang_battleship/player"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
)

func TestGenerateJwtSigningKey(t *testing.T) {
	type args struct {
		keysize int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				12,
			},
			want:    "abcdef",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.GenerateRandomKey(tt.args.keysize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJwtSigningKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.args.keysize {
				t.Errorf("GenerateJwtSigningKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateToken(t *testing.T) {
	type args struct {
		signingKey       []byte
		user             string
		expiresInSeconds int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				signingKey:       []byte("abcdefg"),
				user:             "testuser",
				expiresInSeconds: 10,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQ0NTMyODIsInN1YiI6InRlc3R1c2VyIn0.",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createToken(tt.args.signingKey, tt.args.user, tt.args.expiresInSeconds)
			if (err != nil) != tt.wantErr {
				t.Errorf("createToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.HasPrefix(got, tt.want) {
				t.Errorf("createToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	jwtSigningKey := []byte("abcdefg")
	finish := make(chan struct{})
	jwtm := JWTMiddleware{jwtSigningKey: jwtSigningKey, loginHandler: Login}
	defaultRouter := mux.NewRouter()
	needsAuthRouter := defaultRouter.Path("/games").Subrouter()
	needsAuthRouter.Use(jwtm.CheckJWT)
	defaultRouter.Path("/login").Methods("POST").Handler(jwtm)
	needsAuthRouter.Methods("GET").HandlerFunc(ListGames)
	passwordHashRudolf, _ := hashPassword("passwordrudolf", PASSWORD_REHASH_COUNT)
	player.NewPlayer("Rudolf", passwordHashRudolf)
	passwordHashDagobert, _ := hashPassword("passworddagobert", PASSWORD_REHASH_COUNT)
	player.NewPlayer("Dagobert", passwordHashDagobert)
	game.NewGame(12, 12, 6, "New Game", 2, "Rudolf", "Dagobert")

	go func() {
		if err := http.ListenAndServe("127.0.0.1:8080", defaultRouter); err != nil {
			panic(err)
		}
	}()

	go func() {
		cli := &http.Client{
			Timeout: time.Second * 1,
		}
		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/login").
			Body("{\"playername\":\"Rudolf\",\"password\":\"BADPASSWORD\"}").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()

		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/login").
			Body("{\"playername\":\"Rudolf\"}").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()

		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/login").
			Body("{\"playername\":\"Rudolf\",\"password\":\"\x00\"}").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()

		loginresponse := apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/login").
			Body("{\"playername\":\"Rudolf\",\"password\":\"passwordrudolf\"}").
			Expect(t).
			Status(http.StatusOK).
			CookiePresent(JWT_COOKIE_NAME).
			End().Response

		var jwtCookie string
		for _, c := range loginresponse.Cookies() {
			if c.Name == JWT_COOKIE_NAME {
				jwtCookie = c.Value
			}
		}
		if jwtCookie == "" {
			t.Errorf("no cookie was set after logging in, expected cookie with name %s", JWT_COOKIE_NAME)
			t.Fail()
		}

		apitest.New().
			EnableNetworking(cli).
			Get("http://127.0.0.1:8080/games").
			Cookie(JWT_COOKIE_NAME, jwtCookie).
			Expect(t).
			Status(http.StatusOK).
			End()

		finish <- struct{}{}
	}()
	<-finish
}
