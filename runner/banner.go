package runner

import "fmt"

const VERSION = "v0.2.3"

func ShowBanner() {
	//http://www.network-science.de/ascii/  smslant
	var banner = `
   ___   __   _____  _____                     _       _       
  / __\ / /   \_   \/__   \___ _ __ ___  _ __ | | __ _| |_ ___ 
 / /   / /     / /\/  / /\/ _ \ '_ ' _ \| '_ \| |/ _' | __/ _ \
/ /___/ /___/\/ /_   / / |  __/ | | | | | |_) | | (_| | ||  __/
\____/\____/\____/   \/   \___|_| |_| |_| .__/|_|\__,_|\__\___|
					|_|                %s
`
	fmt.Printf(banner, VERSION)
}
