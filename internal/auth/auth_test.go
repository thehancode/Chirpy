package auth

import (
	"net/http"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// TestHashPassword tests the HashPassword function
func TestHashPassword(t *testing.T) {
	password := "mySecretPassword123"
	hashedPassword, err := HashPassword(password)
	// Ensure no error was returned
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Ensure the hashed password is valid and can be compared correctly
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Fatalf("expected the hashed password to be valid, got %v", err)
	}
}

func TestGetBearerToken(t *testing.T) {
	testCases := []struct {
		name          string
		headers       http.Header
		expectedToken string
		expectedError string
	}{
		{
			name: "Valid token",
			headers: http.Header{
				"Authorization": []string{"Bearer abcdef12345"},
			},
			expectedToken: "abcdef12345",
			expectedError: "",
		},
		{
			name:          "Missing Authorization header",
			headers:       http.Header{},
			expectedToken: "",
			expectedError: "authorization header not found",
		},
		{
			name: "Invalid format (not Bearer)",
			headers: http.Header{
				"Authorization": []string{"Basic abcdef12345"},
			},
			expectedToken: "",
			expectedError: "invalid authorization header format",
		},
		{
			name: "Empty token after Bearer",
			headers: http.Header{
				"Authorization": []string{"Bearer   "},
			},
			expectedToken: "",
			expectedError: "token string is empty",
		},
		{
			name: "Bearer without space",
			headers: http.Header{
				"Authorization": []string{"Bearerabcdef12345"},
			},
			expectedToken: "",
			expectedError: "invalid authorization header format",
		},
		{
			name: "Leading whitespace before Bearer",
			headers: http.Header{
				"Authorization": []string{"   Bearer abcdef12345"},
			},
			expectedToken: "",
			expectedError: "invalid authorization header format",
		},
		{
			name: "Trailing whitespace after token",
			headers: http.Header{
				"Authorization": []string{"Bearer abcdef12345   "},
			},
			expectedToken: "abcdef12345",
			expectedError: "",
		},
		{
			name: "Extra spaces between Bearer and token",
			headers: http.Header{
				"Authorization": []string{"Bearer     abcdef12345"},
			},
			expectedToken: "abcdef12345",
			expectedError: "",
		},
		{
			name: "Bearer in lowercase",
			headers: http.Header{
				"Authorization": []string{"bearer abcdef12345"},
			},
			expectedToken: "",
			expectedError: "invalid authorization header format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GetBearerToken(tc.headers)
			if tc.expectedError != "" {
				if err == nil {
					t.Errorf("expected error %q but got nil", tc.expectedError)
				} else if err.Error() != tc.expectedError {
					t.Errorf("expected error %q but got %q", tc.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %q", err.Error())
				}
				if token != tc.expectedToken {
					t.Errorf("expected token %q but got %q", tc.expectedToken, token)
				}
			}
		})
	}
}
