package spider

import (
	"strings"

	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilibili"
)

const (
	allowTags = "嘉然,向晚,贝拉,珈乐,乃琳,阿草,嘉心糖,嘉然今天吃什么,A-SOUL,a-soul,传说的世界,向晚大魔王,顶晚人,乃琳Queen,乃淇琳,贝拉kira,贝极星,珈乐Carol,ASOUL,asoul,GNK48,超级敏感,传说的世界,A-SOUL二创激励计划,乃贝,嘉晚饭,果丹皮,琳嘉"
)

func calculateScore(info *bilibili.VideoInfoResponse) uint64 {
	score := float64(info.Stat.View)*0.25 +
		float64(info.Stat.Like+info.Stat.Coin+info.Stat.Reply+info.Stat.Like)*0.4 +
		float64(info.Stat.Favorite)*0.3 +
		float64(info.Stat.Share)*0.6
	return uint64(score)
}

// isSkip 判断是否需要跳过此条
func isSkip(sInfo bilibili.VideoSearchInfo, keyword string) bool {
	tags := strings.Split(sInfo.Tag, ",")
	// 防止错误的收录不属于 keyword 的内容
	for _, tag := range tags {
		for _, allowTag := range strings.Split(allowTags, ",") {
			if tag == allowTag || tag == keyword {
				return false
			}
		}
	}

	return true
}

func tagStrToSlice(tagStr, title string) []string {
	tagStr = strings.Replace(tagStr, "['", "", -1)
	tagStr = strings.Replace(tagStr, "']", "", -1)
	tagStr = strings.Replace(tagStr, title, "", -1)

	tags := make([]string, 0, 5)
	for _, tag := range strings.Split(tagStr, ",") {
		if l := len(tag); l < 1 || l > 30 {
			continue
		}

		tags = append(tags, preTag(tag))
	}

	return tags
}

func preTag(tag string) string {
	for _, s := range []string{"[", "]", "'", "#", " "} {
		tag = strings.Replace(tag, s, "", -1)
	}
	return tag
}
