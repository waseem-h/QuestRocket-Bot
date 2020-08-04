package utils

func ChunkInt(array []int, chunkSize int) [][]int {
    var chunks [][]int

    for i := 0; i < len(array); i += chunkSize {
        end := i + chunkSize

        if end > len(array) {
            end = len(array)
        }

        chunks = append(chunks, array[i:end])
    }

    return chunks
}
