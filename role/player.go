package role

import (
	"embed"
	"game/tools"
	"image"
	"runtime"
	"sync"

	"github.com/fzipp/texturepacker"
	"github.com/hajimehoshi/ebiten/v2"
)

var plist_sheet, plist_wea_sheet, plist_skill_sheet *texturepacker.SpriteSheet
var plist_png, plist_wea_png, plist_skill_png *image.NRGBA

type Player struct {
	X         float64
	Y         float64
	State     int
	Direction int
	MouseX    int
	MouseY    int
	image     *embed.FS
}

func NewPlayer(x, y float64, state, dir, mx, my int, images *embed.FS) *Player {
	play := &Player{
		X:         x,
		Y:         y,
		State:     state,
		Direction: dir,
		MouseX:    mx,
		MouseY:    my,
		image:     images,
	}
	return play
}
func (p *Player) LoadImages() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		plist, _ := p.image.ReadFile("resource/man/warrior/man.png")
		plist_json, _ := p.image.ReadFile("resource/man/warrior/man.json")
		plist_sheet, plist_png = tools.GetImageFromPlist(plist, plist_json)
		go func() {
			runtime.GC()
		}()
		wg.Done()
	}()
	go func() {
		plist, _ := p.image.ReadFile("resource/man/weapon/tulong.png")
		plist_json, _ := p.image.ReadFile("resource/man/weapon/tulong.json")
		plist_wea_sheet, plist_wea_png = tools.GetImageFromPlist(plist, plist_json)
		go func() {
			runtime.GC()
		}()
		wg.Done()

	}()
	go func() {
		plist, _ := p.image.ReadFile("resource/man/skill/skill.png")
		plist_json, _ := p.image.ReadFile("resource/man/skill/skill.json")
		plist_skill_sheet, plist_skill_png = tools.GetImageFromPlist(plist, plist_json)
		go func() {
			runtime.GC()
		}()
		wg.Done()

	}()
	wg.Wait()
	p.SetPlayerState(0, 0)
	go func() {
		runtime.GC()
	}()
}

//set state
func (p *Player) SetPlayerState(s, d int) {
	p.State = s
	p.Direction = d
}

//TODO
func (p *Player) Attack() {

}

//TODO
func (p *Player) DeadEvent() {

}

func (p *Player) GetAnimator(flg, name string) (*ebiten.Image, int, int) {
	if flg == "man" {
		return ebiten.NewImageFromImage(plist_png.SubImage(plist_sheet.Sprites[name].Frame)), plist_sheet.Sprites[name].SpriteSourceSize.Min.X, plist_sheet.Sprites[name].SpriteSourceSize.Min.Y
	}
	if flg == "weapon" {
		return ebiten.NewImageFromImage(plist_wea_png.SubImage(plist_wea_sheet.Sprites[name].Frame)), plist_wea_sheet.Sprites[name].SpriteSourceSize.Min.X, plist_wea_sheet.Sprites[name].SpriteSourceSize.Min.Y
	} else {
		return ebiten.NewImageFromImage(plist_skill_png.SubImage(plist_skill_sheet.Sprites[name].Frame)), 0, 0
	}

}
