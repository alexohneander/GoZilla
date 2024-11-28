package handler

import (
	"encoding/base64"
	"fmt"

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
	peer.ID = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", peer.InfoHash, peer.PeerID)))

	// Get IP Address
	if c.Query("ip") != "" {
		peer.IP = c.Query("ip")
	} else {
		peer.IP = c.RemoteIP()
	}

	return peer, nil
}
