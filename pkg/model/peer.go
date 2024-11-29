package model

import (
	"time"
)

type Peer struct {
	ID         string `gorm:"primary_key"`
	InfoHash   string `gorm:"index"`
	PeerID     string `gorm:"index"`
	Category   string `gorm:"index"`
	Port       int32
	Uploaded   int64
	Downloaded int64
	Left       int64
	Compact    bool
	NoPeerID   string
	Event      string
	IP         string
	NumWant    int
	Key        string
	TrackerID  string
	CreatedAt  time.Time
	UpdatedAt  time.Time `gorm:"index"`
}
