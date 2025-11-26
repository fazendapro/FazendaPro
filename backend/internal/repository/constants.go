package repository

const (
	SQLWhereID              = "id = ?"
	SQLWhereFarmID          = "farm_id = ?"
	SQLWhereAnimalID        = "animal_id = ?"
	SQLWhereUserID          = "user_id = ?"
	SQLWhereCreatedAtRange  = "created_at >= ? AND created_at < ?"
	SQLOrderSaleDateDESC    = "sale_date DESC"
	SQLWhereFarmIDAndSex    = "farm_id = ? AND sex = ?"
	SQLWhereUserIDAndFarmID = "user_id = ? AND farm_id = ?"
	SQLWhereFarmIDAndEarTag = "farm_id = ? AND ear_tag_number_local = ?"
	SQLWhereAnimalsFarmID   = "animals.farm_id = ?"
)

const (
	ErrFindingUser                    = "error finding user: %w"
	ErrFindingUserFarms               = "error finding user farms: %w"
	ErrFindingUserFarm                = "error finding user farm: %w"
	ErrCountingUserFarms              = "error counting user farms: %w"
	ErrCreatingPerson                 = "error creating person: %w"
	ErrCreatingUser                   = "error creating user: %w"
	ErrCreatingCompany                = "error creating company: %w"
	ErrCreatingFarm                   = "error creating farm: %w"
	ErrUpdatingPersonData             = "error updating person data: %w"
	ErrCountingDebts                  = "error counting debts: %w"
	ErrFindingDebts                   = "error finding debts: %w"
	ErrCalculatingTotal               = "error calculating total by person: %w"
	ErrCountingMales                  = "error counting males: %w"
	ErrCountingFemales                = "error counting females: %w"
	ErrCountingTotalSold              = "error counting total sold: %w"
	ErrCalculatingRevenue             = "error calculating total revenue: %w"
	ErrSaleNotFoundOrNotBelongsToFarm = "sale not found or does not belong to farm"
)
