package idl

import (
	"time"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/util/query_parser"
)

type BilibiliVideoSearchReq struct {
	Order     BilibiliVideoOrder `form:"order" binding:"required,oneof=pubdate view score"`
	Page      int64              `form:"page,default=1" binding:"omitempty,gt=0"`
	Q         string             `form:"q" binding:"omitempty"`
	Copyright int                `form:"copyright" binding:"omitempty,oneof=1 2"`
	Tname     string             `form:"tname" binding:"omitempty,oneof=animation music dance game live delicacy guichu"`
}

type BilibiliVideoSearchResp struct {
	Page       int64            `json:"page"`
	NumResults int64            `json:"numResults"`
	Result     []*BilibiliVideo `json:"result"`
}

const (
	BilibiliVideoEnabledStatus  = 1
	BilibiliVideoDisabledStatus = 0
)

type BilibiliVideo struct {
	Id        uint64 `json:"-"`
	Bvid      string `json:"bvid"`
	Aid       uint64 `json:"aid"`
	Name      string `json:"name"`
	Mid       uint64 `json:"mid"`
	Face      string `json:"face"`
	Tid       uint64 `json:"tid"`
	Tname     string `json:"tname"`
	Copyright uint64 `json:"copyright"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Pic       string `json:"pic"`
	Tag       string `json:"tag"`
	Pubdate   uint64 `json:"pubdate"`
	Duration  string `json:"duration"`
	View      uint64 `json:"view"`
	Danmaku   uint64 `json:"danmaku"`
	Reply     uint64 `json:"reply"`
	Favorite  uint64 `json:"favorite"`
	Coin      uint64 `json:"coin"`
	Share     uint64 `json:"share"`
	Like      uint64 `json:"like"`
	Score     uint64 `json:"score"`
	Status    uint8  `json:"status"`
	CreatedAt uint64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt uint64 `json:"updated_at" gorm:"autoUpdateTime"`
}

type BilibiliVideoOrder string

const (
	VideoOrderPubdate = "pubdate"
	VideoOrderView    = "view"
	VideoOrderScore   = "score"
)

type BilibiliVideoRepository interface {
	Save(e *BilibiliVideo) error

	FindAllByBvidList(bvidList []string) (list []*BilibiliVideo, err error)
	FindAllByPubDate(from time.Time, to time.Time, page, size int64) (list []*BilibiliVideo, total int64, err error)

	Search(queryItems []query_parser.QueryItem, order BilibiliVideoOrder, page, size int64) (list []*BilibiliVideo, total int64, err error)
}
