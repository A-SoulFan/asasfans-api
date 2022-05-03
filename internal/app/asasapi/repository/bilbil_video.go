package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/util/query_parser"
	"gorm.io/gorm/clause"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	bilibiliVideoTableName    = "bilbil_asoul_video"
	bilibiliVideoTagTableName = "bilbil_video_tag"

	preAliasString = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	onUpdateFields = "name,mid,face,tid,tname,copyright,title,desc,pic,tag,pubdate,duration,view,danmaku,reply,favorite,coin,share,like,score,updated_at"
)

func NewBilibiliVideo(tx *gorm.DB) idl.BilibiliVideoRepository {
	return &BilibiliVideoMysqlImpl{tx: tx}
}

type BilibiliVideoMysqlImpl struct {
	tx *gorm.DB
}

func (impl *BilibiliVideoMysqlImpl) FindAllByPubDate(from, to time.Time, page, size int64) (list []*idl.BilibiliVideo, total int64, err error) {
	result := impl.tx.Table(bilibiliVideoTableName).
		Where("pubdate >= ? AND pubdate <= ?", from.Unix(), to.Unix()).
		Offset(int((page - 1) * size)).Limit(int(size)).
		Order("pubdate DESC").
		Find(&list)

	if result == nil {
		return nil, 0, errors.Wrap(result.Error, fmt.Sprintf("select from %s error", bilibiliVideoTableName))
	}

	result = impl.tx.Table(bilibiliVideoTableName).
		Select("id").
		Where("pubdate >= ? AND pubdate <= ?", from.Second(), to).
		Count(&total)

	if result == nil {
		return nil, 0, errors.Wrap(result.Error, fmt.Sprintf("count from %s error", bilibiliVideoTableName))
	}

	return list, total, nil
}

func (impl *BilibiliVideoMysqlImpl) Search(queryItems []query_parser.QueryItem, order idl.BilibiliVideoOrder, page, size int64) (list []*idl.BilibiliVideo, total int64, err error) {
	resp := builderQueryItems(impl.tx, queryItems).Table(bilibiliVideoTableName).
		Where(fmt.Sprintf("%s.status = ?", bilibiliVideoTableName), idl.BilibiliVideoEnabledStatus).
		Select(fmt.Sprintf("%s.*", bilibiliVideoTableName)).
		Order(fmt.Sprintf("%s DESC", order)).
		Offset(int((page - 1) * size)).Limit(int(size)).
		Find(&list)

	if resp.Error != nil {
		return nil, 0, errors.Wrap(resp.Error, fmt.Sprintf("select from %s error", bilibiliVideoTableName))
	}

	resp = builderQueryItems(impl.tx, queryItems).Table(bilibiliVideoTableName).
		Where(fmt.Sprintf("%s.status = ?", bilibiliVideoTableName), idl.BilibiliVideoEnabledStatus).
		Select(fmt.Sprintf("%s.id", bilibiliVideoTableName)).
		Count(&total)

	if resp.Error != nil {
		return nil, 0, errors.Wrap(resp.Error, fmt.Sprintf("count from %s error", bilibiliVideoTableName))
	}

	return list, total, nil
}

func builderQueryItems(tx *gorm.DB, queryItems []query_parser.QueryItem) *gorm.DB {
	for _, item := range queryItems {
		if strings.ToLower(item.Key) == "tag" {
			switch item.Type {
			case query_parser.TypeAND:
				preSQL := ""
				values := make([]interface{}, 0, 5)
				for i := 0; i < len(item.Values) && i < 5; i++ {
					alias := string(preAliasString[i])
					tmpStr := fmt.Sprintf("(SELECT tag, v_id FROM %s WHERE tag = ?)", bilibiliVideoTagTableName)
					values = append(values, item.Values[i])
					if i == 0 {
						preSQL += fmt.Sprintf("FROM %s AS %s", tmpStr, alias)
					} else {
						preSQL += fmt.Sprintf(" JOIN %s AS %s ON %s.v_id = %s.v_id", tmpStr, alias, alias, string(preAliasString[i-1]))
					}
				}

				tx = tx.Where(fmt.Sprintf("%s.id IN (SELECT %s.v_id %s)", bilibiliVideoTableName, string(preAliasString[0]), preSQL), values...)
			case query_parser.TypeOR:
				tx = tx.Where(fmt.Sprintf("%s.id IN (SELECT v_id FROM %s WHERE tag IN (?)) ", bilibiliVideoTableName, bilibiliVideoTagTableName), item.Values)
			}

			continue
		}

		switch item.Type {
		case query_parser.TypeAND:
			for _, value := range item.Values {
				tx = tx.Where(fmt.Sprintf("%s = ?", item.Key), value)
			}
		case query_parser.TypeOR:
			tx = tx.Where(fmt.Sprintf("%s IN (?)", item.Key), item.Values)
		case query_parser.TypeBetween:
			v := item.GetBetweenValues()
			tx = tx.Where(fmt.Sprintf("%s BETWEEN ? AND ?", item.Key), v[0], v[1])
		}
	}

	return tx
}

func (impl *BilibiliVideoMysqlImpl) Save(e *idl.BilibiliVideo) error {
	return impl.tx.Transaction(func(_tx *gorm.DB) error {
		result := _tx.Table(bilibiliVideoTableName).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "bvid"}},
			DoUpdates: clause.AssignmentColumns(strings.Split(onUpdateFields, ",")),
		}).Create(&e)

		if result.Error != nil {
			return errors.Wrap(result.Error, fmt.Sprintf("insert %s fail", bilibiliVideoTableName))
		}

		// on update
		if e.Id == 0 {
			_tx.Table(bilibiliVideoTableName).Where("bvid = ?", e.Bvid).Select("id").Find(&e.Id)
		}

		tags := strings.Split(e.Tag, ",")
		for _, tag := range tags {
			result = _tx.Table(bilibiliVideoTagTableName).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "v_id"}},
				DoNothing: true,
			}).Create(struct {
				Vid uint64 `gorm:"column:v_id"`
				Tag string
			}{Vid: e.Id,
				Tag: tag})

			if result.Error != nil {
				return errors.Wrap(result.Error, fmt.Sprintf("insert %s fail", bilibiliVideoTagTableName))
			}
		}

		return nil
	})
}

func (impl *BilibiliVideoMysqlImpl) FindAllByBvidList(bvidList []string) (list []*idl.BilibiliVideo, err error) {
	result := impl.tx.Table(bilibiliVideoTableName).Where("bvid IN (?)", bvidList).Find(&list)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("select from %s error", bilibiliVideoTableName))
	}

	return list, nil
}
