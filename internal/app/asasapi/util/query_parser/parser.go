package query_parser

import (
	"strings"
)

const (
	TypeOR      = "OR"
	TypeAND     = "AND"
	TypeBetween = "BETWEEN"
)

type QueryItem struct {
	Key    string
	Values []string
	Type   string
}

func (q *QueryItem) GetBetweenValues() []string {
	return []string{q.Values[0], q.Values[1]}
}

// ParseString 将查询字符串解析为 []QueryItem
//
// 首先使用 ~ 进行分割 产生四组 query string
// 然后使用 . 进行分割 将 query string 转化成 query slice, 此时 key = query[0], 如果 query[len(query) - 1] 是 Type 关键字 (OR, AND, BETWEEN) 则赋值为 Type 否则 默认使用 OR
// 最后使用 + 分割 values, 注意如果 Type = BETWEEN 则 len(values) == 2 否则记为无效, 同时由于 + 是 http query 的保留关键词,用于替换 空格 的场景, 需要上层自行处理
func ParseString(s string, allowKeyItem ...string) []QueryItem {
	if len(s) == 0 {
		return []QueryItem{}
	}

	list := make([]QueryItem, 0, 3)

	for _, itemString := range strings.Split(s, "~") {
		keyStrings := strings.Split(itemString, ".")
		// 无效的
		if len(keyStrings) < 2 {
			continue
		}

		if !checkKeyItem(keyStrings[0], allowKeyItem) {
			continue
		}

		qItem := QueryItem{
			Key:    keyStrings[0],
			Values: strings.Split(keyStrings[1], "+"),
		}

		// 无效的
		if len(qItem.Values) < 1 {
			continue
		}

		if t := strings.ToUpper(keyStrings[len(keyStrings)-1]); t == TypeOR || t == TypeAND {
			qItem.Type = t
		} else if t == TypeBetween {
			// 无效的
			if len(qItem.Values) != 2 {
				continue
			}
			qItem.Type = t
		} else {
			qItem.Type = TypeOR
		}

		list = append(list, qItem)
	}

	return list
}

func Check(qItems []QueryItem, funcMap map[string]func(item QueryItem) (bool, string)) (bool, string) {
	for _, item := range qItems {
		if f, ok := funcMap[item.Key]; ok {
			if ok, msg := f(item); !ok {
				return ok, msg
			}
		}
	}

	return true, ""
}

func checkKeyItem(k string, allowKeyItems []string) bool {
	for _, allowKey := range allowKeyItems {
		if k == allowKey {
			return true
		}
	}

	return false
}
