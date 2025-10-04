package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ingridhq/zebrash"
	"github.com/ingridhq/zebrash/drawers"
)

func main() {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.Default())

	// Health check endpoint
	r.GET("/health", healthHandler)
	
	// ZPL rendering endpoint
	r.POST("/render/:x/:y/:dpmm", renderHandler)
	
	port := ":3009"
	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(port))
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func renderHandler(c *gin.Context) {
	// Extract URL parameters
	x, err := strconv.Atoi(c.Param("x"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid x parameter"})
		return
	}
	
	y, err := strconv.Atoi(c.Param("y"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid y parameter"})
		return
	}
	
	dpmm, err := strconv.Atoi(c.Param("dpmm"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dpi parameter"})
		return
	}

	// Read ZPL data from request body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty ZPL data"})
		return
	}
	
	// Parse ZPL data using zebrash
	parser := zebrash.NewParser()

	res, err := parser.Parse(body)
	if err != nil {
		log.Printf("Error parsing ZPL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse ZPL: %v", err)})
		return
	}
	
	if len(res) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No labels found in ZPL data"})
		return
	}
	
	// Create a buffer to hold the PNG data
	var buff bytes.Buffer
	
	// Create drawer and render to PNG
	drawer := zebrash.NewDrawer()
	
	// x and y are already in millimeters
	labelWidthMm := float64(x)   // x is in mm
	labelHeightMm := float64(y)  // y is in mm
	dpmmi := int(dpmm)  // Convert DPI to dots per mm

	err = drawer.DrawLabelAsPng(res[0], &buff, drawers.DrawerOptions{
		LabelWidthMm:  labelWidthMm,
		LabelHeightMm: labelHeightMm,
		Dpmm:          dpmmi,
	})
	if err != nil {
		log.Printf("Error rendering ZPL to PNG: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to render ZPL to PNG: %v", err)})
		return
	}

	pngData := buff.Bytes()

	// Set appropriate headers for PNG response
	c.Header("Content-Type", "image/png")
	c.Header("Content-Length", strconv.Itoa(len(pngData)))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	
	// Write PNG data to response
	c.Data(http.StatusOK, "image/png", pngData)
	
	log.Printf("Successfully rendered PNG (%d bytes)", len(pngData))
}