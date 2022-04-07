package idl

import "github.com/A-SoulFan/asasfans-api/internal/app/asasapi/util/query_parser"

type BilbilVideoSearchReq struct {
	Order BilbilVideoOrder `form:"order" binding:"required,oneof=pubdate view score"`
	Page  int64            `form:"page,default=1" binding:"omitempty,gt=0"`
	Q     string           `form:"q" binding:"omitempty"`
}

type BilbilVideoSearchResp struct {
	Page       int64          `json:"page"`
	NumResults int64          `json:"numResults"`
	Result     []*BilbilVideo `json:"result"`
}

type BilbilVideo struct {
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
}

type BilbilVideoOrder string

const (
	VideoOrderPubdate = "pubdate"
	VideoOrderView    = "view"
	VideoOrderScore   = "score"
)

type BilbilVideoRepository interface {
	Create(e *BilbilVideo) error

	FindAllByBvidList(bvidList []string) (list []*BilbilVideo, err error)
	Search(queryItems []query_parser.QueryItem, order BilbilVideoOrder, page, size int64) (list []*BilbilVideo, total int64, err error)
}
