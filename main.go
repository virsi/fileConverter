package main

import (
	"fmt"

	"github.com/virsi/fileConverter/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Print(cfg)

	// TODO init logger
	// TODO init storage
	// TODO init router
	// TODO start server
}
