package defines

const (
	// Transactions
	APITransactionAddURL     = "/transactions"
	APITransactionsGetAllURL = "/transactions"
	APITransactionsUpdateURL = "/transactions/{" + ParamMsgID + "}"

	// Users
	APIUsersCreateURL = "/users"

	// Wallets
	APIWalletsCreateURL = "/wallets"
	APIWalletsGetAllURL = "/wallets"
	APIWalletsGetByID   = "/wallet/{" + ParamWalletID + "}"

	// Categories
	APICategoriesCreateURL = "/categories"
	APICategoriesGetAllURL = "/categories"
	APICategoriesGetByID   = "/category/{" + ParamCategoryID + "}"
)
