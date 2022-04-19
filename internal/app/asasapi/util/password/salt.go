package password

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type PasswordHandler interface {
	Check(pwdHash, input string) bool
	Hash(pwd string) string
}

type Config struct {
	Salt string
}

type DefaultPasswordHandler struct {
	salt string
}

func NewDefaultPasswordHandler(viper *viper.Viper) (PasswordHandler, error) {
	conf := &Config{}
	if err := viper.UnmarshalKey("auth", conf); err != nil {
		return nil, errors.Wrap(err, "unmarshal auth error")
	}

	return &DefaultPasswordHandler{salt: conf.Salt}, nil
}

func (d *DefaultPasswordHandler) Check(pwdHash, input string) bool {
	if pwdHash == "" {
		return false
	}

	return d.Hash(input) == pwdHash
}

func (d *DefaultPasswordHandler) Hash(pwd string) string {
	has := sha1.New()
	has.Write([]byte(d.salt + pwd))
	return hex.EncodeToString(has.Sum(nil))
}
