package urlshortener

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type request struct {
	URL string `json:"url"`
}

type response struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func GetHandlers(storage Storage) (func(c *gin.Context), func(c *gin.Context)) {
	shorten := func(c *gin.Context) {
		shortenHandler(c, storage)
	}
	redirect := func(c *gin.Context) {
		redirectHandler(c, storage)
	}

	return shorten, redirect
}

// Expected to get url from body in json format {url: desired_url}
func shortenHandler(c *gin.Context, storage Storage) {
	var req request
	err := json.NewDecoder(c.Request.Body).Decode(&req)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to decode body: %v", err))
		return
	}

	if req.URL == "" {
		c.String(http.StatusBadRequest, "Missing URL")
		return
	}

	key := storage.Store(req.URL)

	c.JSON(http.StatusOK, response{
		URL: req.URL,
		Key: key,
	})
}

// Expected to get key from path
func redirectHandler(c *gin.Context, storage Storage) {
	key := c.Param("key")

	url, found := storage.Get(key)

	if !found {
		c.String(http.StatusBadRequest, "Key not found")
		return
	}

	c.Redirect(http.StatusFound, url)
}
