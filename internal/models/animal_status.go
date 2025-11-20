package models

const (
	AnimalStatusActive   = 0
	AnimalStatusSold     = 1
	AnimalStatusDeceased = 2
)

func GetStatusName(status int) string {
	switch status {
	case AnimalStatusActive:
		return "Ativo"
	case AnimalStatusSold:
		return "Vendido"
	case AnimalStatusDeceased:
		return "Falecido"
	default:
		return "Desconhecido"
	}
}
