package main

import (
	"hrm-backend/auth"
	"hrm-backend/config"
	"hrm-backend/routes"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	if err := auth.Initialize(); err != nil {
		log.Fatalf("Authentication configuration error: %v", err)
	}
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Unable to configure trusted proxies: %v", err)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// 2. ඩේටාබේස් එක කනෙක්ට් කිරීම
	config.ConnectDatabase()

	// 3. පාරවල් (Routes) සෙට් කිරීම
	routes.SetupRoutes(r)

	// 4. සර්වර් එක පෝට් 8080 එකේ ඕන් කිරීම
	r.Run(":8080")
}
