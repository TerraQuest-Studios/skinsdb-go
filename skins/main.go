package skins

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type SkinData struct {
	Filename string `json:"-"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Author   string `json:"author"`
	License  string `json:"license"`
	Url      string `json:"-"`
	Type     string `json:"type"`
	Img      string `json:"img"`
}

//this exists as an abstraction layer over file system so it can be replaced with a database later

// TODO: replace this with env variable or config file
var skinsDbDataPath = os.Getenv("SKINSDB_DATA_PATH")

func bytesToLines(data []byte) []string {
	rawLines := bytes.Split(data, []byte("\n"))
	lines := make([]string, 0, len(rawLines))
	for _, rawLine := range rawLines {
		lines = append(lines, string(rawLine))
	}
	return lines
}

func getId(skinName string) int {
	parts := strings.Split(skinName, "_")
	if len(parts) > 0 {
		num, _ := strconv.Atoi(parts[1])
		return num
	}
	return 0
}

func getName(skinData []byte) string {
	lines := bytesToLines(skinData)
	if len(lines) > 0 {
		return lines[0]
	}
	return "Unknown Skin"
}

func getAuthor(skinData []byte) string {
	lines := bytesToLines(skinData)
	if len(lines) > 1 {
		// split by newline and take the second line
		return lines[1] // return the second line
	}
	return "Unknown Author"
}

func getLicense(skinData []byte) string {
	lines := bytesToLines(skinData)
	if len(lines) > 2 {
		// split by newline and take the second line
		return lines[2] // return the second line
	}
	return "Unknown License"
}

func getImg(skinName string) string {
	skinImagePath := skinsDbDataPath + "/skins/" + skinName + ".png"
	skinImageData, err := os.ReadFile(skinImagePath)
	if err != nil {
		panic("Error reading skin image file: " + err.Error())
	}
	return base64.StdEncoding.EncodeToString(skinImageData)
}

func GetSkinData() []SkinData {
	//TODO: handle errors better

	entries, err := os.ReadDir(skinsDbDataPath + "/meta")
	if err != nil {
		panic("Error reading skin data: " + err.Error())
	}

	//initiallize skinDataMap to an slice of maps
	skinDataMap := make([]SkinData, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		skinName := entry.Name()
		skinName = skinName[:len(skinName)-4] // remove .txt extension
		skinData, err := os.ReadFile(skinsDbDataPath + "/meta/" + skinName + ".txt")
		if err != nil {
			panic("Error reading skin data file: " + err.Error())
		}
		skin := SkinData{
			Filename: skinName,
			Id:       getId(skinName),
			Name:     getName(skinData),
			Author:   getAuthor(skinData),
			License:  getLicense(skinData),
			Url:      "/skinsdata/skins/" + skinName + ".png",
			Type:     "image/png",
			Img:      getImg(skinName),
		}
		skinDataMap = append(skinDataMap, skin)
	}

	return skinDataMap
}

func GetSkinImage(skinName string) []byte {
	skinImagePath := skinsDbDataPath + "/skins/" + skinName
	skinImageData, err := os.ReadFile(skinImagePath)
	if err != nil {
		return nil
	}
	return skinImageData
}

func IsValidSkinImage(skinImage string) bool {
	skinImagePath := skinsDbDataPath + "/skins/" + skinImage
	_, err := os.Stat(skinImagePath)
	return !os.IsNotExist(err)
}

func GetSkinsSlice(per_page int, page int) []SkinData {
	skins := GetSkinData()
	start := per_page * (page - 1)
	end := start + per_page
	if end > len(skins) {
		end = len(skins)
	}
	if start >= len(skins) {
		return []SkinData{}
	}
	return skins[start:end]
}

func MakeSkinsApiResponse(per_page int, page int) string {
	type SkinsApiResponse struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Page    int        `json:"page"`
		Pages   int        `json:"pages"`
		PerPage int        `json:"per_page"`
		Skins   []SkinData `json:"skins"`
	}

	response := SkinsApiResponse{
		Success: true,
		Message: "who knows",
		Page:    page,
		Pages:   (len(GetSkinData()) + per_page - 1) / per_page, // calculate total pages
		PerPage: per_page,
		Skins:   GetSkinsSlice(per_page, page),
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return `{"success": false, "message": "Error generating response"}`
	}
	return string(responseJson)
}
