package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "mypassword123"

		hash, err := HashPassword(password)

		if err != nil {
			t.Errorf("HashPassword failed: %v", err)
		}
		if hash == "" {
			t.Error("hash should not be empty")
		}
		if hash == password {
			t.Error("hash should not be equal to original password")
		}
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("should validate correct password", func(t *testing.T) {
		password := "mypassword123"
		hash, _ := HashPassword(password)

		isValid := CheckPasswordHash(password, hash)

		if !isValid {
			t.Error("password should be valid")
		}
	})

	t.Run("should invalidate wrong password", func(t *testing.T) {
		password := "mypassword123"
		wrongPassword := "wrongpassword"
		hash, _ := HashPassword(password)

		isValid := CheckPasswordHash(wrongPassword, hash)

		if isValid {
			t.Error("password should be invalid")
		}
	})
}

func TestGenerateJWT(t *testing.T) {
	originalSecretKey := os.Getenv("SECRET_KEY")
	os.Setenv("SECRET_KEY", "test-secret-key")
	defer os.Setenv("SECRET_KEY", originalSecretKey)

	t.Run("should generate valid JWT", func(t *testing.T) {
		userId := "123"
		email := "test@example.com"

		token, err := GenerateJWT(userId, email)

		if err != nil {
			t.Errorf("GenerateJWT failed: %v", err)
		}
		if token == "" {
			t.Error("token should not be empty")
		}

		// Validate the generated token
		claims, err := ValidateToken(token)
		if err != nil {
			t.Errorf("Generated token is invalid: %v", err)
		}
		if claims.UserId != userId {
			t.Errorf("Expected userId %s, got %s", userId, claims.UserId)
		}
		if claims.Email != email {
			t.Errorf("Expected email %s, got %s", email, claims.Email)
		}
	})
}

func TestValidateToken(t *testing.T) {
	originalSecretKey := os.Getenv("SECRET_KEY")
	os.Setenv("SECRET_KEY", "test-secret-key")
	defer os.Setenv("SECRET_KEY", originalSecretKey)

	t.Run("should validate correct token", func(t *testing.T) {
		token, _ := GenerateJWT("123", "test@example.com")

		claims, err := ValidateToken(token)

		if err != nil {
			t.Errorf("ValidateToken failed: %v", err)
		}
		if claims == nil {
			t.Error("claims should not be nil")
		}
	})

	t.Run("should reject invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		claims, err := ValidateToken(invalidToken)

		if err == nil {
			t.Error("expected error for invalid token")
		}
		if claims != nil {
			t.Error("claims should be nil for invalid token")
		}
	})
}

func TestExtractToken(t *testing.T) {
	originalSecretKey := os.Getenv("SECRET_KEY")
	os.Setenv("SECRET_KEY", "test-secret-key")
	defer os.Setenv("SECRET_KEY", originalSecretKey)

	t.Run("should extract valid token from header", func(t *testing.T) {
		token, _ := GenerateJWT("123", "test@example.com")
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		claims := ExtractToken(w, req)

		if claims == nil {
			t.Error("claims should not be nil")
		} else {
			if claims.UserId != "123" {
				t.Errorf("Expected userId 123, got %s", claims.UserId)
			}
			if claims.Email != "test@example.com" {
				t.Errorf("Expected email test@example.com, got %s", claims.Email)
			}
		}
	})

	t.Run("should handle missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		claims := ExtractToken(w, req)

		if claims != nil {
			t.Error("claims should be nil")
		}
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("should handle invalid token format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		claims := ExtractToken(w, req)

		if claims != nil {
			t.Error("claims should be nil")
		}
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}
