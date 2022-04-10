package spider

import (
	"strings"

	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilibili"
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
		if tag == keyword {
			return false
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
