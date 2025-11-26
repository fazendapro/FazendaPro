package handlers

const (
	ErrMethodNotAllowed       = "Método não permitido"
	ErrAnimalIDRequired       = "ID do animal é obrigatório"
	ErrInvalidAnimalID        = "ID do animal inválido"
	ErrDecodeJSON             = "Erro ao decodificar JSON: "
	ErrInternalServer         = "Erro interno do servidor"
	ErrGenerateToken          = "Erro ao gerar token"
	ErrFarmIDNotFound         = "Farm ID not found in context"
	ErrInvalidFarmID          = "ID da fazenda inválido"
	ErrInvalidSaleID          = "Invalid sale ID"
	ErrSaleNotFound           = "Sale not found"
	ErrSaleNotBelongsToFarm   = "Sale does not belong to the specified farm"
	ErrAnimalNotBelongsToFarm = "Animal does not belong to the specified farm"
	ErrInvalidMonthsParam     = "Invalid months parameter"
)

const (
	HeaderContentType              = "Content-Type"
	ContentTypeJSON                = "application/json"
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
)

const (
	DateFormatISO      = "2006-01-02"
	DateFormatDateTime = "2006-01-02 15:04:05"
	DateFormatISO8601  = "2006-01-02T15:04:05.000Z"
)
