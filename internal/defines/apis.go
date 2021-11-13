package defines

const (
	// Transactions
	APITransactionAddURL     = "/transaction"
	APITransactionsGetAllURL = "/transactions"
	APITransactionsUpdateURL = "/transactions/{" + ParamMsgID + "}"

	// Users
	APIUsersCreateURL = "/user"

	// Wallets
	APIWalletsCreateURL = "/wallet"
	APIWalletsGetAllURL = "/wallets"
)
