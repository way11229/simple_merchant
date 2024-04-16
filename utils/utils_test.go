package utils

import "testing"

func TestConvertStringToUint(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				str: "100",
			},
			want:    100,
			wantErr: false,
		},
		{
			name: "error case",
			args: args{
				"abc",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertStringToUint(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertStringToUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertStringToUint() = %v, want %v", got, tt.want)
			}
		})
	}
}
