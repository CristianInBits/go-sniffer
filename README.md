# Sniffer ligero (Go + Linux)

Proyecto educativo y largo: construir un sniffer paso a paso para entender cómo Linux gestiona paquetes.
Primera meta: captura en crudo + parseo de cabeceras (Ethernet → IPv4 → TCP/UDP).

## Estado
- [x] Paso 1: Captura con AF_PACKET + lectura de tramas
- [x] Paso 2: Parseo básico Ethernet (MAC src/dst + EtherType)
- [ ] Paso 3: IPv4
- [ ] Paso 4: TCP/UDP
- [ ] Paso 5: filtros y stats

## Docs
Ver `docs/steps/`.
