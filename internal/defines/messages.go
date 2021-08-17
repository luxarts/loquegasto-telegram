package defines

const (
	MessageStart = "Hola %s!\n" +
		"🤓 Soy tu asistente de gastos. Acá vas a poder anotar todas las compras que hagas de una manera rápida para que puedas tener control sobre cómo usas tu dinero.\n" +
		"❓Para conocer todos los comandos disponibles escribí /ayuda."
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
	MessageError         = "Ups! Parece que hubo un error. 😨"
	MessageTotalResponse = "Llevás gastado *$%s*."

	// AddPayment
	MessagePaymentResponse           = "✅Listo, ya anoté tu pago de *%s* por *$%s*."
	MessagePaymentResponseWithSource = "✅Listo, ya anoté tu pago de *%s* por *$%s* con *%s*."
)
