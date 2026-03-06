package engine

import (
    "encoding/json"
    "os"
)

type MapData struct {
    Width  int     `json:"width"`
    Height int     `json:"height"`
    Tiles  [][]int `json:"tiles"`
}

var LoadedMaps [2][][]int
var MapWidth  [2]int
var MapHeight [2]int

func LoadMap(path string, index int) {
    f, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }

    var data MapData
    if err := json.Unmarshal(f, &data); err != nil {
        panic(err)
    }

    MapWidth[index] = data.Width
    MapHeight[index] = data.Height
    LoadedMaps[index] = data.Tiles
}
