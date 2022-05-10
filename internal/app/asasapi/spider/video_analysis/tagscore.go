package video_analysis

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TagScore struct {
	db    *gorm.DB
	tsMap map[string]int
}

func NewTagScore(db *gorm.DB) (*TagScore, error) {
	tagScore := &TagScore{
		db:    db,
		tsMap: map[string]int{},
	}
	if err := tagScore.init(); err != nil {
		return nil, err
	}
	return tagScore, nil
}

func (ts *TagScore) init() error {
	type Storage struct {
		Key   string
		Score int
	}

	var list []*Storage
	result := ts.db.Table("video_analysis").Where("type = ?", Tag).Find(&list)
	if result.Error != nil {
		return errors.Wrap(result.Error, "select video_analysis error")
	}

	tsMap := map[string]int{}
	for _, storage := range list {
		tsMap[storage.Key] = storage.Score
	}

	ts.tsMap = tsMap
	return nil
}

func (ts *TagScore) GetKeyType() string {
	return Tag
}

func (ts *TagScore) GetScore(value string) int {
	if v, ok := ts.tsMap[value]; ok {
		return v
	}
	return 0
}

func (ts *TagScore) Reload() error {
	return ts.init()
}
