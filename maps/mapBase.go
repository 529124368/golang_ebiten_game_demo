package maps

import (
	"embed"
	"game/tools"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

type MapBase struct {
	image   *embed.FS
	OpBg    *ebiten.DrawImageOptions
	BgImage *ebiten.Image
}

//init
func NewMap(images *embed.FS) *MapBase {
	maps := &MapBase{
		image: images,
	}
	return maps
}

//load images
func (m *MapBase) LoadMap() {
	//BG
	go func() {
		s2, _ := m.image.ReadFile("resource/bg/bg1.png")
		img := tools.GetEbitenImage(s2)
		m.BgImage = img
		m.OpBg = &ebiten.DrawImageOptions{}
		m.OpBg.Filter = ebiten.FilterLinear
		m.OpBg.GeoM.Translate(-700, -550)
		runtime.GC()
	}()
}
