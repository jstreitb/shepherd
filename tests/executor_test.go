package tests

import (
	"testing"

	"github.com/jstreitb/baa/internal/executor"
	"github.com/jstreitb/baa/internal/pkgmanager/providers"
)

// ─── ZeroBytes ──────────────────────────────────────────────────────────────

func TestZeroBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"non-empty", []byte("secretpassword")},
		{"single byte", []byte{0x42}},
		{"empty", []byte{}},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor.ZeroBytes(tt.input)
			for i, b := range tt.input {
				if b != 0 {
					t.Errorf("byte at index %d was %d, want 0", i, b)
				}
			}
		})
	}
}

// ─── SecurePassword ─────────────────────────────────────────────────────────

func TestSecurePassword_Basic(t *testing.T) {
	input := []byte("mypassword")
	sp := executor.NewSecurePassword(input)

	// Source should be zeroed after NewSecurePassword.
	for i, b := range input {
		if b != 0 {
			t.Errorf("source byte at index %d was %d, want 0", i, b)
		}
	}

	// Bytes should return the password.
	pw := sp.Bytes()
	if string(pw) != "mypassword" {
		t.Errorf("Bytes() = %q, want %q", pw, "mypassword")
	}

	// Copy should return an independent copy.
	cp := sp.Copy()
	if string(cp) != "mypassword" {
		t.Errorf("Copy() = %q, want %q", cp, "mypassword")
	}
	executor.ZeroBytes(cp) // zero the copy
	// Original should be unaffected.
	if string(sp.Bytes()) != "mypassword" {
		t.Errorf("Bytes() after zeroing copy = %q, want %q", sp.Bytes(), "mypassword")
	}

	// Close should zero and nil the data.
	sp.Close()
	if sp.Bytes() != nil {
		t.Errorf("Bytes() after Close() = %v, want nil", sp.Bytes())
	}
	if sp.Copy() != nil {
		t.Errorf("Copy() after Close() = %v, want nil", sp.Copy())
	}

	// Double close should not panic.
	sp.Close()
}

func TestSecurePassword_NilReceiver(t *testing.T) {
	var sp *executor.SecurePassword
	if sp.Bytes() != nil {
		t.Errorf("nil receiver Bytes() = %v, want nil", sp.Bytes())
	}
	if sp.Copy() != nil {
		t.Errorf("nil receiver Copy() = %v, want nil", sp.Copy())
	}
	// Should not panic.
	sp.Close()
}

// ─── BuildInteractiveCmd ────────────────────────────────────────────────────

func TestBuildInteractiveCmd_Sudo(t *testing.T) {
	apt := &providers.Apt{}
	cmd := executor.BuildInteractiveCmd(apt)

	if cmd.Path == "" {
		t.Fatal("expected non-empty path")
	}
	// Should start with sudo.
	args := cmd.Args
	if len(args) < 1 || args[0] != "sudo" {
		t.Errorf("expected args[0] = sudo, got %v", args)
	}
}

func TestBuildInteractiveCmd_NoSudo(t *testing.T) {
	brew := &providers.Brew{}
	cmd := executor.BuildInteractiveCmd(brew)

	args := cmd.Args
	if len(args) < 1 || args[0] != "bash" {
		t.Errorf("expected args[0] = bash, got %v", args)
	}
	// Should NOT contain sudo.
	for _, a := range args {
		if a == "sudo" {
			t.Error("non-sudo manager should not use sudo")
		}
	}
}
