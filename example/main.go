package main

import (
	"fmt"
	"log"
	"os"

	"github.com/funatsufumiya/ebiten_gvvideo/gvplayer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player *gvplayer.GVPlayer
	err    error
}

func (g *Game) Update() error {
	if g.err != nil {
		return g.err
	}
	if err := g.player.Update(); err != nil {
		g.err = err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.err != nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(20, 20)
	g.player.Draw(screen, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	var gvPath string
	if len(os.Args) > 1 {
		gvPath = os.Args[1]
	} else {
		gvPath = "example/test_asset/test-10px.gv"
		fmt.Println("[INFO] Playing the default GV video. You can specify a .gv file as an argument.")
	}
	player, err := gvplayer.NewGVPlayer(gvPath)
	if err != nil {
		log.Fatal(err)
	}
	g := &Game{player: player}
	player.Play()

	ebiten.SetWindowTitle("GV Video (Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
