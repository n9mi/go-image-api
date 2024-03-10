package entity

import "time"

type History struct {
	ID               int `gorm:"primaryKey"`
	Timestamp        time.Time
	Type             string
	ExtensionBefore  string
	ExtensionAfter   string
	SizeBeforeInMB   float64
	SizeAfterInMB    float64
	HeightBeforeInPx int
	HeightAfterInPx  int
	WidthBeforeInPx  int
	WidthAfterInPx   int
	ImageLinkBefore  string
	ImageLinkAfter   string
}
