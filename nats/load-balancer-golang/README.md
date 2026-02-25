# NATS Load Balancer con Queue Groups

Este proyecto demuestra cómo implementar **load balancing** entre múltiples subscribers usando **NATS Core** con Queue Groups en Go.

## 📋 Requisitos

- Go 1.20+
- Servidor NATS corriendo (local o Docker)

## 🚀 Cómo ejecutar

### 1. Iniciar el servidor NATS

```bash
# Opción A: Docker
docker run -p 4222:4222 nats:latest

# Opción B: Binario local
nats-server
```

### 2. Iniciar múltiples Subscribers (en terminales separadas)

#### Opción A: Probar load balancing dentro de un grupo

Abre **3 o más terminales** y ejecuta el mismo subscriber:

```bash
# Terminal 1
go run ./subscriber/main.go

# Terminal 2
go run ./subscriber/main.go

# Terminal 3
go run ./subscriber/main.go
```

Los mensajes se **distribuyen** entre los 3 (cada uno recibe ~33%).

#### Opción B: Probar múltiples queue groups

Ejecuta subscribers de **diferentes grupos**:

```bash
# Terminal 1 - Queue group "workers"
go run ./subscriber/main.go

# Terminal 2 - Queue group "workers"
go run ./subscriber/main.go

# Terminal 3 - Queue group "queue" (segundo grupo)
go run ./second-subscriber/main.go
```

En este caso:
- Los mensajes se **distribuyen** entre Terminal 1 y 2 (grupo `workers`)
- Terminal 3 recibe **TODOS** los mensajes (grupo `queue`)

> ⚠️ **Importante**: Debes tener **varios subscribers corriendo antes** de enviar mensajes para ver el load balancing en acción.

### 3. Ejecutar el Publisher

En otra terminal:

```bash
go run ./publisher/main.go
```

## 📊 Resultado esperado

### Con 3 subscribers del mismo grupo (workers)

```
Subscriber 1: Mensaje #1 - "1) Hola"
Subscriber 2: Mensaje #1 - "2) Hola lucas!"
Subscriber 3: Mensaje #1 - "3) Hola kimy!"
Subscriber 1: Mensaje #2 - "4) Hola carlos!"
...
```

Los mensajes se **distribuyen automáticamente** entre los subscribers del grupo.

### Con 2 grupos diferentes (workers + queue)

```
# Grupo "workers" (2 instancias) - Se reparten los mensajes:
Subscriber 1 (workers): "1) Hola", "3) Hola kimy!", "5) Hola andres!"
Subscriber 2 (workers): "2) Hola lucas!", "4) Hola carlos!", "6) Hola lele!"

# Grupo "queue" (1 instancia) - Recibe TODOS los mensajes:
Subscriber 3 (queue): "1) Hola", "2) Hola lucas!", "3) Hola kimy!", ...todos
```

## ✅ Ventajas del Load Balancer con Queue Groups

| Ventaja | Descripción |
|---------|-------------|
| **Distribución automática** | NATS distribuye mensajes entre subscribers del mismo queue group |
| **Escalabilidad horizontal** | Agregar más subscribers aumenta la capacidad de procesamiento |
| **Sin configuración adicional** | Solo necesitas usar `QueueSubscribe()` con el mismo nombre de grupo |
| **Alta disponibilidad** | Si un subscriber falla, los demás continúan recibiendo mensajes |
| **Simplicidad** | No requiere brokers externos ni configuración compleja |
| **Baja latencia** | At-most-once delivery, ideal para sistemas en tiempo real |

## ❌ Limitaciones (vs NATS JetStream)

| Limitación | Descripción |
|------------|-------------|
| **Sin persistencia** | Los mensajes NO se guardan en disco. Si no hay subscribers activos, **los mensajes se pierden** |
| **Sin replay** | No puedes reproducir mensajes históricos |
| **Sin acknowledgment** | No hay confirmación de que el mensaje fue procesado correctamente |
| **Sin reintentos** | Si el subscriber falla al procesar, el mensaje se pierde |
| **Sin exactly-once** | Solo garantiza at-most-once delivery |
| **Sin consumer durable** | Los subscribers no recuerdan su posición si se desconectan |

## 🔄 ¿Cuándo usar cada uno?

### Usa NATS Core (este ejemplo) cuando:
- Los mensajes pueden perderse sin consecuencias graves
- Necesitas máxima velocidad y baja latencia
- Los subscribers siempre están activos
- Procesamiento fire-and-forget

### Usa NATS JetStream cuando:
- Los mensajes **no pueden perderse**
- Necesitas garantías de entrega (at-least-once, exactly-once)
- Requieres persistencia y replay de mensajes
- Los subscribers pueden estar offline temporalmente

## 📁 Estructura del proyecto

```
.
├── go.mod
├── README.md
├── publisher/
│   └── main.go           # Envía 7 mensajes a subjects "saludo.*"
├── subscriber/
│   └── main.go           # Queue group: "workers"
└── second-subscriber/
    └── main.go           # Queue group: "queue"
```

> 💡 **Nota**: Las dos carpetas de subscribers usan **diferentes queue groups**:
> - `subscriber/` → queue group `"workers"`
> - `second-subscriber/` → queue group `"queue"`
> 
> Esto significa que **ambos grupos recibirán todos los mensajes** (uno por grupo), demostrando el comportamiento de múltiples queue groups.

## 🔑 Código clave

El load balancing se logra con una sola línea usando `QueueSubscribe`:

```go
// Sin load balancing - todos reciben todos los mensajes
nc.Subscribe("saludo.*", handler)

// CON load balancing - mensajes distribuidos entre el grupo
nc.QueueSubscribe("saludo.*", "workers", handler)
```

En este proyecto:
- `subscriber/main.go` usa: `queue := "workers"`
- `second-subscriber/main.go` usa: `queue := "queue"`

Todos los subscribers que usen el mismo queue group compartirán la carga de trabajo.

## 🎯 Comportamiento según el nombre del Queue Group

El **nombre del queue group** determina cómo se distribuyen los mensajes:

### Mismo Queue Group = Comparten carga

```go
// Subscriber A
nc.QueueSubscribe("saludo.*", "workers", handlerA)

// Subscriber B  
nc.QueueSubscribe("saludo.*", "workers", handlerB)

// Subscriber C
nc.QueueSubscribe("saludo.*", "workers", handlerC)
```

**Resultado**: Cada mensaje llega a **solo uno** de los 3 subscribers (A, B o C).

### Diferentes Queue Groups = Cada grupo recibe todo

```go
// Grupo "workers" - procesamiento general
nc.QueueSubscribe("saludo.*", "workers", handlerWorker)

// Grupo "loggers" - logging
nc.QueueSubscribe("saludo.*", "loggers", handlerLogger)

// Grupo "analytics" - métricas
nc.QueueSubscribe("saludo.*", "analytics", handlerAnalytics)
```

**Resultado**: Cada mensaje se envía a **un subscriber de cada grupo**:
- 1 subscriber del grupo `workers` recibe el mensaje
- 1 subscriber del grupo `loggers` recibe el mensaje  
- 1 subscriber del grupo `analytics` recibe el mensaje

### Diagrama de distribución

```
                    ┌─────────────────────────────────────────────────┐
                    │              Mensaje publicado                  │
                    └─────────────────────────────────────────────────┘
                                          │
                    ┌─────────────────────┼─────────────────────┐
                    ▼                     ▼                     ▼
            ┌───────────────┐     ┌───────────────┐     ┌───────────────┐
            │ Queue:workers │     │ Queue:loggers │     │ Queue:analytics│
            └───────────────┘     └───────────────┘     └───────────────┘
                    │                     │                     │
            ┌───────┴───────┐             │                     │
            ▼               ▼             ▼                     ▼
       Worker A         Worker B      Logger A            Analytics A
       (recibe)        (no recibe)   (recibe)             (recibe)
```

### Caso de uso típico

```go
// Grupo de procesamiento principal (3 instancias)
nc.QueueSubscribe("orders.*", "processors", processOrder)

// Grupo de auditoría (2 instancias) 
nc.QueueSubscribe("orders.*", "auditors", auditOrder)

// Grupo de notificaciones (1 instancia)
nc.QueueSubscribe("orders.*", "notifiers", notifyOrder)
```

Cada orden será:
1. **Procesada** por 1 de los 3 processors
2. **Auditada** por 1 de los 2 auditors
3. **Notificada** por el notifier
