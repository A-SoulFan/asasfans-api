package spider

import (
	"context"
	"fmt"
	"strconv"
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

func (u *Update) Stop() error {
	u.logger.Info("stopping spider server")
	// 平滑关闭,等待5秒钟处理
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := u.stop(ctx); err != nil {
		return errors.Wrap(err, "shutdown spider server error")
	}

	return nil
}

func (u *Update) stop(ctx context.Context) error {
	close(u.stopChan)
	return nil
}

func (u *Update) Run() error {
	go func() {
		err := u.run((*time.Ticker)(time.NewTimer(time.Hour)))
		if err != nil {
			u.logger.Fatal("start spider server error", zap.Error(err))
		}
	}()

	return nil
}

func (u *Update) run(tk *time.Ticker) error {
	if err := u.spider(); err != nil {
		return err
	}

	for {
		select {
		case <-tk.C:
			if err := u.spider(); err != nil {
				return err
			}
		case <-u.stopChan:
			return nil
		}
	}
}

func (u *Update) spider() error {
	repo := repository.NewBilbilVideo(u.db.WithContext(context.TODO()))

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
			// tag 错误的情况下 不更新
			if err != nil {
				u.logger.Error("get video web tag error", zap.String("bvid", video.Bvid), zap.Error(err))
				tags = make([]string, 0)
			}
			tags = tagInfos.ToTagStringSlice()

			// TODO: insertDB
			fmt.Println(tags)
		}

		if len(list) < size {
			break
		}
	}

	return nil
}
