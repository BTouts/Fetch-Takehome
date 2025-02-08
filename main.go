package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var receiptStore = make(map[string]Receipt)
var mutex = sync.Mutex{}

func main() {
	router := gin.Default()

	// Endpoints
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getReceiptPoints)

	router.Run(":8080")
}

func processReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid."})
		return
	}
	receipt.ID = uuid.New().String()

	// Store in memory
	mutex.Lock()
	receiptStore[receipt.ID] = receipt
	mutex.Unlock()

	// Return ID
	c.JSON(http.StatusOK, gin.H{"id": receipt.ID})
}

func getReceiptPoints(c *gin.Context) {
	id := c.Param("id")
	receipt, ok := receiptStore[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "No receipt found for that ID."})
		return
	}

	// Calculate points
	points := calculatePoints(receipt)
	c.JSON(http.StatusOK, gin.H{"points": points})
}

func calculatePoints(receipt Receipt) int {
	points := 0

	// 1️⃣ One point for every alphanumeric character in the retailer name.
	for _, char := range receipt.Retailer {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			points++
		}
	}

	// Convert total to float
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		fmt.Println("Error parsing total:", err)
		return points // Return existing points if there's an error
	}

	// 50 points if the total is a round dollar amount (no cents)
	if totalFloat == float64(int(totalFloat)) {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25
	if int(totalFloat*100)%25 == 0 {
		points += 25
	}

	// 6 points if the day in the purchase date is odd
	dateParts := strings.Split(receipt.PurchaseDate, "-") // Split "YYYY-MM-DD"
	if len(dateParts) == 3 {
		day, err := strconv.Atoi(dateParts[2]) // Convert "DD" to int
		if err == nil && day%2 != 0 {          // Check if the day is odd
			points += 6
		}
	}

	// 10 points if the time of purchase is between 2:00 PM - 4:00 PM
	timeParts := strings.Split(receipt.PurchaseTime, ":") // Split "HH:MM"
	if len(timeParts) == 2 {
		hour, err := strconv.Atoi(timeParts[0]) // Convert "HH" to int
		if err == nil && hour >= 14 && hour < 16 {
			points += 10
		}
	}

	// 5 points for every 2 items on the receipt
	itemCount := len(receipt.Items)
	points += (itemCount / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, multiply price by 0.2, round up, and add that to points.
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				itemPoints := int(math.Ceil(priceFloat * 0.2)) // Ensure rounding up
				points += itemPoints
			}
		}
	}

	return points
}
