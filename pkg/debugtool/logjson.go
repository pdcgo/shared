package debugtool

import (
	"encoding/json"
	"log"
)

func LogJson(v ...any) {
	for _, item := range v {
		data, _ := json.MarshalIndent(item, "", "  ")
		log.Println(string(data))
	}

}
