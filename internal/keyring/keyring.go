package keyring

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	ServiceName    = "cogmoteGO"
	ObsKeyringUser = "obs_password"
)

func SaveCredentials(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}
	return keyring.Set(ServiceName, username, password)
}

func GetPassword(username string) (string, error) {
	return keyring.Get(ServiceName, username)
}

func DeleteCredentials(username string) error {
	return keyring.Delete(ServiceName, username)
}

func SaveObsPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	return keyring.Set(ServiceName, ObsKeyringUser, password)
}

func GetObsPassword() (string, error) {
	return keyring.Get(ServiceName, ObsKeyringUser)
}

func DeleteObsPassword() error {
	return keyring.Delete(ServiceName, ObsKeyringUser)
}
