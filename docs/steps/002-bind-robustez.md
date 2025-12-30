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
