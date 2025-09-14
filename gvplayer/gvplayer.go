package gvplayer

import (
	"time"

	"github.com/funatsufumiya/go-gv-video/gvvideo"
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState int

const (
	Stopped PlayerState = iota
	Playing
	Paused
)

type GVPlayer struct {
	video      *gvvideo.GVVideo
	frameImage *ebiten.Image
	state      PlayerState
	startTime  time.Time
	pauseTime  time.Time
	seekTime   time.Duration
	loop       bool
}

func NewGVPlayer(path string) (*GVPlayer, error) {
	video, err := gvvideo.LoadGVVideo(path)
	if err != nil {
		return nil, err
	}
	img := ebiten.NewImage(int(video.Header.Width), int(video.Header.Height))
	return &GVPlayer{
		video:      video,
		frameImage: img,
		state:      Stopped,
	}, nil
}

func (p *GVPlayer) Play() {
	if p.state == Playing {
		return
	}
	p.state = Playing
	p.startTime = time.Now()
}

func (p *GVPlayer) Pause() {
	if p.state != Playing {
		return
	}
	p.state = Paused
	p.pauseTime = time.Now()
}

func (p *GVPlayer) Stop() {
	p.state = Stopped
	p.seekTime = 0
}

func (p *GVPlayer) Seek(to time.Duration) {
	p.seekTime = to
}

func (p *GVPlayer) Update() error {
	if p.state != Playing {
		return nil
	}
	elapsed := time.Since(p.startTime) + p.seekTime
	frameID := uint32(elapsed.Seconds() * float64(p.video.Header.FPS))
	if frameID >= p.video.Header.FrameCount {
		if p.loop {
			p.startTime = time.Now()
			p.seekTime = 0
			frameID = 0
		} else {
			p.state = Stopped
			return nil
		}
	}
	frame, err := p.video.ReadFrame(frameID)
	if err != nil {
		return err
	}
	p.frameImage.WritePixels(frame)
	return nil
}

func (p *GVPlayer) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if p.frameImage != nil {
		screen.DrawImage(p.frameImage, op)
	}
}
