package badger

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	log "git.randomchars.net/FreeNitori/Log"
	"os"
	"time"
)

var backupTicker *time.Ticker

func setupBackup(db *Badger) error {
	if config.System.BackupInterval > 0 {
		log.Infof("Periodical database backup interval %v second(s).", config.System.BackupInterval)
		if _, err := os.Stat("backup"); os.IsNotExist(err) {
			err = os.Mkdir("backup", 0700)
			if err != nil {
				return err
			}
		}
		backupTicker = time.NewTicker(time.Duration(config.System.BackupInterval) * time.Second)
		go func() {
			for {
				select {
				case <-backupTicker.C:
					intermediate, err := os.Create("backup/.intermediate")
					if err != nil {
						log.Errorf("Error creating intermediate backup file, %s", err)
						continue
					}
					ver, err := db.DB.Backup(intermediate, 0)
					if err != nil {
						log.Errorf("Error generating intermediate backup, %s", err)
						continue
					}
					err = os.Rename("backup/.intermediate", fmt.Sprintf("backup/%v", ver))
					if err != nil {
						log.Errorf("Error renaming intermediate backup, %s", err)
						err = os.Remove("backup/.intermediate")
						if err != nil {
							log.Errorf("Error removing intermediate file, %s", err)
							log.Warnf("Backup has been disabled, please resolve the issue above and restart FreeNitori.")
							break
						}
						continue
					}
					log.Infof("Successfully backed up database, version %v", ver)
				}
			}
		}()
		return nil
	} else {
		log.Infof("Periodical database backup is not enabled.")
		return nil
	}
}
