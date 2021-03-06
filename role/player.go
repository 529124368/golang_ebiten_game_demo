package role

import (
	"embed"
	"game/tools"
	"image"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/fzipp/texturepacker"
	"github.com/hajimehoshi/ebiten/v2"
)

var plist_sheet, plist_wea_sheet, plist_skill_sheet *texturepacker.SpriteSheet
var plist_png, plist_wea_png, plist_skill_png *image.NRGBA
var loadedSkill string

type Player struct {
	X         float64
	Y         float64
	State     int
	Direction int
	MouseX    int
	MouseY    int
	SkillName string
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
		SkillName: "",
		image:     images,
	}
	return play
}
func (p *Player) LoadImages() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	//player load
	go func() {
		plist, _ := p.image.ReadFile("resource/man/warrior/man.png")
		plist_json, _ := p.image.ReadFile("resource/man/warrior/man.json")
		plist_sheet, plist_png = tools.GetImageFromPlist(plist, plist_json)
		runtime.GC()
		wg.Done()
	}()
	//weapon laod
	go func() {
		plist, _ := p.image.ReadFile("resource/man/weapon/tulong.png")
		plist_json, _ := p.image.ReadFile("resource/man/weapon/tulong.json")
		plist_wea_sheet, plist_wea_png = tools.GetImageFromPlist(plist, plist_json)
		runtime.GC()
		wg.Done()

	}()
	//skill load
	go func() {
		loadedSkill = "liehuo"
		plist, _ := p.image.ReadFile("resource/man/skill/liehuo.png")
		plist_json, _ := p.image.ReadFile("resource/man/skill/liehuo.json")
		plist_skill_sheet, plist_skill_png = tools.GetImageFromPlist(plist, plist_json)
		runtime.GC()
		wg.Done()
	}()
	wg.Wait()
	p.SetPlayerState(0, 0)
	go func() {
		runtime.GC()
	}()
}

func (p *Player) loadSkillImages(name string) {
	go func() {
		loadedSkill = name
		plist, _ := p.image.ReadFile("resource/man/skill/" + name + ".png")
		plist_json, _ := p.image.ReadFile("resource/man/skill/" + name + ".json")
		plist_skill_sheet, plist_skill_png = tools.GetImageFromPlist(plist, plist_json)
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
		if p.SkillName != loadedSkill {
			p.loadSkillImages(p.SkillName)
		}
		xy := strings.Split(plist_skill_sheet.Meta.Version, "_")
		x, _ := strconv.Atoi(xy[0])
		y, _ := strconv.Atoi(xy[1])
		return ebiten.NewImageFromImage(plist_skill_png.SubImage(plist_skill_sheet.Sprites[name].Frame)), x, y
	}

}
