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
	"github.com/marksamman/bencode"
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

	peers, err := getPeerListForInfoHash(peer)
	if err != nil {
		c.String(503, "Database Error")
		return
	}

	bencodeDict := writeBencodeDict(peers)

	c.String(200, bencodeDict)
}

func parseAnnounceRequest(c *gin.Context) (model.Peer, error) {
	peer := model.Peer{}

	peer.InfoHash = base64.StdEncoding.EncodeToString([]byte(c.Query("info_hash")))
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

	if c.Query("compact") != "" {
		parsedCompact, err := strconv.ParseBool(c.Query("compact"))
		if err != nil {
			return model.Peer{}, err
		}
		peer.Compact = parsedCompact
	}

	if c.Query("no_peer_id") != "" {
		peer.NoPeerID = c.Query("no_peer_id")
	}

	if c.Query("event") != "" {
		peer.Event = c.Query("event")
	}

	if c.Query("numwant") != "" {
		parsedNumWant, err := strconv.ParseInt(c.Query("numwant"), 0, 32)
		if err != nil {
			return model.Peer{}, err
		}
		peer.NumWant = int32(parsedNumWant)
	}

	if c.Query("key") != "" {
		peer.Key = c.Query("key")
	}

	// if c.Query("tracker") != "" {
	// 	peer.Key = c.Query("key")
	// }

	peer.UpdatedAt = time.Now()
	return peer, nil
}

func getPeerListForInfoHash(peer model.Peer) ([]model.Peer, error) {
	var peers []model.Peer

	db, err := database.GetDB()
	if err != nil {
		return []model.Peer{}, nil
	}

	db.Where("info_hash = ?", peer.InfoHash).Find(&peers)

	return peers, nil
}

func writeBencodeDict(peers []model.Peer) string {
	peersArray := make([]interface{}, len(peers))

	for index, peer := range peers {
		peerDict := make(map[string]interface{})
		if peer.NoPeerID != "" {
			peerDict["id"] = ""
		} else {
			peerDict["id"] = peer.PeerID
		}

		peerDict["ip"] = peer.IP
		peerDict["port"] = peer.Port
		peersArray[index] = peerDict
	}

	dict := make(map[string]interface{})
	dict["interval"] = 60
	dict["peers"] = peersArray

	bencodeDict := bencode.Encode(dict)
	fmt.Printf("bencode encoded dict: %s\n", bencodeDict)

	return string(bencodeDict)
}
