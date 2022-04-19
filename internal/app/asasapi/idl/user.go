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

type UserUpdateReq struct {
	Avatar    string  `json:"avatar" binding:"omitempty,uuid"`
	Cover     string  `json:"cover" binding:"omitempty,uuid"`
	Nickname  string  `json:"nickname" binding:"omitempty,min=2,max=16"`
	Gender    *uint8  `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Birthday  *int64  `json:"birthday" binding:"omitempty,numeric"`
	ShortDesc *string `json:"short_desc" binding:"omitempty,max=500"`
}

func NewUserInfoResp(u *User) *UserInfoResp {
	return &UserInfoResp{
		Id:            u.Id,
		Avatar:        u.Avatar,
		Cover:         u.Cover,
		Nickname:      u.Nickname,
		Status:        u.Status,
		Email:         u.Email,
		Gender:        u.Gender,
		FollowersNum:  u.FollowersNum,
		FollowingsNum: u.FollowingsNum,
		Birthday:      u.Birthday,
		ShortDesc:     u.ShortDesc,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

type UserRepository interface {
	Create(u *User) (isNew bool, err error)
	Update(u *User, fields ...string) error

	FindUserByID(uid uint64) (*User, error)
	FindUserByEmail(email string) (*User, error)
}
