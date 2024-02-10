package runner

import "fmt"

const version = "v0.0.1"

func ShowBanner() {
	//http://www.network-science.de/ascii/  smslant
	var banner = `banner %s`
	fmt.Printf(banner, version)
}