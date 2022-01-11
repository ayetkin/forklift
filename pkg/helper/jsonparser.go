package helper

import (
	"encoding/json"
	"strings"
)

func JsonToArray(out []byte) ([]string, error) {

	var arr []string
	if err := json.Unmarshal(out, &arr); err != nil{
		return nil, err
	}

	var splitSlashArray []string
	for _, v := range arr {
		delSlash := strings.Split(v, "/")
		splitSlashArray = append(splitSlashArray, delSlash[len(delSlash)-1])
	}

	var items []string
	for _, v := range splitSlashArray {
		if v != "" {
			items = append(items, v)
		}
	}
	return items,nil
}
