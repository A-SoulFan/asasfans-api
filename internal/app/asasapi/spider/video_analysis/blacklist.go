package video_analysis

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Blacklist struct {
	db           *gorm.DB
	blacklistMap map[string]int
}

func NewBlacklist(db *gorm.DB) (*Blacklist, error) {
	blacklist := &Blacklist{
		db:           db,
		blacklistMap: map[string]int{},
	}

	return blacklist, nil
}

func (b *Blacklist) init() error {
	type Storage struct {
		Key   string
		Score int
	}

	var list []*Storage
	result := b.db.Table("video_analysis").Where("type = ?", Mid).Find(&list)
	if result.Error != nil {
		return errors.Wrap(result.Error, "select video_analysis error")
	}

	blacklistMap := map[string]int{}
	for _, storage := range list {
		blacklistMap[storage.Key] = storage.Score
	}

	b.blacklistMap = blacklistMap
	return nil
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

func (b *Blacklist) Reload() error {
	return b.init()
}
