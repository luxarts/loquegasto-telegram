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
	MessageHelp = "Para agregar un pago podés decírmelo usando la sintaxis:\n" +
		"`valor descripción`\n" +
		"Ejemplo:\n" +
		"`$100.50 comida`\n" +
		"También podés decirme el método que usaste para pagar.\n" +
		"Ejemplo:\n" +
		"`$200 combustible (efectivo)`\n\n" +
		"Los comandos disponibles son:\n" +
		"/ayuda - _Muestra este mensaje_\n" +
		"/ping - _Prueba tu conexión con el bot_\n" +
		"/total - _Muestra el total gastado_"
	MessageError            = "Ups! Parece que hubo un error. 😨"
	MessageConsumosResponse = "*%s:* $%.2f."

	// AddPayment
	MessagePaymentResponse           = "Listo, ya anoté tu pago de *%s* por *$%.2f*."
	MessagePaymentResponseWithSource = "Listo, ya anoté tu pago de *%s* por *$%.2f* con *%s*."
	// UpdatePayment
	MesssagePaymentUpdatedResponse = "Listo, ya modifiqué tu pago."
)
