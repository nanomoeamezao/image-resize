package flags

import "flag"

func ReadFlags() int {
	flagLimit := flag.Int("limit", 2, "amount of parallel requests allowed")
	flag.Parse()
	return *flagLimit
}
