package engine

var ScreenWidth = 640
var ScreenHeight = 400

// func GetMap(index int) *[16][16]int {
    // return &LoadedMaps[index]
//}

func GetMap(index int) [][]int {
    return LoadedMaps[index]
}

func GetMapWidth(index int) int {
    return MapWidth[index]
}

func GetMapHeight(index int) int {
    return MapHeight[index]
}
