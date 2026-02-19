package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	fmt.Println("NATS Publisher - Enviando mensajes...")
	
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

	// Lista de mensajes para enviar
	messages := []struct {
		subject string
		message string
	}{
		{"greet.joe", "Hola Joe!"},
		{"greet.pam", "Hola Pam!"},
		{"greet.bob", "Hola Bob!"},
		{"greet.alice", "Hola Alice!"},
	}

	fmt.Println("Enviando mensajes...")
	
	for i, msg := range messages {
		err := nc.Publish(msg.subject, []byte(msg.message))
		if err != nil {
			fmt.Printf("Error enviando mensaje: %v\n", err)
			continue
		}
		
		fmt.Printf("[%d] Enviado: %q al subject %q\n", i+1, msg.message, msg.subject)
		time.Sleep(1 * time.Second) // Pausa entre mensajes
	}

	// Forzar el envío de todos los mensajes
	err = nc.Flush()
	if err != nil {
		fmt.Printf("Error forzando envío: %v\n", err)
	}

	fmt.Println("Todos los mensajes enviados. Publisher terminando...")
}