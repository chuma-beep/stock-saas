package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/chuma-beep/stock-saas/internal/database"
	"github.com/gin-gonic/gin"
)

type Feedback struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Feedback  string    `json:"feedback"`
	CreatedAt time.Time `json:"created_at"`
}

func SubmitFeedback(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Feedback string `json:"feedback"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Feedback == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feedback is required"})
		return
	}

	query := `
		INSERT INTO feedback (name, email, feedback)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	var id int
	var createdAt time.Time
	err := database.DB.QueryRow(query, req.Name, req.Email, req.Feedback).Scan(&id, &createdAt)
	if err != nil {
		log.Printf("Error saving feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Feedback submitted successfully",
		"id":      id,
	})
}

func GetFeedback(c *gin.Context) {
	query := `
		SELECT id, name, email, feedback, created_at
		FROM feedback
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feedback"})
		return
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var f Feedback
		err := rows.Scan(&f.ID, &f.Name, &f.Email, &f.Feedback, &f.CreatedAt)
		if err != nil {
			continue
		}
		feedbacks = append(feedbacks, f)
	}

	c.JSON(http.StatusOK, gin.H{"feedbacks": feedbacks})
}
