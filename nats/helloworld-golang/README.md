# NATS Publisher/Subscriber Example

Este proyecto contiene ejemplos separados de un publisher y subscriber usando NATS.

## Estructura del proyecto

```
├── publisher/main.go    # Programa que envía mensajes
├── subscriber/main.go   # Programa que recibe mensajes
└── cmd/main.go         # Ejemplo original combinado
```

## Cómo usar

### 1. Ejecutar el Subscriber (en una terminal)
```bash
go run ./subscriber/main.go
```

El subscriber se mantendrá corriendo y esperando mensajes. Mostrará cada mensaje recibido con timestamp.

### 2. Ejecutar el Publisher (en otra terminal)
```bash
go run ./publisher/main.go
```

El publisher enviará una serie de mensajes y terminará.

## Funcionalidad

### Publisher
- Conecta a NATS en `nats://host.docker.internal:4222`
- Envía mensajes a diferentes subjects: `greet.joe`, `greet.pam`, `greet.bob`, `greet.alice`
- Incluye pausa entre mensajes para observar mejor el flujo
- Termina después de enviar todos los mensajes

### Subscriber
- Conecta a NATS en `nats://host.docker.internal:4222` 
- Se suscribe al patrón `greet.*` (recibe todos los mensajes que empiecen con "greet.")
- Muestra cada mensaje recibido con timestamp
- Se mantiene corriendo hasta que se presiona Ctrl+C

## Ejemplo de ejecución

**Terminal 1 (Subscriber):**
```
NATS Subscriber - Esperando mensajes...
Conectando a NATS en: nats://host.docker.internal:4222
Conectado exitosamente a NATS
Suscribiéndose al subject: greet.*
Subscriber activo. Presiona Ctrl+C para salir...
Esperando mensajes...
```

**Terminal 2 (Publisher):**
```
NATS Publisher - Enviando mensajes...
Conectando a NATS en: nats://host.docker.internal:4222
Conectado exitosamente a NATS
Enviando mensajes...
[1] Enviado: "Hola Joe!" al subject "greet.joe"
[2] Enviado: "Hola Pam!" al subject "greet.pam"
[3] Enviado: "Hola Bob!" al subject "greet.bob"
[4] Enviado: "Hola Alice!" al subject "greet.alice"
Todos los mensajes enviados. Publisher terminando...
```

**Terminal 1 después del publisher:**
```
[15:04:05] Mensaje #1 recibido:
  Subject: greet.joe
  Data: "Hola Joe!"
  ---
[15:04:06] Mensaje #2 recibido:
  Subject: greet.pam
  Data: "Hola Pam!"
  ---
[15:04:07] Mensaje #3 recibido:
  Subject: greet.bob
  Data: "Hola Bob!"
  ---
[15:04:08] Mensaje #4 recibido:
  Subject: greet.alice
  Data: "Hola Alice!"
  ---
```