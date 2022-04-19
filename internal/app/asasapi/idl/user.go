package idl

import "github.com/google/uuid"

const (
	UserGenderDefault uint8 = 0
	UserGenderMale    uint8 = 1
	UserGenderFemale  uint8 = 2

	UserStatusDisabled uint8 = 0
	UserStatusNormal   uint8 = 1

	UserHistorySaveMaxNum = 100
)

type User struct {
	Id            uint64    `json:"id"`
	Avatar        uuid.UUID `json:"avatar"`
	Cover         uuid.UUID `json:"cover"`
	Nickname      string    `json:"nickname"`
	Password      string    `json:"-"`
	Status        uint8     `json:"status"`
	Email         string    `json:"email,omitempty"`
	Gender        uint8     `json:"gender"`
	Birthday      *int64    `json:"birthday,omitempty"`
	FollowersNum  uint64    `json:"followers_num"`
	FollowingsNum uint64    `json:"followings_num"`
	ShortDesc     string    `json:"short_desc"`
	CreatedAt     uint64    `json:"created_at"`
	UpdatedAt     uint64    `json:"updated_at"`
}

type UserInfoResp struct {
	Id            uint64    `json:"id"`
	Avatar        uuid.UUID `json:"avatar"`
	Cover         uuid.UUID `json:"cover"`
	Nickname      string    `json:"nickname"`
	Status        uint8     `json:"status"`
	Email         string    `json:"email,omitempty"`
	Gender        uint8     `json:"gender"`
	FollowersNum  uint64    `json:"followers_num"`
	FollowingsNum uint64    `json:"followings_num"`
	Birthday      *int64    `json:"birthday,omitempty"`
	ShortDesc     string    `json:"short_desc"`
	CreatedAt     uint64    `json:"created_at"`
	UpdatedAt     uint64    `json:"updated_at"`
}

type UserRepository interface {
	Create(u *User) (isNew bool, err error)
	Update(u *User, fields ...string) error

	FindUserByID(uid uint64) (*User, error)
	FindUserByEmail(email string) (*User, error)
}
