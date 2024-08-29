package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func convertFromB64() {

}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.File("./build/index.html")
	})

	router.POST("/ingest_picture", func(c *gin.Context) {

		type pictureRequest struct {
			Image string `json:"image"`
		}

		var jsonData pictureRequest

		if err := c.BindJSON(&jsonData); err != nil {
			// If there is an error, respond with bad request
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := base64.StdEncoding.DecodeString(jsonData.Image)
		if err != nil {
			fmt.Printf("Invalid b64 image: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		requestURL := "http://localhost:999/captionImage"
		req, err := http.Post(requestURL, "application/json", strings.NewReader(fmt.Sprintf(`{"image":"%s"}`, jsonData.Image)))
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		finalB64Image := jsonData.Image

		// if req.StatusCode != 200 {
		// 	fmt.Printf("Didn't get 200 response from backend: (%s)\n", req.Status)
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error querying /captionImage"})
		// }

		fmt.Printf("client: got response!\n")
		fmt.Printf("Response status code: %d\n", req.StatusCode)
		fmt.Printf("Response body: %s\n", string(body))
		c.JSON(http.StatusOK, gin.H{"success": finalB64Image})

	})

	router.Run("0.0.0.0:999")
}
