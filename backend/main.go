package main

import (
	"hrm-backend/config"
	"hrm-backend/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 1. CORS Setup (Frontend එකට බාධාවක් නැතුව කතා කරන්න ඉඩ දීම)
	// For local development allow requests from the frontend dev server.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
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
