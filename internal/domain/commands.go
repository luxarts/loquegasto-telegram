package domain

import (
	"loquegasto-telegram/internal/defines"
	"strconv"
	"strings"
)

type CommandTransactionPayload struct {
	Amount      float64
	Description string
	WalletName  string
}

func (ctp *CommandTransactionPayload) Parse(message string) error {
	// Search for amount and description
	result := defines.RegexTransaction.FindAllStringSubmatch(message, -1)

	// Validate results
	if len(result) != 1 || len(result[0]) != 4 {
		return defines.ErrInvalidSyntax
	}

	// Amount capture group 1
	amountStr := result[0][1]

	// Parse decimal as dot for internal usage and colon for response
	amountStr = strings.Replace(amountStr, ",", ".", 1)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return err
	}

	// Description capture group 2
	description := result[0][2]

	// Wallet will be empty if capture group 3 isn't set
	walletName := result[0][3]

	ctp.WalletName = walletName
	ctp.Amount = amount
	ctp.Description = description

	return nil
}
