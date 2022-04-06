package spider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilbil"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBMigrate struct {
	db        *gorm.DB
	sdk       *bilbil.SDK
	logger    *zap.Logger
	isRunning bool
}

func NewDbMigrate(db *gorm.DB, sdk *bilbil.SDK, logger *zap.Logger) *DBMigrate {
	return &DBMigrate{
		db:        db,
		sdk:       sdk,
		logger:    logger,
		isRunning: false,
	}
}

func (m *DBMigrate) Run() {
	go m.run()
}

func (m *DBMigrate) run() {
	// 爬取失败的 bvid 存储在文件
	failBvFile, err := os.OpenFile(failByListFilename, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		m.logger.Fatal(fmt.Sprintf("open %s fail", failByListFilename), zap.Error(err))
		return
	}
	defer failBvFile.Close()

	oldDB, err := gorm.Open(sqlite.Open("ASOUL.db"))

	if err != nil {
		m.logger.Fatal("open sqlite db error", zap.Error(err))
		return
	}

	m.isRunning = true
	size := 100

	type oldInfo struct {
		Title string
		Bvid  string
		Tags  string
	}

	tx := m.db.WithContext(context.TODO())
	for p := 0; m.isRunning; p++ {
		var list []oldInfo
		result := oldDB.Raw("SELECT bvid, tags, title FROM ASOUL_ALL_API LIMIT ?, ?", (p-1)*size, size).Find(&list)
		if result.Error != nil {
			m.logger.Fatal("select ASOUL_ALL_API error", zap.Error(result.Error), zap.Int("page", p))
			return
		}

		for _, info := range list {
			var bInfo *bilbil.VideoInfoResponse
			for retry := 0; retry < 3; retry++ {
				if bInfo, err = m.sdk.VideoWebInfo(info.Bvid); err != nil {
					break
				}
				time.Sleep(time.Duration(retry) * 300 * time.Millisecond)
			}

			if err != nil || bInfo == nil {
				_, _ = failBvFile.WriteString(fmt.Sprintf("%s,%s\n", info.Bvid, err.Error()))
				m.logger.Warn("request VideoWebInfo fail", zap.String("bvid", info.Bvid), zap.Error(err))
				continue
			}

			tags := tagStrToSlice(info.Tags, info.Title)

			if err := insertDB(tx, bInfo, strings.Join(tags, ",")); err != nil {
				_, _ = failBvFile.WriteString(fmt.Sprintf("%s,%s\n", info.Bvid, err.Error()))
				m.logger.Error("insertDB fail", zap.String("bvid", info.Bvid), zap.Error(err))
				continue
			}

			m.logger.Info("insert success", zap.String("bvid", info.Bvid), zap.String("title", info.Title))
		}

		if len(list) < size {
			break
		}
	}
}

func tagStrToSlice(tagStr, title string) []string {
	//if strings.HasPrefix(tagStr, "['") {
	tagStr = strings.Replace(tagStr, "['", "", -1)
	tagStr = strings.Replace(tagStr, "']", "", -1)
	tagStr = strings.Replace(tagStr, title, "", -1)
	//}

	tags := make([]string, 0, 5)
	for _, tag := range strings.Split(tagStr, ",") {
		if len(tag) > 30 {
			continue
		}

		tags = append(tags, tag)
	}

	return tags
}

func (m *DBMigrate) Stop() error {
	m.isRunning = false
	return nil
}
