package defines

const (
	// Chat individual
	MessageStart = "Hola %s!\n" +
		"🤓 Soy tu asistente de gastos. Acá vas a poder anotar todas las transacciones que hagas de una manera rápida para que puedas tener control sobre cómo usás tu dinero.\n" +
		"❓ Usá el comando " + CommandHelp + " para ver todo lo que podés hacer."
	MessageHelp = "Para registrar un pago podés decírmelo usando la sintaxis:\n" +
		"`valor descripción`\n" +
		"Ejemplo:\n" +
		"`$100.50 pizza`\n\n" +
		"Si querés registrar un ingreso de dinero podés decírmelo usando el signo `-` antes del valor.\n" +
		"Ejemplo:\n" +
		"`$-1000 Mamá`\n\n" +
		"_Nota: Podés indicarme el valor sin usar el signo $._\n\n" +
		"Los comandos disponibles son:\n" +
		CommandHelp + " - Muestra este mensaje\n" +
		CommandGetWallets + " - Muestra tus billeteras"

	// Errores
	MessageError               = "Ups! Parece que hubo un error. 😨"
	MessageErrorResponse       = "Ups! Parece que hubo un error. 😨 \n```\n%s\n```"
	MessageErrorWalletNotFound = "No encontré una billetera con ese nombre."

	// AddTransaction
	MessageAddPaymentResponseWithWalletAndCategory = "Listo, ya anoté tu pago de *%s* (%s) por *$%.2f* con *%s*."
	MessageAddMoneyResponse                        = "Listo, ya anoté ingreso de dinero de *%s* por *$%.2f* en *%s*."

	// UpdateTransaction
	MesssageUpdatePaymentResponse = "Listo, ya modifiqué la transacción."

	// CreateWallet
	MessageCreateWallet = "Tu billetera *%s* está lista!"
)
