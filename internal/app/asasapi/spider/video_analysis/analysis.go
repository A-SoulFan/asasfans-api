package video_analysis

import (
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Score interface {
	GetKeyType() string
	GetScore(value string) int
	Reload() error
}

const (
	Tag = "tag"
	Mid = "mid"
)

type Analysis struct {
	logger   *zap.Logger
	scoreMap map[string]Score
}

func NewAnalysis(logger *zap.Logger, scores ...Score) *Analysis {
	analysis := &Analysis{
		logger:   logger,
		scoreMap: map[string]Score{},
	}

	for _, score := range scores {
		analysis.scoreMap[score.GetKeyType()] = score
	}

	go analysis.tick(time.NewTicker(10 * time.Minute))
	return analysis
}

func NewFxAnalysis(logger *zap.Logger, blacklist *Blacklist, tagScore *TagScore) *Analysis {
	return NewAnalysis(logger, blacklist, tagScore)
}

func (als *Analysis) Calculate(key, value string) int {
	if score, ok := als.scoreMap[key]; ok {
		return score.GetScore(value)
	}
	return 0
}

func (als *Analysis) tick(tk *time.Ticker) {
	for {
		select {
		case <-tk.C:
			for k := range als.scoreMap {
				als.logger.Info("reload score start", zap.String("key", k))
				err := als.scoreMap[k].Reload()
				if err != nil {
					als.logger.Error("reload score error", zap.String("key", k), zap.Error(err))
				}
			}
		}
	}
}

func Provide() fx.Option {
	return fx.Provide(NewFxAnalysis, NewTagScore, NewBlacklist)
}
