package handlers

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestCheckToken(t *testing.T) {
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckToken(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CheckToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
