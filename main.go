package main

import (
    "os"
	"fmt"

	"github.com/mikeunge/Scripts/wallpaper-engine/cmd"
)

func main() {
    err := cmd.App()
    if err != nil {
        fmt.Printf("Wallpaper-Engine exited with error %+v\n", err)
        os.Exit(1)
    }
    os.Exit(0)
}
