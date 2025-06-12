package utils

import (
	"context"
	"testing"
)

func TestSetAndGetContextValue(t *testing.T) {
	t.Run("should set and get string value successfully", func(t *testing.T) {
		ctx := context.Background()
		key := ctxKey("test")
		value := "test value"

		ctx = SetContextValue(ctx, key, value)
		result, ok := GetContextValue[string](ctx, key)

		if !ok {
			t.Error("expected to get value successfully")
		}
		if result != value {
			t.Errorf("expected value %v, got %v", value, result)
		}
	})

	t.Run("should set and get struct value successfully", func(t *testing.T) {
		type TestStruct struct {
			Name string
			Age  int
		}

		ctx := context.Background()
		key := ctxKey("test")
		value := TestStruct{Name: "John", Age: 30}

		ctx = SetContextValue(ctx, key, value)
		result, ok := GetContextValue[TestStruct](ctx, key)

		if !ok {
			t.Error("expected to get value successfully")
		}
		if result != value {
			t.Errorf("expected value %v, got %v", value, result)
		}
	})

	t.Run("should handle nil context value", func(t *testing.T) {
		ctx := context.Background()
		key := ctxKey("nonexistent")

		result, ok := GetContextValue[string](ctx, key)

		if ok {
			t.Error("expected get value to fail")
		}
		if result != "" {
			t.Errorf("expected empty string, got %v", result)
		}
	})

	t.Run("should handle type mismatch", func(t *testing.T) {
		ctx := context.Background()
		key := ctxKey("test")
		value := 42 // storing int

		ctx = SetContextValue(ctx, key, value)
		result, ok := GetContextValue[string](ctx, key) // trying to get as string

		if ok {
			t.Error("expected type assertion to fail")
		}
		if result != "" {
			t.Errorf("expected empty string, got %v", result)
		}
	})
}
