package tests

// HTTP Headers
const (
	HeaderContentType     = "Content-Type"
	ContentTypeJSON       = "application/json"
	HeaderAuthorization   = "Authorization"
	BearerPrefix            = "Bearer "
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
)

// API Endpoints
const (
	EndpointAnimals              = "/animals"
	EndpointAnimalsPhoto         = "/animals/photo"
	EndpointAnimalsWithID        = "/animals?id=1"
	EndpointAPIAuthLogin         = "/api/auth/login"
	EndpointAPIAuthRegister      = "/api/auth/register"
	EndpointAPIAuthRefresh       = "/api/auth/refresh"
	EndpointAPIAuthLogout        = "/api/auth/logout"
	EndpointDebts                = "/debts"
	EndpointAPIv1FarmsUser      = "/api/v1/farms/user"
	EndpointAPIv1FarmsSelect    = "/api/v1/farms/select"
	EndpointAPIv1Farms          = "/api/v1/farms"
	EndpointFarmWithID           = "/farm?id=1"
	EndpointMilkCollections      = "/milk-collections"
	EndpointMilkCollectionsFarm  = "/milk-collections/farm/1"
	EndpointReproductions        = "/reproductions"
	EndpointReproductionsPhase   = "/reproductions/phase"
	EndpointReproductionsWithID  = "/reproductions?id=1"
	EndpointAPIv1ReproductionsNextToCalve = "/api/v1/reproductions/next-to-calve?farmId=1"
	EndpointSales                = "/sales"
	EndpointSalesWithID          = "/sales/{id}"
	EndpointSalesMonthlyStats    = "/sales/monthly-stats"
	EndpointSalesOverview        = "/sales/overview"
	EndpointSalesID              = "/sales/1"
	EndpointUsers                = "/users"
	EndpointUsersPerson          = "/users/person"
	EndpointUsersPersonWithID    = "/users/person?id=1"
	EndpointHealth               = "/health"
	EndpointAPIv1AuthRegister    = "/api/v1/auth/register"
	EndpointAPIv1Weights         = "/api/v1/weights"
	EndpointAPIv1WeightsAnimal  = "/api/v1/weights/animal/1"
	EndpointAPIv1WeightsFarm    = "/api/v1/weights/farm/1"
)

// Test Data - Emails
const (
	TestEmailExample     = "test@example.com"
	TestEmailJoao        = "joao@example.com"
	TestEmailJoaoFazenda = "joao@fazenda.com"
	TestEmailInvalid     = "invalid-email"
	TestEmailJohn        = "john@test.com"
)

// Test Data - Names
const (
	TestNameBoiTeste      = "Boi Teste"
	TestNameAnimalAtualizado = "Animal Atualizado"
	TestNameJoaoSilva     = "João Silva"
	TestNameTataSalt      = "Tata Salt"
	TestNameTestAnimal    = "Test Animal"
	TestNameTestSale      = "Test sale"
	TestNameVacaTeste     = "Vaca Teste"
	TestNameAnimalTeste   = "Animal Teste"
	TestNameBoiJoao      = "Boi João"
	TestNameMariaSantos   = "Maria Santos"
	TestNameCompradorTeste = "Comprador Teste"
)

// Test Data - Files and Paths
const (
	TestFileTestJPG                    = "test.jpg"
	TestFileFakeImageData              = "fake image data"
	TestFileNewLogoPNG                 = "new-logo.png"
	TestFileLogo1PNG                   = "logo1.png"
	TestPathTataPNG                    = "src/assets/images/mocked/cows/tata.png"
	TestPathLaysPNG                    = "src/assets/images/mocked/cows/lays.png"
)

// Test Data - Tokens and Secrets
const (
	TestSecret              = "test-secret"
	TestRefreshTokenValid   = "valid-refresh-token"
	TestRefreshToken        = "refresh-token"
)

// Test Data - Error Messages
const (
	TestErrorInvalidJSON                    = "invalid json"
	TestErrorDatabaseError                  = "database error"
	TestErrorSaleNotFound                   = "sale not found"
	TestErrorReproductionNotFound           = "registro de reprodução não encontrado"
	TestErrorUpdateError                   = "update error"
	TestErrorAnimalNotFound                 = "animal not found"
	TestErrorErroAoBuscar                  = "erro ao buscar"
)

// Test Data - Dates
const (
	TestDate20240115 = "2024-01-15"
	TestDateFormat   = "2006-01-02"
)

// Test Data - Company and Farm
const (
	TestCompanyName = "Test Company"
)

// Test Data - CORS
const (
	TestOriginLocalhost = "http://localhost:3000"
)

// Model Types (for type assertions in tests)
const (
	TypeModelsAnimal        = "*models.Animal"
	TypeModelsUser          = "*models.User"
	TypeModelsPerson        = "*models.Person"
	TypeModelsUserFarm      = "*models.UserFarm"
	TypeModelsMilkCollection = "*models.MilkCollection"
	TypeModelsReproduction  = "*models.Reproduction"
	TypeModelsSale          = "*models.Sale"
	TypeModelsWeight        = "*models.Weight"
	TypeTimeTime            = "time.Time"
	TypeTimeTimePtr         = "*time.Time"
)

