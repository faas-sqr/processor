package signals

import "os"

var Sigs chan os.Signal

var FunctionSigs1 chan os.Signal

var FunctionSigs2 chan os.Signal

func init() {
	Sigs = make(chan os.Signal, 1)
	FunctionSigs1 = make(chan os.Signal, 1)
	FunctionSigs2 = make(chan os.Signal, 1)
}
