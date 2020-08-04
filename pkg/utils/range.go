package utils

func MakeRange(min int, max int) []int {
    array := make([]int, max-min+1)
    for index := range array {
        array[index] = min + index
    }
    return array
}
