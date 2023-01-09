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
		CommandCreateWallet + " - Crea una billetera con un nombre y un monto inicial. Ej: `/crearbilletera Débito $0.00`\n" +
		CommandCreateCategory + " - Crea una categoría\n" +
		CommandGetWallets + " - Muestra tus billeteras\n" +
		CommandExport + " - Exporta todas las transacciones en formato CSV"
	MessageCancel = "❌ Operación cancelada."

	// Errores
	MessageError               = "😨 Ups! Parece que hubo un error. "
	MessageErrorResponse       = MessageError + "\n```\n%+v\n```"
	MessageErrorWalletNotFound = "No encontré una billetera con ese nombre."

	// AddTransaction
	MessageAddPaymentResponse = "✅ Listo, ya anoté tu pago de *%s* (%s) por *$%.2f* con *%s*."
	MessageAddMoneyResponse   = "Listo, ya anoté ingreso de dinero de *%s* (%s) por *$%.2f* en *%s*."

	// UpdateTransaction
	MesssageUpdatePaymentResponse = "Listo, ya modifiqué la transacción."

	// CreateWallet
	MessageCreateWalletWaitingName   = "¿Como se va a llamar la billetera?"
	MessageCreateWalletWaitingAmount = "¿Cuál es el monto inicial de la billetera?"
	MessageCreateWalletSuccess       = "¡Tu billetera *%s* está lista!"

	// CreateCategory
	MessageCreateCategoryWaitingName  = "¿Como se va a llamar la categoría?"
	MessageCreateCategoryWaitingEmoji = "¿Con qué emoji querés representar la categoría?"
	MessageCreateCategorySuccess      = "¡La categoría %s (%s) está lista!"
)
