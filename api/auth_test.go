package api

import (
	"golang_battleship/cmd"
	"strings"
	"testing"
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
			got, err := cmd.GenerateJwtSigningKey(tt.args.keysize)
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

func Test_createToken(t *testing.T) {
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
