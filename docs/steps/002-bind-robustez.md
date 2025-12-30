# 002 - Bind a interfaz y robustez

## Objetivo
Capturar tráfico solo de una interfaz concreta y cerrar de forma limpia, con estadísticas fiables.

## Qué añade este paso
- Bind del socket AF_PACKET a una interfaz por nombre.
- Parada con Ctrl+C sin `os.Exit` (cierre limpio).
- El bucle sale cerrando el socket para desbloquear `Recvfrom`.
- Estadísticas por EtherType: IPv4, IPv6, ARP, VLAN y Otros.
- Robustez: detecta y cuenta tramas cortas/corruptas sin contaminar stats.

## Limitaciones (a propósito)
- Sin modo promiscuo explícito.
- Sin parseo de VLAN (solo conteo por tag).
- Solo capa 2: todavía no se parsea IPv4/TCP/UDP.

## Ejecución
```bash
sudo go run .
```

```bash
--- Sniffer v0.2: Robustez y Calidad ---
Escuchando en: enp6s0
Capturando... (Ctrl+C para parar)
Capturados: 10400 | Bytes: 7623174^C

Solicitud de parada recibida...

============= RESUMEN =============
Duración:        35.39790627s
Total Paquetes:  10442
Total Bytes:     7631390
-----------------------------------
IPv4:            5191
IPv6:            6
ARP:             24
VLAN (Tag):      0
Otros:           0
Short/Corrupt:   0
===================================

--- Limpiando recursos (Cerrando Socket) ---
```
