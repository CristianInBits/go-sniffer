# 001 - Ethernet (Capa 2)

## Objetivo
Ver tramas reales y extraer MAC src/dst + EtherType.

## Qué ya funciona
- Captura con AF_PACKET (raw socket)
- Parseo Ethernet (14 bytes)
- Log por paquete: MACs, EtherType, tamaño

## Limitaciones (a propósito)
- Sin bind a interfaz
- Sin promiscuo explícito
- Sin soporte VLAN (0x8100/0x88A8)
- Buffer fijo (posible truncado)

## Próximo
Bind a interfaz + contadores por EtherType + preparar salto a IPv4.
