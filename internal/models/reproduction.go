package models

import (
	"time"
)

// ReproductionPhase representa as fases de reprodução de uma vaca
type ReproductionPhase int

const (
	PhaseLactacao ReproductionPhase = iota // 0 - Lactação
	PhaseSecando                           // 1 - Secando
	PhaseVazias                            // 2 - Vazias
	PhasePrenhas                           // 3 - Prenhas
)

// String retorna a representação em string da fase
func (rp ReproductionPhase) String() string {
	switch rp {
	case PhaseLactacao:
		return "Lactação"
	case PhaseSecando:
		return "Secando"
	case PhaseVazias:
		return "Vazias"
	case PhasePrenhas:
		return "Prenhas"
	default:
		return "Desconhecida"
	}
}

// Reproduction representa o registro de reprodução de um animal
type Reproduction struct {
	ID                     uint              `gorm:"primaryKey"`
	AnimalID               uint              `gorm:"not null"`
	Animal                 Animal            `gorm:"foreignKey:AnimalID"`
	CurrentPhase           ReproductionPhase `gorm:"not null;default:2"` // 2 = Vazias por padrão
	InseminationDate       *time.Time        // Data da inseminação
	InseminationType       string            // Natural, Artificial Insemination, etc.
	PregnancyDate          *time.Time        // Data de confirmação da prenhez
	ExpectedBirthDate      *time.Time        // Data prevista do parto
	ActualBirthDate        *time.Time        // Data real do parto
	LactationStartDate     *time.Time        // Data de início da lactação
	LactationEndDate       *time.Time        // Data de fim da lactação
	DryPeriodStartDate     *time.Time        // Data de início do período seco
	VeterinaryConfirmation bool              `gorm:"default:false"` // Confirmação veterinária
	Observations           string            // Observações gerais
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
