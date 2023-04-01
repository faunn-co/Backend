package encrypt

import (
	"context"
	"fmt"
	"testing"
)

func TestHashAndSalt(t *testing.T) {
	type args struct {
		ctx context.Context
		pwd string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "mock",
			args: args{
				ctx: nil,
				pwd: "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashAndSalt(tt.args.ctx, tt.args.pwd); got != tt.want {
				fmt.Println(got)
				t.Errorf("HashAndSalt() = %v, want %v", got, tt.want)
			}
		})
	}
}
