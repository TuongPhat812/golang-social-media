package user

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

func TestNewTokenID(t *testing.T) {
	token := "test-token-string"
	tokenID := NewTokenID(token)

	if tokenID.value == "" {
		t.Error("TokenID.value should not be empty")
	}

	// Verify it's a SHA256 hash
	expectedHash := sha256.Sum256([]byte(token))
	expectedValue := hex.EncodeToString(expectedHash[:])

	if tokenID.value != expectedValue {
		t.Errorf("TokenID.value = %v, want %v", tokenID.value, expectedValue)
	}
}

func TestTokenID_String(t *testing.T) {
	token := "test-token"
	tokenID := NewTokenID(token)

	if tokenID.String() != tokenID.value {
		t.Errorf("TokenID.String() = %v, want %v", tokenID.String(), tokenID.value)
	}
}

func TestTokenID_Consistency(t *testing.T) {
	token := "same-token"
	tokenID1 := NewTokenID(token)
	tokenID2 := NewTokenID(token)

	if tokenID1.String() != tokenID2.String() {
		t.Error("Same token should produce same TokenID")
	}
}

func TestNewPassword(t *testing.T) {
	plainPassword := "mypassword123"
	password := NewPassword(plainPassword)

	if password.hashed == "" {
		t.Error("Password.hashed should not be empty")
	}

	if password.hashed == plainPassword {
		t.Error("Password should be hashed, not plain text")
	}
}

func TestPassword_String(t *testing.T) {
	plainPassword := "mypassword123"
	password := NewPassword(plainPassword)

	if password.String() != password.hashed {
		t.Errorf("Password.String() = %v, want %v", password.String(), password.hashed)
	}
}

func TestPassword_Verify(t *testing.T) {
	plainPassword := "mypassword123"
	password := NewPassword(plainPassword)

	tests := []struct {
		name          string
		inputPassword string
		wantValid     bool
	}{
		{
			name:          "correct password",
			inputPassword: plainPassword,
			wantValid:     true,
		},
		{
			name:          "incorrect password",
			inputPassword: "wrongpassword",
			wantValid:     false,
		},
		{
			name:          "empty password",
			inputPassword: "",
			wantValid:     false,
		},
		{
			name:          "similar password",
			inputPassword: "mypassword1234",
			wantValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := password.Verify(tt.inputPassword)
			if valid != tt.wantValid {
				t.Errorf("Password.Verify() = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

func TestPassword_Consistency(t *testing.T) {
	plainPassword := "samepassword"
	password1 := NewPassword(plainPassword)
	password2 := NewPassword(plainPassword)

	// Same password should produce same hash
	if password1.String() != password2.String() {
		t.Error("Same password should produce same hash")
	}
}

func TestNewEmail(t *testing.T) {
	emailStr := "test@example.com"
	email := NewEmail(emailStr)

	if email.value != emailStr {
		t.Errorf("Email.value = %v, want %v", email.value, emailStr)
	}
}

func TestEmail_String(t *testing.T) {
	emailStr := "test@example.com"
	email := NewEmail(emailStr)

	if email.String() != emailStr {
		t.Errorf("Email.String() = %v, want %v", email.String(), emailStr)
	}
}

func TestNewName(t *testing.T) {
	nameStr := "John Doe"
	name := NewName(nameStr)

	if name.value != nameStr {
		t.Errorf("Name.value = %v, want %v", name.value, nameStr)
	}
}

func TestName_String(t *testing.T) {
	nameStr := "John Doe"
	name := NewName(nameStr)

	if name.String() != nameStr {
		t.Errorf("Name.String() = %v, want %v", name.String(), nameStr)
	}
}

func TestNewTimestamp(t *testing.T) {
	now := time.Now()
	timestamp := NewTimestamp(now)

	if timestamp.value != now {
		t.Errorf("Timestamp.value = %v, want %v", timestamp.value, now)
	}
}

func TestTimestamp_Value(t *testing.T) {
	now := time.Now()
	timestamp := NewTimestamp(now)

	// Access the value field directly (assuming it's exported or we have a getter)
	// Since we don't have a getter, we'll just verify the struct is created correctly
	if timestamp.value.IsZero() {
		t.Error("Timestamp.value should not be zero")
	}
}

