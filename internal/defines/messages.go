package defines

const (
	MessageStart = "Hola %s!\n" +
		"🤓 Soy tu asistente de gastos. Acá vas a poder anotar todas las compras que hagas de una manera rápida para que puedas tener control sobre cómo usas tu dinero.\n\n" +
		"Para que puedas tener controlados tus gastos es necesario que me indiques desde cual _billetera_ estás haciendo el movimiento. " +
		"Las _billeteras_ son únicamente para que sepas con qué medio hiciste la transacción, podés crear la cantidad que quieras.\n" +
		"Te creé una billetera *Efectivo* con un balance inicial de $0.00 para que uses por defecto pero podés crear otra en cualquier momento usando el comando /billetera.\n" +
		"Por ejemplo: `/billetera Débito $1234,56`\n" +
		"\n" +
		"❓Para conocer todos los comandos disponibles escribí /ayuda.\n"
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
	MessageError               = "Ups! Parece que hubo un error. 😨"
	MessageErrorWalletNotFound = "No encontré una billetera con ese nombre."

	// AddTransaction
	MessageAddPaymentResponseWithWallet = "Listo, ya anoté tu pago de *%s* por *$%.2f* con *%s*."
	MessageAddMoneyResponse             = "Listo, ya anoté ingreso de dinero de *%s* por *$%.2f* en *%s*."
	// UpdateTransaction
	MesssageUpdatePaymentResponse = "Listo, ya modifiqué la transacción."
)
