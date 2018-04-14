package ramstore

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"time"

	"github.com/fe0b6/config"
)

func saveData() {

	fh, err := os.OpenFile(config.GetStr("path", "db"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("[fatal]", err)
		return
	}
	defer fh.Close()

	gh := gob.NewEncoder(fh)

	for i := range data {
		data[i].Lock()

		for k, v := range data[i].Data {
			if v.Deleted {
				if time.Now().After(time.Unix(0, v.Time).Add(deletedTimeout * time.Hour)) {
					continue
				}
			}

			gh.Encode(saveObj{Key: k, Obj: v})
		}

		data[i].Unlock()
	}

}

func readData() {
	fh, err := os.Open(config.GetStr("path", "db"))
	if err != nil {
		return
	}
	defer fh.Close()

	gh := gob.NewDecoder(fh)

	for {
		var d saveObj
		err := gh.Decode(&d)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("[error]", err)
			return
		}

		Set(d.Key, d.Obj)
	}

}
