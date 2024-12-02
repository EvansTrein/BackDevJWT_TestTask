package utils

import "testing"

func TestIsGUID(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test valid GUID",
			args: args{s: "f3a4d346-2596-46bb-a9aa-d33e7f67b127"},
			want: true,
		},
		{
			name: "Test valid GUID",
			args: args{s: "16e0b0ec-1d98-43ec-84a9-66433e291ef7"},
			want: true,
		},
		{
			name: "Test valid GUID",
			args: args{s: "4aa7ca1b-ae64-45d0-89bb-a1dfff4708d9"},
			want: true,
		},
		{
			name: "Test invalid GUID",
			args: args{s: "dfe-a7e1-4593a-a3e3-30e8321236c275ea"},
			want: false,
		},
		{
			name: "Test invalid GUID",
			args: args{s: "32-a23-4d63-8a9f-t5"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGUID(tt.args.s); got != tt.want {
				t.Errorf("IsGUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
