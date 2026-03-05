package engine

import (
    "os"
    "strconv"
    "strings"
)

func LoadHighScore() int {
    data, err := os.ReadFile("highscore.txt")
    if err != nil {
        return 0
    }
    score, err := strconv.Atoi(strings.TrimSpace(string(data)))
    if err != nil {
        return 0
    }
    return score
}

func SaveHighScore(score int) {
    os.WriteFile("highscore.txt", []byte(strconv.Itoa(score)), 0644)
}
