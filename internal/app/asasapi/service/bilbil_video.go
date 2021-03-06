package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/repository"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/util/query_parser"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

const (
	defaultQuerySize = 20

	allowKeyItems = "tag,name,mid,view,pubdate"
)

type BilbilVideo struct {
	tnameMaps map[string][]string
	db        *gorm.DB
}

func NewBilbilVideo(db *gorm.DB) *BilbilVideo {
	return &BilbilVideo{tnameMaps: tnameMaps(), db: db}
}

func tnameMaps() map[string][]string {
	return map[string][]string{
		"animation": {"1", "24", "25", "47", "210", "86", "27"},
		"music":     {"3", "28", "31", "30", "194", "59", "193", "29", "130"},
		"dance":     {"20", "198", "199", "200", "154", "156"},
		"game":      {"4", "17", "171", "172", "65", "173", "121", "136", "19"},
		"live":      {"160", "138", "239", "161", "162", "21"},
		"delicacy":  {"211", "76", "212", "213", "214", "215"},
		"guichu":    {"119", "22", "26", "126", "216", "127"},
		//"others":{},
	}
}

func (b *BilbilVideo) Search(ctx context.Context, req idl.BilibiliVideoSearchReq) (*idl.BilibiliVideoSearchResp, error) {
	queryItems := query_parser.ParseString(strings.Replace(req.Q, " ", "+", -1), strings.Split(allowKeyItems, ",")...)

	if ok, msg := query_parser.Check(queryItems, queryCheck()); !ok {
		return nil, apperrors.NewValidationError(404, msg)
	}

	if req.Copyright != 0 {
		queryItems = append(queryItems, query_parser.QueryItem{
			Key:    "copyright",
			Values: []string{strconv.Itoa(req.Copyright)},
			Type:   query_parser.TypeAND,
		})
	}

	if req.Tname != "" {
		if ids, ok := b.tnameMaps[req.Tname]; ok {
			queryItems = append(queryItems, query_parser.QueryItem{
				Key:    "tid",
				Values: ids,
				Type:   query_parser.TypeOR,
			})
		}
	}

	tx := b.db.WithContext(ctx)
	bvRepository := repository.NewBilibiliVideo(tx)

	list, total, err := bvRepository.Search(queryItems, req.Order, req.Page, defaultQuerySize)
	if err != nil {
		return nil, err
	}

	return &idl.BilibiliVideoSearchResp{
		Page:       req.Page,
		NumResults: total,
		Result:     list,
	}, nil
}

var (
	_stringCheck = func(item query_parser.QueryItem) (bool, string) {
		if err := validator.New().Var(item.Values, "omitempty,dive,max=32"); err != nil {
			return false, fmt.Sprintf("%s verification failed.", item.Key)
		}
		return true, "ok"
	}

	_numericCheck = func(item query_parser.QueryItem) (bool, string) {
		if err := validator.New().Var(item.Values, "omitempty,dive,max=64,numeric"); err != nil {
			return false, fmt.Sprintf("%s verification failed.", item.Key)
		}
		return true, "ok"
	}
)

func queryCheck() map[string]func(item query_parser.QueryItem) (bool, string) {
	return map[string]func(item query_parser.QueryItem) (bool, string){
		"tag":     _stringCheck,
		"name":    _stringCheck,
		"mid":     _numericCheck,
		"view":    _numericCheck,
		"pubdate": _numericCheck,
	}
}
