package models

import (
	"time"
)

type ReproductionPhase int

const (
	PhaseLactacao ReproductionPhase = iota
	PhaseSecando
	PhaseVazias
	PhasePrenhas
)

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

type Reproduction struct {
	ID                     uint              `gorm:"primaryKey"`
	AnimalID               uint              `gorm:"not null"`
	Animal                 Animal            `gorm:"foreignKey:AnimalID"`
	CurrentPhase           ReproductionPhase `gorm:"not null;default:2"`
	InseminationDate       *time.Time
	InseminationType       string
	PregnancyDate          *time.Time
	ExpectedBirthDate      *time.Time
	ActualBirthDate        *time.Time
	LactationStartDate     *time.Time
	LactationEndDate       *time.Time
	DryPeriodStartDate     *time.Time
	VeterinaryConfirmation bool `gorm:"default:false"`
	Observations           string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
