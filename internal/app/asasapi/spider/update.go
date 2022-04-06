package spider

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/repository"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilbil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	failUpdateByListFilename = "fail_update_bv_list.log"
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

	go u.awaitSignal()

	return nil
}

func (u *Update) awaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-c:
		u.logger.Info("receive server signal", zap.String("signal", s.String()))
		if err := u.Stop(); err != nil {
			u.logger.Warn("stop spider server error", zap.Error(err))
		}
		os.Exit(0)
	}
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
	//get time three days ago

	repo := repository.NewBilbilVideo(u.db)
	data, err := repo.Read(time.Now().AddDate(0, 0, -3), time.Now())
	if err != nil {
		return err
	}
	bvList := getUpdateBVList(data)
	for _, bv := range bvList {
		data, err := u.sdk.VideoWebInfo(bv)
		if err != nil {
			u.logger.Error("get video web info error", zap.Error(err))
			continue
		}

		// 处理 tag == title 的情况

		if err := u.updateDB(info, &sInfo); err != nil {
			u.logger.Error("insertDB error", zap.String("bvid", sInfo.Bvid))
			_, _ = failBvFile.WriteString(sInfo.Bvid + "\n")
		}

	}

}

func (u *Update) updateDB(info *bilbil.VideoInfoResponse, sInfo *bilbil.VideoSearchInfo) error {
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
		Tag:       sInfo.Tag,
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

	if err := repository.NewBilbilVideo(u.db.WithContext(context.TODO())).Create(e); err != nil {
		return err
	}

	return nil
}

func getUpdateBVList(videos []*idl.BilbilVideo) []string {
	result := make([]string, 0)
	for _, v := range videos {
		result = append(result, v.Bvid)
	}
	return result
}

//func calculateScore(info *bilbil.VideoInfoResponse) uint64 {
//	score := float64(info.Stat.View)*0.25 +
//		float64(info.Stat.Like+info.Stat.Coin+info.Stat.Reply+info.Stat.Like)*0.4 +
//		float64(info.Stat.Favorite)*0.3 +
//		float64(info.Stat.Share)*0.6
//	return uint64(score)
//}
//
//// isSkip 判断是否需要跳过此条
//func isSkip(sInfo bilbil.VideoSearchInfo, keyword string) bool {
//	tags := strings.Split(sInfo.Tag, ",")
//	// 如果 用户昵称存在 keyword 则进一步检查 tag 是否具有 keyword
//	// 防止错误的收录不属于 keyword 的内容
//	if strings.Index(sInfo.Author, keyword) != -1 {
//		for _, tag := range tags {
//			if tag == keyword {
//				return false
//			}
//		}
//		return true
//	}
//
//	return false
//}
