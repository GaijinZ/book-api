package models

import (
	"testing"
)

func TestUser_ValidateUser(t *testing.T) {
	type fields struct {
		ID        int
		Firstname string
		Lastname  string
		Email     string
		Password  string
		Role      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid user",
			fields: fields{
				Firstname: "John",
				Lastname:  "Doe",
				Password:  "securepassword",
			},
			wantErr: false,
		},

		{
			name: "Invalid firstname",
			fields: fields{
				Firstname: "John123",
				Lastname:  "Doe",
				Password:  "securepassword",
			},
			wantErr: true,
		},

		{
			fields: fields{
				Firstname: "John",
				Lastname:  "Doe123",
				Password:  "securepassword",
			},
			wantErr: true,
		},

		{
			name: "Blank password",
			fields: fields{
				Firstname: "John",
				Lastname:  "Doe",
				Password:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Firstname: tt.fields.Firstname,
				Lastname:  tt.fields.Lastname,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				Role:      tt.fields.Role,
			}
			if err := u.ValidateUser(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
