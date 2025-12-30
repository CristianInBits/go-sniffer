package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fmt.Println("--- Sniffer v0.1: Analizando Ethernet ---")

	// 1. Crear el Socket (igual que antes)
	// 0x0300 es ETH_P_ALL (capturar todo)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0x0300)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	buffer := make([]byte, 2048)

	for {
		// 2. Leer del cable
		n, _, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			continue
		}

		// Solo nos interesan los datos válidos (de 0 a n)
		rawPacket := buffer[:n]

		// 3. ¡AQUI USAMOS NUESTRO PARSER!
		ethHeader, payload := ParseEthernet(rawPacket)

		// Si el payload es nil, es que el paquete era erróneo o muy corto
		if payload == nil {
			continue
		}

		// 4. Imprimir resultados
		// Usamos %x para imprimir el EtherType en Hexadecimal (ej: 0800)
		fmt.Printf("[Ethernet] MAC Origen: %s -> MAC Destino: %s | Tipo: %04x | Payload restante: %d bytes\n",
			macToString(ethHeader.SrcMac),
			macToString(ethHeader.DstMac),
			ethHeader.EtherType,
			len(payload))
	}
}
