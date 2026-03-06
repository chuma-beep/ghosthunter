package engine

const ScreenWidth = 320
const ScreenHeight = 200


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
