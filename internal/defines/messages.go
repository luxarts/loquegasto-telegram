package defines

const (
	// Chat individual
	MessageStart = "Hola %s!\n" +
		"🤓 Soy tu asistente de gastos. Acá vas a poder anotar todas las transacciones que hagas de una manera rápida para que puedas tener control sobre cómo usás tu dinero.\n" +
		"Simplemente tenes que escribir el monto de la transacción seguido de una descripción.\n" +
		"Para que puedas tener controlados tus gastos es necesario que me indiques desde cual _billetera_ estás haciendo el movimiento. " +
		"Las _billeteras_ son únicamente para que sepas con qué medio hiciste la transacción, podés crear la cantidad que quieras usando el comando " + CommandCreateWallet + ".\n" +
		"Te creé la billetera *Efectivo* con un balance inicial de *$0.00* que voy a usar por defecto si no me indicás otra.\n\n" +
		"❓Para conocer todos los comandos disponibles escribí " + CommandHelp + ".\n"
	MessageHelp = "Para registrar un pago podés decírmelo usando la sintaxis:\n" +
		"`valor descripción`\n" +
		"Ejemplo:\n" +
		"`$100.50 pizza`\n" +
		"También podés decirme el método que usaste para pagar.\n" +
		"Ejemplo:\n" +
		"`$200 taxi (crédito)`\n" +
		"Si querés registrar un ingreso de dinero podés decírmelo indicando el monto con el signo _-_.\n" +
		"Ejemplo:\n" +
		"`$-1000 Mamá (débito)`\n\n" +
		"Los comandos disponibles son:\n" +
		"/ayuda - _Muestra este mensaje_\n" +
		"/billeteras - _Muestra tus billeteras_"

	// Chat grupal
	MessageStartGroup = "Hola!\n" +
		"🤓 Soy el asistente de gastos y mi tarea es ayudarlos a dividir los gastos que me vayan diciendo.\n" +
		"Para comenzar necesito que cada uno me envie el comando " + CommandStart + " para poder registrarlos.\n\n" +
		"*Cómo usarme*\n" +
		"Para anotar un gasto simplemente escriban el comando " + CommandAddTransaction + " seguido del monto y de una descripción.\n" +
		"Ejemplo: `" + CommandAddTransaction + " $1234.56 Comida` o `" + CommandAddTransaction + " 1234 Bebidas`.\n\n" +
		"Para consultar cuánto dinero gastó cada uno pueden usar el comando " + CommandSplit + ".\n" +
		"Si quieren reiniciar la cuenta me lo pueden decir con el comando " + CommandReset + "\n\n" +
		"❓Para conocer todos los comandos disponibles escribí " + CommandHelp + ".\n"

	// Errores
	MessageError               = "Ups! Parece que hubo un error. 😨"
	MessageErrorResponse       = "Ups! Parece que hubo un error. 😨 \n```\n%s\n```"
	MessageErrorWalletNotFound = "No encontré una billetera con ese nombre."

	// AddTransaction
	MessageAddPaymentResponse           = "Listo, ya anoté tu pago de *%s* por *$%.2f*."
	MessageAddPaymentResponseWithWallet = "Listo, ya anoté tu pago de *%s* por *$%.2f* con *%s*."
	MessageAddMoneyResponse             = "Listo, ya anoté ingreso de dinero de *%s* por *$%.2f* en *%s*."

	// UpdateTransaction
	MesssageUpdatePaymentResponse = "Listo, ya modifiqué la transacción."

	// CreateWallet
	MessageCreateWallet = "Tu billetera *%s* está lista!"
)
