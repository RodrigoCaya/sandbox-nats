package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	fmt.Println("NATS Subscriber - Esperando mensajes...")
	
	// Para conectarse a NATS con port forwarding desde host local en dev container
	url := "nats://host.docker.internal:4222"
	fmt.Printf("Conectando a NATS en: %s\n", url)
	
	nc, err := nats.Connect(url, nats.Timeout(5*time.Second))
	if err != nil {
		fmt.Printf("Error conectando a NATS: %v\n", err)
		os.Exit(1)
	}
	defer nc.Drain()
	
	fmt.Println("Conectado exitosamente a NATS")

	// Suscribirse al patrón greet.*
	subject := "greet.*"
	fmt.Printf("Suscribiéndose al subject: %s\n", subject)
	
	messageCount := 0
	
	// Crear suscripción asíncrona
	sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		messageCount++
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("[%s] Mensaje #%d recibido:\n", timestamp, messageCount)
		fmt.Printf("  Subject: %s\n", msg.Subject)
		fmt.Printf("  Data: %q\n", string(msg.Data))
		fmt.Println("  ---")
	})
	
	if err != nil {
		fmt.Printf("Error creando suscripción: %v\n", err)
		os.Exit(1)
	}
	defer sub.Unsubscribe()

	fmt.Println("Subscriber activo. Presiona Ctrl+C para salir...")
	fmt.Println("Esperando mensajes...")

	// Configurar manejo de señales para salida limpia
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Mantener el programa corriendo hasta que se reciba una señal
	<-c
	
	fmt.Printf("\nShutdown signal recibido. Total de mensajes recibidos: %d\n", messageCount)
	fmt.Println("Subscriber terminando...")
}