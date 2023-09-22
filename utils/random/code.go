package random

import (
	"fmt"
	"golang.org/x/exp/rand"
	"time"
)

var stringCode = ""

func Code(length int) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	return fmt.Sprintf("%4v", rand.Intn(10000))
}
