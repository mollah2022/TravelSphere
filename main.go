package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"

	_ "TravelSphere/routers"
	"TravelSphere/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] .env file not found, using system environment variables")
	}

	if port := os.Getenv("APP_PORT"); port != "" {
		if n, err := parseInt(port); err == nil {
			web.BConfig.Listen.HTTPPort = n
		}
	}

	web.BConfig.Listen.Graceful = true

	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionProvider = "memory"
	web.BConfig.WebConfig.Session.SessionName = "travelsphere_session"
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 86400

	web.SetStaticPath("/static", "static")

	services.InitServices()
	log.Println("[INFO] Services initialized")

	addr := ":8080"
	if web.BConfig.Listen.HTTPPort != 0 {
		addr = fmt.Sprintf(":%d", web.BConfig.Listen.HTTPPort)
	}
	log.Printf("[INFO] TravelSphere starting on %s", addr)

	// Handle graceful shutdown on Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("[INFO] Received signal: %v", sig)
		log.Println("[INFO] Shutting down gracefully...")
		os.Exit(0)
	}()

	web.Run()
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
