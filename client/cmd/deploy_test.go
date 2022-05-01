package cmd

import (
	"testing"
)

func Test_sendtocrypt(t *testing.T) {
	type args struct {
		pass         string
		username     string
		hostname     string
		serialnumber string
		environment  string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendtocrypt(tt.args.pass, tt.args.username, tt.args.hostname, tt.args.serialnumber, tt.args.environment)
		})
	}
}

func Test_getConfigInt(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfigInt(tt.args.key); got != tt.want {
				t.Errorf("getConfigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConfigStr(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfigStr(tt.args.key); got != tt.want {
				t.Errorf("getConfigStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generatePassword(t *testing.T) {
	type args struct {
		passwordLength int
		minSpecialChar int
		minNum         int
		minUpperCase   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generatePassword(tt.args.passwordLength, tt.args.minSpecialChar, tt.args.minNum, tt.args.minUpperCase); got != tt.want {
				t.Errorf("generatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
