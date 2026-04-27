// Package executor provides secure command execution helpers.
package executor

// SecurePassword wraps a password behind a pointer so that Bubble Tea's
// by-value Model copies share one backing array instead of each holding
// an independent, unzeroed copy. Call Close() to zero and release memory.
type SecurePassword struct {
	data []byte
}

// NewSecurePassword copies input into a new SecurePassword and zeroes
// the source slice. The caller must not use input after this call.
func NewSecurePassword(input []byte) *SecurePassword {
	cp := make([]byte, len(input))
	copy(cp, input)
	ZeroBytes(input)
	return &SecurePassword{data: cp}
}

// Bytes returns the raw password bytes. Returns nil if closed.
func (s *SecurePassword) Bytes() []byte {
	if s == nil || s.data == nil {
		return nil
	}
	return s.data
}

// Copy returns a fresh copy of the password for one-time use.
// The caller is responsible for zeroing the returned slice.
func (s *SecurePassword) Copy() []byte {
	if s == nil || s.data == nil {
		return nil
	}
	cp := make([]byte, len(s.data))
	copy(cp, s.data)
	return cp
}

// Close zeroes and releases the underlying password memory.
// Safe to call multiple times or on a nil receiver.
func (s *SecurePassword) Close() {
	if s == nil || s.data == nil {
		return
	}
	ZeroBytes(s.data)
	s.data = nil
}
