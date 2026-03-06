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

var LoadedMaps [2][16][16]int

func LoadMap(path string, index int) {
    f, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }

    var data MapData
    if err := json.Unmarshal(f, &data); err != nil {
        panic(err)
    }

    for row := 0; row < 16; row++ {
        for col := 0; col < 16; col++ {
            LoadedMaps[index][row][col] = data.Tiles[row][col]
        }
    }
}
