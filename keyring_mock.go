package keyring
import "strings"
type mockProvider struct {
	mockStore map[string]string
}

// Set stores args and pass in the keyring under the defined label
// name.
func (m *mockProvider) Set(label string, args []string, pass string) error {
	if m.mockStore == nil {
		m.mockStore = make(map[string]string)
	}
	
	keyValuePairs := strings.Join(args, " ")
	m.mockStore[keyValuePairs] = pass
	return nil
}

// Get gets a secret from the keyring given  args.
func (m *mockProvider) Get(args []string) (string, error) {
	keyValuePairs := strings.Join(args, " ")
	if b, ok := m.mockStore[keyValuePairs]; ok {
			return b, nil
	}
	return "", ErrNotFound
}

// Delete deletes a secret, identified by args, from the keyring.
func (m *mockProvider) Delete(args []string) error {
	keyValuePairs := strings.Join(args, " ")
	if m.mockStore != nil {
		if _, ok := m.mockStore[keyValuePairs]; ok {
				delete(m.mockStore,keyValuePairs)
				return nil
		}
	}
	return ErrNotFound
}

// MockInit sets the provider to a mocked memory store
func MockInit() {
	provider = &mockProvider{}
}