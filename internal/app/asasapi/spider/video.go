package spider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/repository"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilbil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	spiderKeywords = "嘉然,向晚,贝拉,珈乐,乃琳,asoul"
	pageDepth      = 10

	failByListFilename = "fail_bv_list.log"
)

type Video struct {
	stopChan  chan bool
	db        *gorm.DB
	logger    *zap.Logger
	sdk       *bilbil.SDK
	isRunning bool
}

func NewVideo(db *gorm.DB, logger *zap.Logger, sdk *bilbil.SDK) *Video {
	return &Video{
		stopChan: make(chan bool),
		db:       db,
		logger:   logger,
		sdk:      sdk,
	}
}

func (v *Video) Stop(ctx context.Context) error {
	v.logger.Info("stopping spider server")

	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown spider server timeout")
		default:
			if err := v.stop(); err != nil {
				return errors.Wrap(err, "shutdown spider server error")
			}
		}
	}
}

func (v *Video) stop() error {
	v.stopChan <- true
	v.isRunning = false
	return nil
}

func (v *Video) Run(ctx context.Context) error {
	tk := time.NewTimer(time.Hour)
	v.isRunning = true

	go func() {
		if err := v.spider(); err != nil {
			v.logger.Fatal("start spider server error", zap.Error(err))
		}
	}()

	go func(_tk *time.Timer) {
		for {
			select {
			case <-tk.C:
				if err := v.spider(); err != nil {
					v.logger.Fatal("start spider server error", zap.Error(err))
				}
			case <-v.stopChan:
				return
			}
		}
	}(tk)

	return nil
}

func (v *Video) spider() error {
	// 爬取失败的 bvid 存储在文件
	failBvFile, err := os.OpenFile(failByListFilename, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("open %s fail", failByListFilename))
	}
	defer failBvFile.Close()

	// 从指定的 keyword 中进行 搜索
	for _, keyword := range strings.Split(spiderKeywords, ",") {
		for p := 1; p <= pageDepth; p++ {
			list, totalPage, err := v.sdk.VideoWebSearchToInfoList(keyword, p)
			if err != nil {
				v.logger.Error("VideoWebSearchToInfoList error", zap.String("keyword", keyword), zap.Int("page", p))
				continue
			}

			for _, sInfo := range list {
				if !v.isRunning {
					return nil
				}

				if isSkip(sInfo, keyword) {
					continue
				}

				time.Sleep(400 * time.Millisecond)
				info, err := v.sdk.VideoWebInfo(sInfo.Bvid)
				if err != nil {
					v.logger.Error("VideoWebInfo error", zap.String("bvid", sInfo.Bvid))
					continue
				}

				// 处理 tag == title 的情况
				tags := make([]string, 0, 5)
				for _, tag := range strings.Split(sInfo.Tag, ",") {
					if tag != sInfo.Title {
						tags = append(tags, tag)
					}
				}
				sInfo.Tag = strings.Join(tags, ",")

				if err := insertDB(v.db.WithContext(context.TODO()), info, sInfo.Tag); err != nil {
					v.logger.Error("insertDB error", zap.String("bvid", sInfo.Bvid))
					_, _ = failBvFile.WriteString(sInfo.Bvid + "\n")
				}

				v.logger.Info("insert success", zap.String("bvid", info.Bvid), zap.String("title", info.Title), zap.String("tag", sInfo.Tag))
			}

			if p >= totalPage {
				break
			}
		}
	}

	return nil
}

func insertDB(tx *gorm.DB, info *bilbil.VideoInfoResponse, strTag string) error {
	e := &idl.BilbilVideo{
		Bvid:      info.Bvid,
		Aid:       uint64(info.Aid),
		Name:      info.Owner.Name,
		Mid:       uint64(info.Owner.Mid),
		Face:      info.Owner.Face,
		Tid:       uint64(info.Tid),
		Tname:     info.Tname,
		Copyright: uint64(info.Copyright),
		Title:     info.Title,
		Desc:      info.Desc,
		Pic:       info.Pic,
		Tag:       strTag,
		Pubdate:   uint64(info.Pubdate),
		Duration:  strconv.Itoa(info.Duration),
		View:      uint64(info.Stat.View),
		Danmaku:   uint64(info.Duration),
		Reply:     uint64(info.Stat.Reply),
		Favorite:  uint64(info.Stat.Favorite),
		Coin:      uint64(info.Stat.Coin),
		Share:     uint64(info.Stat.Share),
		Like:      uint64(info.Stat.Like),
		Score:     calculateScore(info),
	}

	if err := repository.NewBilbilVideo(tx).Create(e); err != nil {
		return err
	}

	return nil
}

func calculateScore(info *bilbil.VideoInfoResponse) uint64 {
	score := float64(info.Stat.View)*0.25 +
		float64(info.Stat.Like+info.Stat.Coin+info.Stat.Reply+info.Stat.Like)*0.4 +
		float64(info.Stat.Favorite)*0.3 +
		float64(info.Stat.Share)*0.6
	return uint64(score)
}

// isSkip 判断是否需要跳过此条
func isSkip(sInfo bilbil.VideoSearchInfo, keyword string) bool {
	tags := strings.Split(sInfo.Tag, ",")
	// 防止错误的收录不属于 keyword 的内容
	for _, tag := range tags {
		if tag == keyword {
			return false
		}
	}

	return true
}
