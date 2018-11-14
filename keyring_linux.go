package keyring

import (
	"strings"
	"github.com/godbus/dbus"
	"github.com/divinumDolo/go-keyring/secret_service"
)

type secretServiceProvider struct{}

// Set stores args and pass in the keyring under the defined label
func (s secretServiceProvider) Set(label string, args []string, pass string) error {
	svc, err := ss.NewSecretService()
	if err != nil {
		return err
	}

	// open a session
	session, err := svc.OpenSession()
	if err != nil {
		return err
	}
	defer svc.Close(session)

	var attributes = make(map[string]string)
	for i:=0 ;i < len(args);i=i+2{
		attributes[args[i]] = args[i+1]
	}
	secret := ss.NewSecret(session.Path(), pass)

	collection := svc.GetLoginCollection()

	err = svc.Unlock(collection.Path())
	if err != nil {
		return err
	}

	err = svc.CreateItem(collection, label, attributes, secret)
	if err != nil {
		return err
	}

	return nil
}

// findItem looksup an item by args and user.
//searches in two formats
func (s secretServiceProvider) findItem(svc *ss.SecretService, args []string) (dbus.ObjectPath, error) {
	collection := svc.GetLoginCollection()
	var search = make(map[string]string)
	for i:=0 ;i < len(args);i=i+2{
		search[args[i]] = args[i+1]
	}

	searchOld := map[string]string{
		"username": strings.Join(args," "),
	}
	err := svc.Unlock(collection.Path())
	if err != nil {
		return "", err
	}

	results, err := svc.SearchItems(collection, search)
	results2, err2 := svc.SearchItems(collection, searchOld)

	if err != nil {
		return "", err
	}
	if err2 != nil {
		return "", err2
	}
	if ( len(results) != 0 ){
		return results[0], nil
	}
	if ( len(results2) != 0 ){
		return results2[0], nil
	}
	return "", ErrNotFound
}

// Get gets a secret from the keyring given the args.
func (s secretServiceProvider) Get(args []string) (string, error) {
	svc, err := ss.NewSecretService()
	if err != nil {
		return "", err
	}

	item, err := s.findItem(svc, args)
	if err != nil {
		return "", err
	}

	// open a session
	session, err := svc.OpenSession()
	if err != nil {
		return "", err
	}
	defer svc.Close(session)

	secret, err := svc.GetSecret(item, session.Path())
	if err != nil {
		return "", err
	}

	return string(secret.Value), nil
}

// Delete deletes a secret, identified by args, from the keyring.
func (s secretServiceProvider) Delete(args []string) error {
	svc, err := ss.NewSecretService()
	if err != nil {
		return err
	}

	item, err := s.findItem(svc, args)
	if err != nil {
		return err
	}

	return svc.Delete(item)
}

func init() {
	provider = secretServiceProvider{}
}
