package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/funatsufumiya/ebiten_gvvideo/gvplayer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MultiGame struct {
	players      []*gvplayer.GVPlayer
	errs         []error
	windowWidth  int
	windowHeight int
	async        bool
	gvPaths      []string
	startTimes   []time.Time
	loop         bool
}

func (mg *MultiGame) Update() error {
    if inpututil.IsKeyJustPressed(ebiten.KeyA) {
        mg.toggleAsync()
    }
	for i, player := range mg.players {
		if mg.errs[i] != nil {
			continue
		}
		if err := player.Update(); err != nil {
			mg.errs[i] = err
		}
	}
	return nil
}

func (mg *MultiGame) toggleAsync() {
	mg.async = !mg.async
	for i, path := range mg.gvPaths {
		mg.players[i].Stop()
		player, err := gvplayer.NewGVPlayerWithOption(path, mg.async)
		if err != nil {
			mg.errs[i] = err
			continue
		}
		player.Loop = mg.loop
		mg.players[i] = player
		mg.startTimes[i] = time.Now()
		player.Play()
	}
}

func (mg *MultiGame) Draw(screen *ebiten.Image) {
	n := len(mg.players)
	if n == 0 {
		return
	}
	cols := int(1)
	for cols*cols < n {
		cols++
	}
	rows := (n + cols - 1) / cols
	w := mg.windowWidth / cols
	h := mg.windowHeight / rows
	for i, player := range mg.players {
		if mg.errs[i] != nil {
			continue
		}
		row := i / cols
		col := i % cols
		op := &ebiten.DrawImageOptions{}
		videoW := player.Width()
		videoH := player.Height()
		scaleX := float64(w) / float64(videoW)
		scaleY := float64(h) / float64(videoH)
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}
		op.GeoM.Scale(scale, scale)
		tx := float64(col*w) + float64(w)/2 - float64(videoW)*scale/2
		ty := float64(row*h) + float64(h)/2 - float64(videoH)*scale/2
		op.GeoM.Translate(tx, ty)
		player.Draw(screen, op)
		videoTime := player.CurrentTime().Seconds()
		elapsed := time.Since(mg.startTimes[i]).Seconds()
		msg := fmt.Sprintf("Video %d: %.2fs | Elapsed: %.2fs", i+1, videoTime, elapsed)
		ebitenutil.DebugPrintAt(screen, msg, col*w, row*h+16)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f | Async: %v", ebiten.ActualFPS(), mg.async))
}

func (mg *MultiGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	mg.windowWidth = outsideWidth
	mg.windowHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	var gvPaths []string
	loop := true
	args := os.Args[1:]
	if len(args) > 0 {
		gvPaths = args
	} else {
		gvPaths = []string{"example/test_asset/test-10px.gv", "example/test_asset/test-10px.gv"}
		fmt.Println("[INFO] Playing default GV videos. You can specify multiple .gv files as arguments.")
	}
	players := make([]*gvplayer.GVPlayer, len(gvPaths))
	errs := make([]error, len(gvPaths))
	startTimes := make([]time.Time, len(gvPaths))
	for i, path := range gvPaths {
		player, err := gvplayer.NewGVPlayerWithOption(path, true)
		if err != nil {
			errs[i] = err
			continue
		}
		player.Loop = loop
		players[i] = player
		startTimes[i] = time.Now()
		player.Play()
	}
	mg := &MultiGame{players: players, errs: errs, async: true, gvPaths: gvPaths, startTimes: startTimes, loop: loop}

	ebiten.SetWindowTitle("GV Video Multiple (Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(mg); err != nil {
		log.Fatal(err)
	}
}
