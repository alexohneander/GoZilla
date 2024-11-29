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
		result := db.Find(&peers).Where("updated_at < ?", thirtyMinutesAgo)

		if result.Error != nil {
			fmt.Println("Fehler beim Abrufen der Peers:", result.Error)
			return
		}

		fmt.Println("Peers älter als 30 Minuten:", len(peers))
		time.Sleep(10 * time.Second)
	}

}
