package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// User is the API representation of a user.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// users is intentionally in memory for this first Gin exercise.
// The MySQL lesson will replace it with persistent storage.
var users = []User{
	{ID: 1, Name: "Alice", Email: "alice@example.com"},
	{ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func main() {
	router := gin.Default()

	router.GET("/health", healthCheck)
	router.GET("/users", listUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", createUser)
	router.PATCH("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func listUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func getUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
}

func createUser(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name 必填，email 必须是有效邮箱"})
		return
	}

	user := User{
		ID:    nextUserID(),
		Name:  input.Name,
		Email: input.Email,
	}
	users = append(users, user)
	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	var input struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体必须是合法 JSON"})
		return
	}
	if input.Name == nil && input.Email == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "至少提供 name 或 email"})
		return
	}

	for i := range users {
		if users[i].ID == id {
			if input.Name != nil {
				users[i].Name = *input.Name
			}
			if input.Email != nil {
				users[i].Email = *input.Email
			}
			c.JSON(http.StatusOK, users[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
}

func deleteUser(c *gin.Context) {
	id, ok := parseUserID(c)
	if !ok {
		return
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
}

func parseUserID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id 必须是正整数"})
		return 0, false
	}
	return id, true
}

func nextUserID() int {
	maxID := 0
	for _, user := range users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}
	return maxID + 1
}
