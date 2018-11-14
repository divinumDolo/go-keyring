package keyring

import "fmt"

// provider set in the init function by the relevant os file e.g.:
// keyring_linux.go
var provider Keyring

var (
	// ErrNotFound is the expected error if the secret isn't found in the
	// keyring.
	ErrNotFound = fmt.Errorf("secret not found in keyring")
)

// Keyring provides a simple set/get interface for a keyring service.
type Keyring interface {
	// Set password in keyring for user.
	Set(label string, args []string, password string) error
	// Get password from keyring given label and args list
	Get(args []string) (string, error)
	// Delete secret from keyring.
	Delete(args []string) error
}

// Set password in keyring for args.
func Set(label string, args []string, password string) error {
	return provider.Set(label, args, password)
}

// Get password from keyring given label and args list.
func Get(args []string) (string, error) {
	return provider.Get(args)
}

// Delete secret from keyring.
func Delete(args []string) error {
	return provider.Delete(args)
}
