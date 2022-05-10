package video_analysis

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TagScore struct {
	tsMap map[string]int
}

func NewTagScore(db *gorm.DB) (*TagScore, error) {
	tagScore := &TagScore{tsMap: map[string]int{}}

	type Storage struct {
		Key   string
		Score int
	}

	var list []*Storage
	result := db.Table("video_analysis").Where("type = ?", Tag).Find(&list)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "select video_analysis error")
	}

	for _, storage := range list {
		tagScore.tsMap[storage.Key] = storage.Score
	}

	return tagScore, nil
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
