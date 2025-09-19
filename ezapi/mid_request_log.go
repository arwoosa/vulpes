package ezapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 印出所有 Header
		var buffer bytes.Buffer
		buffer.WriteString("===== Headers =====\n")
		for k, v := range c.Request.Header {
			buffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}

		// 如果是 POST，印出 Body
		if c.Request.Method == http.MethodPost {
			buffer.WriteString("===== Body =====\n")
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				buffer.WriteString(string(bodyBytes))
				// 讀過後要重設 body，讓後續 handler 還能用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}
		fmt.Println(strings.TrimSpace(buffer.String()))
		c.Next()
	}
}
