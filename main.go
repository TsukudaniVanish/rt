package main
import (
	"log"
	"context"


	"github.com/anaseto/gruid"
	sdl "github.com/anaseto/gruid-sdl"
)

func main() {
	gd := gruid.NewGrid(UIWidth, UIHight)
	m := &Model{ Grid: gd}
	// Specify a driver among the provided ones.
	tile, err := GetTileDrawer()
	if err != nil {
		log.Fatal(err)
	}
	dr := sdl.NewDriver(sdl.Config{
		TileManager: tile,
	})

	app := gruid.NewApp(gruid.AppConfig{
		Driver: dr,
		Model: m,
	})
	// Start the main loop of the application.
	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
