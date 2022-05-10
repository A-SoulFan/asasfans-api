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
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/spider/video_analysis"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilibili"
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
	sdk       *bilibili.SDK
	analysis  *video_analysis.Analysis
	isRunning bool
}

func NewVideo(db *gorm.DB, logger *zap.Logger, sdk *bilibili.SDK, analysis *video_analysis.Analysis) *Video {
	return &Video{
		stopChan: make(chan bool),
		db:       db,
		logger:   logger,
		sdk:      sdk,
		analysis: analysis,
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
			return nil
		}
	}
}

func (v *Video) stop() error {
	v.stopChan <- true
	v.isRunning = false
	return nil
}

func (v *Video) Run(ctx context.Context) error {
	tk := time.NewTicker(30 * time.Minute)
	v.isRunning = true

	go func() {
		if err := v.spider(); err != nil {
			v.logger.Error("start spider server error", zap.Error(err))
		}
	}()

	go func(_tk *time.Ticker) {
		for {
			select {
			case <-_tk.C:
				v.logger.Info("[tick] video spider", zap.Time("time", time.Now()))
				if err := v.spider(); err != nil {
					v.logger.Error("start spider server error", zap.Error(err))
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
				v.logger.Error("VideoWebSearchToInfoList error", zap.String("keyword", keyword), zap.Int("page", p), zap.Error(err))
				continue
			}

			for _, sInfo := range list {
				if !v.isRunning {
					return nil
				}

				if isSkip(strings.Split(sInfo.Tag, ","), strconv.Itoa(sInfo.Mid), v.analysis) {
					continue
				}

				time.Sleep(400 * time.Millisecond)
				info, err := v.sdk.VideoWebInfo(sInfo.Bvid)
				if err != nil {
					if bErr, ok := err.(*bilibili.Error); ok {
						v.logger.Warn("VideoWebInfo error", zap.String("bvid", sInfo.Bvid), zap.Int("code", bErr.Code), zap.String("message", bErr.Message))
					} else {
						v.logger.Error("VideoWebInfo error", zap.String("bvid", sInfo.Bvid), zap.Error(err))
					}
					continue
				}

				if err := insertDB(v.db.WithContext(context.TODO()), info, sInfo.Tag); err != nil {
					v.logger.Error("insertDB error", zap.String("bvid", sInfo.Bvid), zap.Error(err))
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

func insertDB(tx *gorm.DB, info *bilibili.VideoInfoResponse, strTag string) error {
	e := &idl.BilibiliVideo{
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
		Status:    idl.BilibiliVideoEnabledStatus,
	}

	if err := repository.NewBilibiliVideo(tx).Save(e); err != nil {
		return err
	}

	return nil
}
