package encrypt

import (
	"context"
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
			name: "1",
			args: args{nil, "123456"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashAndSalt(tt.args.ctx, tt.args.pwd); got != tt.want {
				t.Errorf("HashAndSalt() = %v, want %v", got, tt.want)
			}
		})
	}
}
