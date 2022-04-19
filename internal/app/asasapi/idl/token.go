package idl

type Token struct {
	UserId uint64 `json:"user_id"`
}

type UniqueConflictError struct {
	Message string `json:"message"`
}

func (err UniqueConflictError) Error() string {
	return err.Message
}

type TokenStorage interface {
	Get(key string) (*Token, error)
	Set(key string, token *Token) error
	Del(key string) error
	ClearTokens(token *Token) error
}
