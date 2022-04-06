package service

import (
	"context"
	"strings"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/repository"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/util/query_parser"

	"gorm.io/gorm"
)

const (
	defaultQuerySize = 20

	allowKeyItems = "tag,name,mid,"
)

type BilbilVideo struct {
	db *gorm.DB
}

func NewBilbilVideo(db *gorm.DB) *BilbilVideo {
	return &BilbilVideo{db: db}
}

func (b *BilbilVideo) Search(ctx context.Context, req idl.BilbilVideoSearchReq) (*idl.BilbilVideoSearchResp, error) {
	queryItems := query_parser.ParseString(strings.Replace(req.Q, " ", "+", -1), strings.Split(allowKeyItems, ",")...)

	if ok, msg := query_parser.Check(queryItems, queryCheck()); !ok {
		return nil, apperrors.NewError(404, msg)
	}

	tx := b.db.WithContext(ctx)
	bvRepository := repository.NewBilbilVideo(tx)

	list, total, err := bvRepository.Search(queryItems, req.Order, req.Page, defaultQuerySize)
	if err != nil {
		return nil, err
	}

	return &idl.BilbilVideoSearchResp{
		Page:       req.Page,
		NumResults: total,
		Result:     list,
	}, nil
}

func queryCheck() map[string]func(item query_parser.QueryItem) (bool, string) {
	// TODO: check query params
	return map[string]func(item query_parser.QueryItem) (bool, string){}
}
