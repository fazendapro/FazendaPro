package handlers

const (
	ErrMethodNotAllowed = "Método não permitido"
	ErrAnimalIDRequired = "ID do animal é obrigatório"
	ErrDecodeJSON       = "Erro ao decodificar JSON: "
	ErrInternalServer   = "Erro interno do servidor"
	ErrGenerateToken    = "Erro ao gerar token"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

const (
	DateFormatISO      = "2006-01-02"
	DateFormatDateTime = "2006-01-02 15:04:05"
	DateFormatISO8601  = "2006-01-02T15:04:05.000Z"
)
