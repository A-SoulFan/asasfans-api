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
	bilbilVideoTableName    = "bilbil_asoul_video"
	bilbilVideoTagTableName = "bilbil_video_tag"
)

func NewBilbilVideo(tx *gorm.DB) idl.BilbilVideoRepository {
	return &BilbilVideoMysqlImpl{tx: tx}
}

type BilbilVideoMysqlImpl struct {
	tx *gorm.DB
}

func (impl *BilbilVideoMysqlImpl) FindAllByPubDate(from, to time.Time, page, size int64) (list []*idl.BilbilVideo, total int64, err error) {
	result := impl.tx.Table(bilbilVideoTableName).
		Where("pubdate >= ? AND pubdate <= ?", from.Unix(), to.Unix()).
		Offset(int((page - 1) * size)).Limit(int(size)).
		Order("pubdate DESC").
		Find(&list)

	if result == nil {
		return nil, 0, errors.Wrap(result.Error, fmt.Sprintf("select from %s error", bilbilVideoTableName))
	}

	result = impl.tx.Table(bilbilVideoTableName).
		Select("id").
		Where("pubdate >= ? AND pubdate <= ?", from.Second(), to).
		Count(&total)

	if result == nil {
		return nil, 0, errors.Wrap(result.Error, fmt.Sprintf("count from %s error", bilbilVideoTableName))
	}

	return list, total, nil
}

func (impl *BilbilVideoMysqlImpl) Search(queryItems []query_parser.QueryItem, order idl.BilbilVideoOrder, page, size int64) (list []*idl.BilbilVideo, total int64, err error) {
	resp := builderQueryItems(impl.tx, queryItems).Table(bilbilVideoTableName).
		Select(fmt.Sprintf("%s.*", bilbilVideoTableName)).
		Order(fmt.Sprintf("%s DESC", order)).
		Offset(int((page - 1) * size)).Limit(int(size)).
		Find(&list)

	if resp.Error != nil {
		return nil, 0, errors.Wrap(resp.Error, fmt.Sprintf("select from %s error", bilbilVideoTableName))
	}

	resp = builderQueryItems(impl.tx, queryItems).Table(bilbilVideoTableName).
		Select(fmt.Sprintf("%s.id", bilbilVideoTableName)).
		Count(&total)

	if resp.Error != nil {
		return nil, 0, errors.Wrap(resp.Error, fmt.Sprintf("count from %s error", bilbilVideoTableName))
	}

	return list, total, nil
}

func builderQueryItems(tx *gorm.DB, queryItems []query_parser.QueryItem) *gorm.DB {
	for _, item := range queryItems {
		if strings.ToLower(item.Key) == "tag" {
			tx = tx.Where(fmt.Sprintf("%s.id IN (SELECT v_id FROM %s WHERE tag IN (?)) ", bilbilVideoTableName, bilbilVideoTagTableName), item.Values)
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

func (impl *BilbilVideoMysqlImpl) Create(e *idl.BilbilVideo) error {
	return impl.tx.Transaction(func(_tx *gorm.DB) error {
		result := _tx.Table(bilbilVideoTableName).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "bvid"}},
			UpdateAll: true,
		}).Create(&e)

		if result.Error != nil {
			return errors.Wrap(result.Error, fmt.Sprintf("insert %s fail", bilbilVideoTableName))
		}

		tags := strings.Split(e.Tag, ",")
		for _, tag := range tags {
			result = _tx.Table(bilbilVideoTagTableName).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "v_id"}},
				DoNothing: true,
			}).Create(struct {
				Vid uint64 `gorm:"column:v_id"`
				Tag string
			}{Vid: e.Id,
				Tag: tag})

			if result.Error != nil {
				return errors.Wrap(result.Error, fmt.Sprintf("insert %s fail", bilbilVideoTagTableName))
			}
		}

		return nil
	})
}

func (impl *BilbilVideoMysqlImpl) FindAllByBvidList(bvidList []string) (list []*idl.BilbilVideo, err error) {
	result := impl.tx.Table(bilbilVideoTableName).Where("bvid IN (?)", bvidList).Find(&list)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("select from %s error", bilbilVideoTableName))
	}

	return list, nil
}
