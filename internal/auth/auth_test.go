package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGetAuthIdentifier(t *testing.T) {
	userID := uuid.New().String()
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{
				ctx: context.WithValue(context.TODO(), JwtUserIDContextKey, userID),
			},
			want:    userID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAuthIdentifier(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAuthIdentifier() got = %v, want %v", got, tt.want)
			}
		})
	}
}
