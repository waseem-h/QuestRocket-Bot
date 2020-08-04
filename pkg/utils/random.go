package utils

import "math/rand"

func GetNRandom(min int, max int, count int) IntArray {
    randoms := make(IntArray, 0)
    for i := 0; i < count; i++ {
        random := rand.Intn(max-min) + min
        if randoms.Contains(random) {
            i--
            continue
        }
        randoms = append(randoms, random)
    }

    return randoms
}

type IntArray []int

func (ia IntArray) Contains(number int) bool {
    for _, num := range ia {
        if num == number {
            return true
        }
    }

    return false
}
