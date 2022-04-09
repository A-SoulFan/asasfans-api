package spider

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/repository"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilbil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Update struct {
	stopChan chan bool
	db       *gorm.DB
	logger   *zap.Logger
	sdk      *bilbil.SDK
}

func NewUpdate(db *gorm.DB, logger *zap.Logger, sdk *bilbil.SDK) *Update {
	return &Update{
		stopChan: make(chan bool),
		db:       db,
		logger:   logger,
		sdk:      sdk,
	}
}

func (u *Update) Stop(ctx context.Context) error {
	u.logger.Info("stopping spider server")

	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown video update spider server timeout")
		default:
			close(u.stopChan)
			return nil
		}
	}
}

func (u *Update) Run(ctx context.Context) error {
	go func() {
		if err := u.spider(); err != nil {
			u.logger.Fatal("start video update error", zap.Error(err))
		}
	}()

	tk := time.NewTicker(60 * time.Minute)
	go func(_tk *time.Ticker) {
		for {
			select {
			case <-_tk.C:
				u.logger.Info("[tick] video update spider", zap.Time("time", time.Now()))
				if err := u.spider(); err != nil {
					u.logger.Fatal("start video update error", zap.Error(err))
				}
			case <-u.stopChan:
				return
			}
		}
	}(tk)

	return nil
}

func (u *Update) spider() error {
	tx := u.db.WithContext(context.TODO())
	repo := repository.NewBilbilVideo(tx)

	size := 100
	for p := 1; true; p++ {
		list, _, err := repo.FindAllByPubDate(time.Now().Add(-(3 * 24 * time.Hour)), time.Now(), int64(p), int64(size))
		if err != nil {
			u.logger.Error("FindAllByPubDate error", zap.Int("page", p), zap.Error(err))
			return nil
		}

		for _, video := range list {
			time.Sleep(400 * time.Millisecond)
			// 获取视频信息
			vInfo, err := u.sdk.VideoWebInfo(video.Bvid)
			if err != nil {
				u.logger.Error("get video web info error", zap.String("bvid", video.Bvid), zap.Error(err))
				continue
			}

			time.Sleep(200 * time.Millisecond)
			var tags []string
			// 获取视频 tag
			tagInfos, err := u.sdk.VideoWebTagInfo(strconv.Itoa(vInfo.Aid))
			// tag 错误的情况下 不更新 tag
			if err != nil {
				u.logger.Error("get video web tag error", zap.String("bvid", video.Bvid), zap.Error(err))
				tags = strings.Split(video.Tag, ",")
			}
			tags = tagInfos.ToTagStringSlice()

			if err := insertDB(tx, vInfo, strings.Join(tags, ",")); err != nil {
				u.logger.Error("insertDB error", zap.String("bvid", video.Bvid))
			}
		}

		if len(list) < size {
			break
		}
	}

	return nil
}
