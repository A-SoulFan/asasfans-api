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

func (m *DBMigrate) Run(ctx context.Context) error {
	go m.run()
	return nil
}

func (m *DBMigrate) run() {
	// 爬取失败的 bvid 存储在文件
	failBvFile, err := os.OpenFile(failByListFilename, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		m.logger.Fatal(fmt.Sprintf("open %s fail", failByListFilename), zap.Error(err))
		return
	}
	defer failBvFile.Close()

	oldDB, err := gorm.Open(sqlite.Open("./ASOUL.db"))

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
	for p := 1; m.isRunning; p++ {
		m.logger.Info("start new page ", zap.Int("page", p))

		var list []oldInfo
		result := oldDB.Raw("SELECT bvid, tags, title FROM ASOUL_ALL_API LIMIT ?, ?", (p-1)*size, size).Find(&list)
		if result.Error != nil {
			m.logger.Fatal("select ASOUL_ALL_API error", zap.Error(result.Error), zap.Int("page", p))
			return
		}

		bvidList := make([]string, 0, size)
		for _, info := range list {
			bvidList = append(bvidList, info.Bvid)
		}

		var inList []*struct {
			Bvid string
		}

		oList := make([]oldInfo, 0, 100)
		if tx.Table("bilbil_asoul_video").Where("bvid IN (?)", bvidList).Select("bvid").Find(&inList).Error == nil {
			for _, info := range list {
				var hit bool
				for _, e := range inList {
					if e.Bvid == info.Bvid {
						hit = true
						break
					}
				}

				if !hit {
					oList = append(oList, info)
				}
			}
		}

		for _, info := range oList {
			if !m.isRunning {
				break
			}

			var bInfo *bilbil.VideoInfoResponse
			time.Sleep(400 * time.Millisecond)
			if bInfo, err = m.sdk.VideoWebInfo(info.Bvid); err != nil || bInfo == nil {
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

	m.logger.Info("db migrate success")
}

func (m *DBMigrate) Stop(ctx context.Context) error {
	m.isRunning = false
	return nil
}
