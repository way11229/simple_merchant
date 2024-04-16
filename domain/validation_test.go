package domain

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "happy case",
			args: args{
				email: "test@gmail.com",
			},
			want: true,
		},
		{
			name: "fail case",
			args: args{
				email: "aaabbsssccc",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmail(tt.args.email); got != tt.want {
				t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateUserPassword(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "happy case",
			args: args{
				pwd: "aaaAAA()",
			},
			want: true,
		},
		{
			name: "too short",
			args: args{
				pwd: "a",
			},
			want: false,
		},
		{
			name: "too long",
			args: args{
				pwd: "aaaaaAAAA!##AAaaa",
			},
			want: false,
		},
		{
			name: "without lower char",
			args: args{
				pwd: "AAAAAA@@",
			},
			want: false,
		},
		{
			name: "without upper char",
			args: args{
				pwd: "aaaaaa!!",
			},
			want: false,
		},
		{
			name: "without special char",
			args: args{
				pwd: "aaaaaaAAA",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateUserPassword(tt.args.pwd); got != tt.want {
				t.Errorf("ValidateUserPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
