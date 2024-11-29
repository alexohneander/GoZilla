package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alexohneander/GoZilla/internal/database"
	"github.com/alexohneander/GoZilla/pkg/model"
	"github.com/gin-gonic/gin"
)

func Announce(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		c.String(503, "Database Error")
		return
	}

	peer, err := parseAnnounceRequest(c)
	if err != nil {
		c.String(400, "Invalid Request")
		return
	}

	db.Save(&peer)

	c.String(200, "Announce")
}

func parseAnnounceRequest(c *gin.Context) (model.Peer, error) {
	peer := model.Peer{}

	peer.InfoHash = c.Query("info_hash")
	peer.PeerID = c.Query("peer_id")

	if peer.InfoHash == "" || peer.PeerID == "" {
		err := errors.New("invalid request")
		return model.Peer{}, err
	}
	peer.ID = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", peer.InfoHash, peer.PeerID)))

	// Get IP Address
	if c.Query("ip") != "" {
		peer.IP = c.Query("ip")
	} else {
		peer.IP = c.RemoteIP()
	}

	if c.Query("port") == "" {
		err := errors.New("invalid request")
		return model.Peer{}, err
	}

	// parse into int32
	parsedPort, err := strconv.ParseInt(c.Query("port"), 0, 32)
	if err != nil {
		return model.Peer{}, err
	}
	peer.Port = int32(parsedPort)

	if c.Query("uploaded") != "" {
		parsedUploaded, err := strconv.ParseInt(c.Query("uploaded"), 0, 64)
		if err != nil {
			return model.Peer{}, err
		}
		peer.Uploaded = int64(parsedUploaded)
	}

	if c.Query("downloaded") != "" {
		parsedDownloaded, err := strconv.ParseInt(c.Query("downloaded"), 0, 64)
		if err != nil {
			return model.Peer{}, err
		}
		peer.Downloaded = int64(parsedDownloaded)
	}

	if c.Query("left") != "" {
		parsedLeft, err := strconv.ParseInt(c.Query("left"), 0, 64)
		if err != nil {
			return model.Peer{}, err
		}
		peer.Left = int64(parsedLeft)
	}

	peer.UpdatedAt = time.Now()
	return peer, nil
}
