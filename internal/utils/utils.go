package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

type Spinner struct {
	s *spinner.Spinner
}

func NewSpinner(suffix string) *Spinner {
	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond)
	s.Suffix = " " + suffix
	s.Start()
	return &Spinner{s: s}
}

func (sp *Spinner) Stop() {
	sp.s.Stop()
}
