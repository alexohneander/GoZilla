package task

import (
	"fmt"
	"time"

	"github.com/alexohneander/GoZilla/internal/database"
	"github.com/alexohneander/GoZilla/pkg/model"
)

func CleanPeers() {
	db, err := database.GetDB()
	if err != nil {
		fmt.Println(err)
	}

	for {
		thirtyMinutesAgo := time.Now().Add(-30 * time.Minute)

		var peers []model.Peer
		result := db.Where("updated_at < ?", thirtyMinutesAgo).Find(&peers)

		if result.Error != nil {
			fmt.Println("Fehler beim Abrufen der Peers:", result.Error)
			return
		}

		for _, peer := range peers {
			fmt.Println("removed dead peers:", len(peers))
			db.Delete(peer)
		}

		time.Sleep(10 * time.Second)
	}

}
