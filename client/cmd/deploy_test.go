package cmd

import (
	"testing"
)

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

func Test_logger(t *testing.T) {
	type args struct {
		message string
		clr     string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger(tt.args.message, tt.args.clr)
		})
	}
}
