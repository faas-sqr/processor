package signals

import "os"

var Sigs chan os.Signal

func init() {
	Sigs = make(chan os.Signal, 1)
}
