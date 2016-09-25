// By izanbf1803
// http://izanbf.es/

package err

import (
	"fmt"
	"os"
)

// Log : Print error to console
func Log(txt string, a ...interface{}) {
	msg := fmt.Sprintf(txt, a...)
	fmt.Printf("\n\n\tERROR: %v\n\n", msg)
}

// Exit : Print error to console and stop program execution [exit(1)]
func Exit(txt string, a ...interface{}) {
	Log(txt, a...)
	os.Exit(0)
}
