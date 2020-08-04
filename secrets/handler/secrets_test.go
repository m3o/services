package handler

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	t.Run("MissingSecret", func(t *testing.T) {
		_, err := new(Secrets).encrypt("foo")
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("InvalidSecret", func(t *testing.T) {
		_, err := (&Secrets{secret: "bar"}).encrypt("foo")
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})

	t.Run("ValidSecret", func(t *testing.T) {
		h := &Secrets{secret: "6368616e676520746869732070617373"}
		bytes, err := h.encrypt("foo")
		if err != nil {
			t.Errorf("Expected nil error but got %v", err)
			return
		}
		if bytes == nil {
			t.Errorf("Expected byted but got nil")
			return
		}
		if string(bytes) == "foo" {
			t.Errorf("Result was the same as the input")
			return
		}

		result, err := h.decrypt(bytes)
		if err != nil {
			t.Errorf("Expected nil error but got %v", err)
			return
		}
		if len(result) == 0 {
			t.Errorf("Expected a result but got a blank string")
			return
		}
		if result != "foo" {
			t.Errorf("Expected foo but got %v", result)
		}
	})
}
