package role

import (
	"embed"
	"game/tools"
	"runtime"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	X           float64
	Y           float64
	State       int
	Direction   int
	MouseX      int
	MouseY      int
	ImgNow      map[int]*ebiten.Image
	ImgWeaNow   map[int]*ebiten.Image
	ImgSkillNow map[int]*ebiten.Image
	asset       map[string]*ebiten.Image //man
	assetWea    map[string]*ebiten.Image //weapon
	skill_asset map[string]*ebiten.Image //skill
	image       *embed.FS
}

func NewPlayer(x, y float64, state, dir, mx, my int, images *embed.FS) *Player {
	play := &Player{
		X:           x,
		Y:           y,
		State:       state,
		Direction:   dir,
		MouseX:      mx,
		MouseY:      my,
		ImgNow:      make(map[int]*ebiten.Image, 8),
		ImgWeaNow:   make(map[int]*ebiten.Image, 8),
		ImgSkillNow: make(map[int]*ebiten.Image, 6),
		asset:       make(map[string]*ebiten.Image, 160),
		assetWea:    make(map[string]*ebiten.Image, 160),
		skill_asset: make(map[string]*ebiten.Image, 48),
		image:       images,
	}
	return play
}
func (p *Player) LoadImages() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	//load skill images
	go func() {
		for j := 0; j < 8; j++ {
			for i := 0; i < 6; i++ {
				s, _ := p.image.ReadFile("resource/man/skill/" + strconv.Itoa(j) + "_skill_" + strconv.Itoa(i) + ".png")
				mg := tools.GetEbitenImage(s)
				name := strconv.Itoa(j) + "_" + strconv.Itoa(i)
				p.skill_asset[name] = mg
			}
		}
		go func() {
			runtime.GC()
		}()
		wg.Done()
	}()
	go func() {
		//load idle images
		for j := 0; j < 8; j++ {
			for i := 0; i < 6; i++ {
				s, _ := p.image.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_stand_" + strconv.Itoa(i) + ".png")
				mg := tools.GetEbitenImage(s)
				name := "ms_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				p.asset[name] = mg
				s, _ = p.image.ReadFile("resource/man/weapon/" + strconv.Itoa(j) + "_stand_" + strconv.Itoa(i) + ".png")
				mg = tools.GetEbitenImage(s)
				name = "ws_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				p.assetWea[name] = mg
				s, _ = p.image.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_attack_" + strconv.Itoa(i) + ".png")
				mg = tools.GetEbitenImage(s)
				name = "ma_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				p.asset[name] = mg
				s, _ = p.image.ReadFile("resource/man/weapon/" + strconv.Itoa(j) + "_attack_" + strconv.Itoa(i) + ".png")
				mg = tools.GetEbitenImage(s)
				name = "wa_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				p.assetWea[name] = mg
			}
		}
		runtime.GC()
		wg.Done()
	}()

	go func() {
		for j := 0; j < 8; j++ {
			for i := 0; i < 8; i++ {
				//man
				s, _ := p.image.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_run_" + strconv.Itoa(i) + ".png")
				mg := tools.GetEbitenImage(s)
				name := "mr_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				p.asset[name] = mg

				//weapon
				s1, _ := p.image.ReadFile("resource/man/weapon/" + strconv.Itoa(j) + "_run_" + strconv.Itoa(i) + ".png")
				mg = tools.GetEbitenImage(s1)
				name = "wr_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				p.assetWea[name] = mg
			}
		}
		runtime.GC()
		wg.Done()
	}()
	wg.Wait()
	runtime.GC()
}

//set Anm
func (p *Player) SetAnimator(i int) {
	if i == 0 {
		tools.MapCopy(p.skill_asset, p.ImgSkillNow, p.State, p.Direction, "skill")
	}
	tools.MapCopy(p.asset, p.ImgNow, p.State, p.Direction, "man")
	tools.MapCopy(p.assetWea, p.ImgWeaNow, p.State, p.Direction, "weapon")
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
