package utils

import (
	"fmt"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func CreateProgressBar(steps int64, text string) (*mpb.Progress, *mpb.Bar) {
	barWidth := 64
	p := mpb.New(mpb.WithWidth(barWidth))

	return p, p.New(steps,
		mpb.BarStyle().Lbound("[").Filler("=").Tip(">").Padding("-").Rbound("]"),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf(text)),
			decor.CountersNoUnit("%d/%d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)
}
