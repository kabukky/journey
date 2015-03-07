package timer

import (
	"log"
	"time"
)

func Track(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
