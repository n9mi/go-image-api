package model

import "time"

type History struct {
	ID               int       `json:"id"`
	Timestamp        time.Time `json:"timestamp"`
	Type             string    `json:"type"`
	ExtensionBefore  string    `json:"extension_before"`
	ExtensionAfter   string    `json:"extension_after"`
	SizeBeforeInMB   int       `json:"size_before_in_mb"`
	SizeAfterInMB    int       `json:"size_after_in_mb"`
	HeightBeforeInPx int       `json:"height_before_in_px"`
	HeightAfterInPx  int       `json:"height_after_in_px"`
	WidthBeforeInPx  int       `json:"width_before_in_px"`
	WidthAfterInPx   int       `json:"width_after_in_px"`
	IsSuccess        bool      `json:"is_success"`
	ImageLink        string    `json:"image_link"`
}
