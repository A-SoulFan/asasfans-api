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
	RealName      string    `json:"real_name,omitempty"`
	Password      string    `json:"-"`
	MobilePhone   string    `json:"mobile_phone,omitempty"`
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
