package digitalclock

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	timeString, err := parseTime(c.Query("time"))
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to parse time from query: %v", err))
		return
	}

	k, err := parseK(c.Query("k"))
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to parse k from query: %v", err))
		return
	}

	img := formImage(k, timeString)

	c.Header("Content-Type", "image/png")

	png.Encode(c.Writer, img)
}
