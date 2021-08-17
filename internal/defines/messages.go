package defines

const (
	MessageStart = "Hola %s!\n📊Soy tu asistente de gastos. Acá vas a poder anotar todas las compras que hagas de una manera rápida para que puedas tener control sobre cómo usas tu dinero.\n🆘Para conocer todos los comandos disponibles escribí /ayuda."
	MessageHelp = "Para agregar un pago podés decírmelo usando la sintaxis:\n`valor descripción`\nEjemplo:\n`100 comida`\n\nLos comandos disponibles son:\n/ayuda\n/ping"
	MessageError = "Ups! Hubo un error."

	// AddPayment
	MessagePaymentResponse = "Listo, ya anoté tu pago de _%s_ por $%d."
	MessagePaymentResponseWithSource = "Listo, ya anoté tu pago de _%s_ por $%d con %s."
)

