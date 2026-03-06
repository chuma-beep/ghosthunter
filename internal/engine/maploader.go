package engine

import (
    "encoding/json"
    "os"
)

type MapData struct {
    Width  int     `json:"width"`
    Height int     `json:"height"`
    Name   int     `json:"name"`
	Tiles  [][]int `json:"tiles"`
}

var LoadedMaps [5][][]int
var MapWidth  [5]int
var MapHeight [5]int
var MapNames  [5]string


func LoadMap(path string, index int) {
    f, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }

    var data MapData
    if err := json.Unmarshal(f, &data); err != nil {
        panic(err)
    }
    
	MapNames[index] = data.Name
    MapWidth[index] = data.Width
    MapHeight[index] = data.Height
    LoadedMaps[index] = data.Tiles
}
