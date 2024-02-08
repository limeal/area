package context

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/websocket/v2"
)

// It reads the log file of the applet, and sends the new lines to the client
func GetAppletLogs(c *websocket.Conn) {

	appletID := c.Params("applet_id", "")

	var nbLines int
	for {
		time.Sleep(2 * time.Second)

		content, err := os.ReadFile("logs/" + appletID + ".log")
		if err != nil {
			fmt.Printf("Error while reading: %v. Closing http connection.\n", err)
			break
		}

		lines := strings.Split(string(content), "\n")
		if nbLines == 0 {
			nbLines = len(lines) - 1
		} else if len(lines) == nbLines {
			continue
		}

		nbNbLines := len(lines)
		diff := nbNbLines - nbLines
		// Get the n last lines (where n is the diff)
		content = []byte(strings.Join(lines[len(lines)-diff:], "\n"))
		if len(content) == 0 {
			continue
		}

		nbLines = nbNbLines
		if err = c.WriteMessage(websocket.TextMessage, content); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
