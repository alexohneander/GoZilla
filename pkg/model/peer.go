package model

import (
	"time"

	"gorm.io/gorm"
)

type Peer struct {
	ID         string `gorm:"primary_key"`
	InfoHash   string `gorm:"index"`
	PeerID     string `gorm:"index"`
	Port       int32
	Uploaded   int64
	Downloaded int64
	Left       int64
	Compact    bool
	NoPeerID   string
	Event      string
	IP         string
	NumWant    int32
	Key        string
	TrackerID  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
