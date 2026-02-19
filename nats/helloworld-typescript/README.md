# NATS Hello World - TypeScript

Este proyecto es un ejemplo de aprendizaje para entender el patrón **Publish/Subscribe** usando **NATS** con **TypeScript** y **Node.js**.

## Descripción

NATS es un sistema de mensajería de alto rendimiento que permite la comunicación asíncrona entre servicios usando el patrón publish/subscribe. Este proyecto incluye dos componentes principales:

- **Publisher** (`src/publisher.ts`): Publica mensajes a un tema específico
- **Subscriber** (`src/subscriber.ts`): Se suscribe a un tema y recibe mensajes

## Requisitos Previos

- Node.js (versión 16 o superior)
- npm
- Un servidor NATS ejecutándose (por defecto en `localhost:4222`)

## Instalación

1. Clona o descarga este proyecto
2. Navega al directorio del proyecto:
   ```bash
   cd nats/helloworld
   ```

3. Instala las dependencias necesarias:
   ```bash
   npm install nats@latest
   npm install -D typescript @types/node ts-node
   ```

## Configuración de Scripts

Agrega los siguientes scripts a tu `package.json` para facilitar la ejecución:

```json
{
  "scripts": {
    "publisher": "ts-node src/publisher.ts",
    "subscriber": "ts-node src/subscriber.ts",
    "build": "tsc",
    "start:pub": "node dist/publisher.js",
    "start:sub": "node dist/subscriber.js"
  }
}
```

## Servidor NATS

Antes de ejecutar los ejemplos, necesitas tener un servidor NATS ejecutándose:

### Opción 1: Docker (Recomendado)
```bash
docker run -p 4222:4222 -ti nats:latest
```

### Opción 2: Instalación local
Descarga NATS desde https://nats.io/download/ y ejecuta:
```bash
nats-server
```

## Cómo Ejecutar

### Usando ts-node (Desarrollo)

1. **Ejecutar el Subscriber** (en una terminal):
   ```bash
   npm run subscriber
   # o directamente:
   npx ts-node src/subscriber.ts
   ```

2. **Ejecutar el Publisher** (en otra terminal):
   ```bash
   npm run publisher
   # o directamente:
   npx ts-node src/publisher.ts
   ```

### Usando TypeScript compilado (Producción)

1. **Compilar el proyecto**:
   ```bash
   npm run build
   ```

2. **Ejecutar el Subscriber** (en una terminal):
   ```bash
   npm run start:sub
   ```

3. **Ejecutar el Publisher** (en otra terminal):
   ```bash
   npm run start:pub
   ```

## Documentación Oficial

Para más información sobre NATS y el cliente de Node.js:

- **Documentación oficial de NATS Node.js**: https://github.com/nats-io/nats.node
- **Documentación general de NATS**: https://docs.nats.io/
- **Ejemplos adicionales**: https://github.com/nats-io/nats.node/tree/main/examples
