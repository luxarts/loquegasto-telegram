# loquegasto-telegram
LoQueGasto Telegram Bot backend

### Dividir gastos
1. El bot se agrega a un grupo
2. Las personas del grupo deben enviar el comando `/start` para que el bot las registre.  
 _Nota: Esto no es necesario para las personas que se agreguen al grupo después de que se haya agregado el bot._  
 _Opcional: El usuario configura la cantidad de personas por las que aportará usando el comando `/aportarpor <cantidad de personas>`_
3. Las personas envían los gastos usando el comando `/anotar <monto> <descripción>`
4. Se le pide al bot la divisón de gastos con el comando `/dividir`
5. Se reinician los gastos del grupo con `/reiniciar`