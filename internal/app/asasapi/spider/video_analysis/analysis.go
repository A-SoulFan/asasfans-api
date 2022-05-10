package video_analysis

import "go.uber.org/fx"

type Score interface {
	GetKeyType() string
	GetScore(value string) int
}

const (
	Tag = "tag"
	Mid = "mid"
)

type Analysis struct {
	scoreMap map[string]Score
}

func NewAnalysis(scores ...Score) *Analysis {
	analysis := &Analysis{map[string]Score{}}
	for _, score := range scores {
		analysis.scoreMap[score.GetKeyType()] = score
	}

	return analysis
}

func NewFxAnalysis(blacklist *Blacklist, tagScore *TagScore) *Analysis {
	return NewAnalysis(blacklist, tagScore)
}

func (als *Analysis) Calculate(key, value string) int {
	if score, ok := als.scoreMap[key]; ok {
		return score.GetScore(value)
	}
	return 0
}

func Provide() fx.Option {
	return fx.Provide(NewFxAnalysis, NewTagScore, NewBlacklist)
}
