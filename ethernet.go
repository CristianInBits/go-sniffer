package main

import (
	"encoding/binary" // Necesario para el "Endianness" (Big Endian)
	"fmt"
)

// Definimos el "Molde" de la cabecera Ethernet (14 bytes)
type EthernetHeader struct {
	DstMac    [6]byte // 6 bytes: Destino
	SrcMac    [6]byte // 6 bytes: Origen
	EtherType uint16  // 2 bytes: Tipo de protocolo (IPv4, ARP, IPv6...)
}

// ParseEthernet toma los bytes crudos y nos devuelve la estructura rellena
// También devuelve el 'payload' (los bytes que sobran y que contienen la IP)
func ParseEthernet(data []byte) (EthernetHeader, []byte) {
	
	// 1. Validación de seguridad: ¿Tenemos al menos 14 bytes?
	if len(data) < 14 {
		return EthernetHeader{}, nil // Si es muy corto, devolvemos vacío
	}

	// Creamos una variable del tipo de nuestra estructura
	var h EthernetHeader

	// 2. Extracción (Copia directa de bytes)
	// copy(destino, origen)
	copy(h.DstMac[:], data[0:6]) // Copiamos bytes 0-5 a DstMac
	copy(h.SrcMac[:], data[6:12]) // Copiamos bytes 6-11 a SrcMac

	// 3. Conversión de Endianness
	// Leemos bytes 12 y 13, y los giramos de Big Endian a Little Endian
	h.EtherType = binary.BigEndian.Uint16(data[12:14])

	// 4. Devolvemos la cabecera y el resto de los datos (del 14 al final)
	return h, data[14:]
}

// Helper visual: Convierte los bytes de la MAC en texto legible (ej: "00:1a:2b...")
func macToString(b [6]byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", b[0], b[1], b[2], b[3], b[4], b[5])
}
