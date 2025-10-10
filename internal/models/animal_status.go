package models

const (
	AnimalStatusActive   = 0 // Animal ativo (não vendido)
	AnimalStatusSold     = 1 // Animal vendido
	AnimalStatusDeceased = 2 // Animal falecido
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
