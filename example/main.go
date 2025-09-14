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
	player       *gvplayer.GVPlayer
	err          error
	windowWidth  int
	windowHeight int
	async        bool
	gvPath       string
}

func (g *Game) Update() error {
	if g.err != nil {
		return g.err
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.toggleAsync()
	}
	if err := g.player.Update(); err != nil {
		g.err = err
	}
	return nil
}

func (g *Game) toggleAsync() {
	g.async = !g.async
	g.player.Stop()
	player, err := gvplayer.NewGVPlayerWithOption(g.gvPath, g.async)
	if err != nil {
		g.err = err
		return
	}
	g.player = player
	g.player.Play()
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.err != nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	// scale calculation
	videoW := g.player.Width()
	videoH := g.player.Height()
	scaleX := float64(g.windowWidth) / float64(videoW)
	scaleY := float64(g.windowHeight) / float64(videoH)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	op.GeoM.Scale(scale, scale)
	// position center
	tx := float64(g.windowWidth)/2 - float64(videoW)*scale/2
	ty := float64(g.windowHeight)/2 - float64(videoH)*scale/2
	op.GeoM.Translate(tx, ty)
	g.player.Draw(screen, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f | Async: %v (Key A to toggle)", ebiten.ActualFPS(), g.async))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.windowWidth = outsideWidth
	g.windowHeight = outsideHeight
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
	player, err := gvplayer.NewGVPlayerWithOption(gvPath, true) // デフォルト非同期
	if err != nil {
		log.Fatal(err)
	}
	g := &Game{player: player, async: true, gvPath: gvPath}
	player.Play()

	ebiten.SetWindowTitle("GV Video (Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
