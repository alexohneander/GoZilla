package handler

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/alexohneander/GoZilla/internal/database"
	"github.com/alexohneander/GoZilla/internal/helper"
	"github.com/alexohneander/GoZilla/pkg/model"
	"github.com/gin-gonic/gin"
	bencode "github.com/jackpal/bencode-go"
)

var (
	announceInterval    = 1800
	minAnnounceInterval = 900
	maxAccounceDrift    = 300
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
		c.Writer.Header().Set("Content-Type", "text/plain")
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
	} else if c.Query("ipv4") != "" {
		peer.IP = c.Query("ipv4")
	} else if c.Query("ipv6") != "" {
		peer.IP = c.Query("ipv6")
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

		/* We want to check if the Client is a Seeder or a Leecher.
		We can do this by checking what the client needs to download. Is the left > 0 then we know the client is a Leecher.
		Otherwise the client is a Seeder. */
		if parsedLeft > 0 {
			peer.Category = "Leecher"
		} else if c.Query("event") == "completed" {
			peer.Category = "Seeder"
		} else if parsedLeft == 0 {
			peer.Category = "Seeder"
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
		peer.NumWant = int(parsedNumWant)
	}

	if c.Query("key") != "" {
		peer.Key = c.Query("key")
	}

	peer.UpdatedAt = time.Now()
	return peer, nil
}

func getPeerListForInfoHash(peer model.Peer) ([]model.Peer, error) {
	var peers []model.Peer

	db, err := database.GetDB()
	if err != nil {
		return []model.Peer{}, nil
	}

	/* If the NumWant is greater than zero then we want to return that many peers.
	If the NumWant is less than or equal to zero then we want to return our default limit (50 peers). */
	var limit int = 50
	if peer.NumWant > 0 {
		limit = int(peer.NumWant)
	}

	/* If the client is a Leecher then we want to return the peers that are Seeders.
	If the client is a Seeder then we want to return the peers that are Leechers.
	Otherwise we returning all peers that are Leechers or Seeders. */
	if peer.Category == "Leecher" {
		db.Where("info_hash = ?", peer.InfoHash).Where("category = ?", "Seeder").Limit(limit).Find(&peers)
	} else if peer.Category == "Seeder" {
		db.Where("info_hash =?", peer.InfoHash).Where("category = ?", "Leecher").Limit(limit).Find(&peers)
	} else {
		db.Where("info_hash =?", peer.InfoHash).Limit(limit).Find(&peers)
	}

	return peers, nil
}

func writeBencodeDict(peers []model.Peer) string {
	peersArray := make([]interface{}, len(peers))

	for index, peer := range peers {
		peerDict := make(map[string]interface{})
		if peer.NoPeerID != "" {
			peerDict["peer id"] = ""
		} else {
			peerDict["peer id"] = peer.PeerID
		}

		peerDict["ip"] = peer.IP
		peerDict["port"] = peer.Port
		peersArray[index] = peerDict
	}

	dict := make(map[string]interface{})

	/* We ask clients to announce each interval seconds. In order to spread the load on tracker,
	we will vary the interval given to client by random number of seconds between 0 and value
	specified in the Variables */
	dict["interval"] = announceInterval + helper.UnsafeIntn(maxAccounceDrift)
	dict["min interval"] = minAnnounceInterval

	dict["peers"] = peersArray

	var bencodedDict bytes.Buffer
	err := bencode.Marshal(&bencodedDict, dict)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("bencode encoded dict: %s\n", bencodedDict.String())

	return string(bencodedDict.String())
}
