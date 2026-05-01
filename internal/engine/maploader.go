package engine

import (
	"encoding/json"
)

type MapData struct {
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Name   string  `json:"name"`
	Tiles  [][]int `json:"tiles"`
}

var LoadedMaps [5][][]int
var MapWidth [5]int
var MapHeight [5]int
var MapNames [5]string

func LoadMap(path string, index int) {
	data, err := maps.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var mapData MapData
	if err := json.Unmarshal(data, &mapData); err != nil {
		panic(err)
	}

	MapNames[index] = mapData.Name
	MapWidth[index] = mapData.Width
	MapHeight[index] = mapData.Height
	LoadedMaps[index] = mapData.Tiles
}
