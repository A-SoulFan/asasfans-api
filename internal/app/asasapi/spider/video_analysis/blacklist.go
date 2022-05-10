package video_analysis

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Blacklist struct {
	blacklistMap map[string]int
}

func NewBlacklist(db *gorm.DB) (*Blacklist, error) {
	blacklist := &Blacklist{map[string]int{}}

	type Storage struct {
		Key   string
		Score int
	}

	var list []*Storage
	result := db.Table("video_analysis").Where("type = ?", Mid).Find(&list)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "select video_analysis error")
	}

	for _, storage := range list {
		blacklist.blacklistMap[storage.Key] = storage.Score
	}

	return blacklist, nil
}

func (b *Blacklist) GetKeyType() string {
	return Mid
}

func (b *Blacklist) GetScore(value string) int {
	if v, ok := b.blacklistMap[value]; ok {
		return v
	}
	return 0
}
