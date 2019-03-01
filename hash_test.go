package main

import "testing"

func Test_md5Hash(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				in: "yay$&%/12  a",
			},
			want: "1905e8d2f3f43318215f5d15683e8f98",
		},
		{
			args: args{
				in: "drailing.net",
			},
			want: "69683b9e11968b34118cc8665388f0eb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := md5Hash(tt.args.in); got != tt.want {
				t.Errorf("md5Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
