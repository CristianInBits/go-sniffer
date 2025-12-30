package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings" // Lo usamos para detectar el error de cierre
	"syscall"
	"time"
)

func main() {
	// --- CONFIGURACIÓN ---
	const interfaceName = "enp6s0" // ¡PON AQUÍ TU INTERFAZ!

	fmt.Printf("--- Sniffer v0.2: Robustez y Calidad ---\n")
	fmt.Printf("Escuchando en: %s\n", interfaceName)

	// 1. Obtener interfaz
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		panic(fmt.Sprintf("No encuentro la interfaz %s: %v", interfaceName, err))
	}

	// 2. Crear Socket Raw
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0x0300)
	if err != nil {
		panic(fmt.Sprintf("Error al crear socket: %v", err))
	}

	// --- CAMBIO 1: El defer se ejecutará porque ya no usamos os.Exit ---
	defer func() {
		fmt.Println("\n--- Limpiando recursos (Cerrando Socket) ---")
		syscall.Close(fd)
	}()

	// 3. Bind a la interfaz
	sll := syscall.SockaddrLinklayer{
		Protocol: 0x0300,
		Ifindex:  iface.Index,
	}
	if err := syscall.Bind(fd, &sll); err != nil {
		panic(fmt.Sprintf("Error bind: %v", err))
	}

	// Variables para estadísticas
	var (
		packetsTotal uint64
		bytesTotal   uint64
		countIPv4    uint64
		countIPv6    uint64
		countARP     uint64
		countVLAN    uint64 // ¡Nuevo!
		countOther   uint64
		countShort   uint64 // ¡Nuevo! Robustez
	)

	// Canal para la señal de parada
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Goroutine de control
	go func() {
		<-sigChan // Espera Ctrl+C
		fmt.Println("\n\nSolicitud de parada recibida...")
		// TRUCO: Cerramos el socket AQUÍ para desbloquear el Recvfrom del main
		syscall.Close(fd)
	}()

	fmt.Println("Capturando... (Ctrl+C para parar)")
	startTime := time.Now()
	buffer := make([]byte, 2048)

	// --- BUCLE PRINCIPAL ---
	for {
		n, _, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			// Si el error dice "bad file descriptor" es que nosotros cerramos el socket
			// para salir. Es el momento de romper el bucle limpiamente.
			if strings.Contains(err.Error(), "bad file descriptor") || strings.Contains(err.Error(), "closed") {
				break
			}
			// Si es otro error real, lo imprimimos y seguimos
			fmt.Printf("Error leyendo: %v\n", err)
			continue
		}

		// --- PROCESAMIENTO ---
		packetsTotal++
		bytesTotal += uint64(n)

		// --- PROCESAMIENTO ---
		packetsTotal++
		bytesTotal += uint64(n)

		// Feedback visual cada 100 paquetes para saber que está vivo
		if packetsTotal%100 == 0 {
			// \r hace "retorno de carro": vuelve al principio de la línea y sobrescribe
			fmt.Printf("\rCapturados: %d | Bytes: %d", packetsTotal, bytesTotal)
		}

		rawPacket := buffer[:n]

		// Usamos el parser
		ethHeader, payload := ParseEthernet(rawPacket)

		// --- CAMBIO 2: Robustez ante frames cortos ---
		if payload == nil {
			countShort++
			continue // Saltamos al siguiente ciclo sin contaminar estadísticas
		}

		// --- CAMBIO 3: Clasificación mejorada (VLANs) ---
		switch ethHeader.EtherType {
		case 0x0800:
			countIPv4++
			// fmt.Print(".") // Comentado para limpiar la salida
		case 0x86DD:
			countIPv6++
		case 0x0806:
			countARP++
		case 0x8100, 0x88A8: // Standard VLAN y Service VLAN
			countVLAN++
		default:
			countOther++
		}
	}

	// --- INFORME FINAL (Se ejecuta al salir del bucle) ---
	duration := time.Since(startTime)
	fmt.Println("\n============= RESUMEN =============")
	fmt.Printf("Duración:        %s\n", duration)
	fmt.Printf("Total Paquetes:  %d\n", packetsTotal)
	fmt.Printf("Total Bytes:     %d\n", bytesTotal)
	fmt.Println("-----------------------------------")
	fmt.Printf("IPv4:            %d\n", countIPv4)
	fmt.Printf("IPv6:            %d\n", countIPv6)
	fmt.Printf("ARP:             %d\n", countARP)
	fmt.Printf("VLAN (Tag):      %d\n", countVLAN)
	fmt.Printf("Otros:           %d\n", countOther)
	fmt.Printf("Short/Corrupt:   %d\n", countShort)
	fmt.Println("===================================")
}
