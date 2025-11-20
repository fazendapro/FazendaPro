package models

const (
	Batch1 = 1
	Batch2 = 2
	Batch3 = 3
)

const (
	Batch1MinLiters = 30.0
	Batch2MinLiters = 20.0
	Batch2MaxLiters = 30.0
)

func GetBatchByLiters(liters float64) int {
	if liters > Batch1MinLiters {
		return Batch1
	}
	if liters >= Batch2MinLiters && liters <= Batch2MaxLiters {
		return Batch2
	}
	return Batch3
}
