package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
)

func TestRegister(t *testing.T) {
	srv := &http.Server{Addr: "127.0.0.1:8080"}
	finish := make(chan struct{})

	http.HandleFunc("/home", home)
	http.HandleFunc("/register", register)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	go func() {

		cli := &http.Client{
			Timeout: time.Second * 1,
		}

		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/register").
			JSON(`{"name": "Rudolf"}`).
			Expect(t).
			Status(http.StatusOK).
			End()

		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/register").
			JSON(`{"name": "Dagobert"}`).
			Expect(t).
			Status(http.StatusOK).
			End()

		apitest.New().
			EnableNetworking(cli).
			Get("http://127.0.0.1:8080/home").
			Expect(t).
			Body(`["Rudolf","Dagobert"]`).
			Status(http.StatusOK).
			End()
		finish <- struct{}{}
	}()
	<-finish
}
