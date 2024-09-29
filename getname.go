package main

import (
	"fmt"
	"strings"
)

func GetEncryptedFilePath(filePath string, videoUrl string) string {
	id := (strings.Split(videoUrl[strings.Index(videoUrl, "v=")+2:len(videoUrl)-1], ""))

	for i, j := 0, len(id)-1; i < j; i, j = i+1, j-1 {
		id[i], id[j] = id[j], id[i]
	}

	videoName := strings.Join(id, "")

	finalPath := strings.Replace(filePath, filePath[strings.LastIndex(filePath, "\\"):], fmt.Sprintf("\\%s.enc", videoName), 1)

	return finalPath
}
