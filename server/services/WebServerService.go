package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func generateTaskView(c *gin.Context) {
	taskSize, isFound := c.Params.Get("taskSize")
	if isFound {
		ts, err := strconv.Atoi(taskSize)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "taskSize must be integer"})
			return
		}
		task := GenerateTask(ts)
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Task with ID %d sucessfully generated", task.ID)})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "taskSize not found"})
	return
}

func StartWebServerListener() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/healthCheck", healthCheck)
	r.GET("/generateTask/:taskSize", generateTaskView)
	fmt.Println("http listener started")
	err := r.Run(":3333")
	FailOnError(err, "Cant start http listener")
}
