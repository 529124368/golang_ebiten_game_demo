package main

import (
	"embed"
	"fmt"
	"game/tools"
	_ "image/png"
	"log"
	"runtime"
	"strconv"
	"sync"

	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	IDLE          int     = 0
	RUN           int     = 1
	ATTACK        int     = 2
	SCREENWIDTH   int     = 450
	SCREENHEIGHT  int     = 300
	OFFSETX       int     = 200
	OFFSETY       int     = 80
	PLAYERCENTERX int64   = 361
	PLAYERCENTERY int64   = 219
	SPEED         float64 = 2
	LAYOUTX       int     = 720
	LAYOUTY       int     = 480
	WEOFFSETX     int     = 127
	WEOFFSETY     int     = 14
	SKILLOFFSETX  int     = 120
	SKILLOFFSETY  int     = 10
)

//Game
type Game struct {
	count  int
	player *Player
}

//Player
type Player struct {
	x         float64
	y         float64
	state     int
	direction int
	MouseX    int
	MouseY    int
}

var (
	counts      int  = 0
	frameNums   int  = 4
	flg         bool = false
	imgNow      map[int]*ebiten.Image
	imgWeaNow   map[int]*ebiten.Image
	imgSkillNow map[int]*ebiten.Image
)
var op, opWea, opBg, opUI, opSkill *ebiten.DrawImageOptions

//man
var asset map[string]*ebiten.Image

//weapon
var assetWea map[string]*ebiten.Image

//skill
var skill_asset map[string]*ebiten.Image

var bgImage, UI *ebiten.Image

//go:embed resource
var images embed.FS

//factory
func NewGame() *Game {
	gameStart := &Game{
		count: 0,
		player: &Player{
			x:         float64(SCREENWIDTH / 2),
			y:         float64(SCREENHEIGHT / 2),
			state:     IDLE,
			direction: 0,
			MouseX:    0,
			MouseY:    0,
		},
	}
	return gameStart
}

func init() {
	wg := sync.WaitGroup{}
	wg.Add(4)
	//get images
	imgNow = make(map[int]*ebiten.Image, 8)
	imgWeaNow = make(map[int]*ebiten.Image, 8)
	imgSkillNow = make(map[int]*ebiten.Image, 6)

	//man
	asset = make(map[string]*ebiten.Image, 160)

	//weapon
	assetWea = make(map[string]*ebiten.Image, 160)

	//skill
	skill_asset = make(map[string]*ebiten.Image, 48)

	//UI load
	s, _ := images.ReadFile("resource/UI/attack.png")
	mgUI := tools.GetEbitenImage(s)
	UI = mgUI

	//load skill images
	go func() {
		for j := 0; j < 8; j++ {
			for i := 0; i < 6; i++ {
				s, _ := images.ReadFile("resource/man/skill/" + strconv.Itoa(j) + "_skill_" + strconv.Itoa(i) + ".png")
				mg := tools.GetEbitenImage(s)
				name := strconv.Itoa(j) + "_" + strconv.Itoa(i)
				skill_asset[name] = mg
			}
		}
		wg.Done()
	}()
	go func() {
		//load idle images
		for j := 0; j < 8; j++ {
			for i := 0; i < 6; i++ {
				s, _ := images.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_stand_" + strconv.Itoa(i) + ".png")
				mg := tools.GetEbitenImage(s)
				name := "ms_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				asset[name] = mg
				s, _ = images.ReadFile("resource/man/weapon/" + strconv.Itoa(j) + "_stand_" + strconv.Itoa(i) + ".png")
				mg = tools.GetEbitenImage(s)
				name = "ws_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				assetWea[name] = mg
				s, _ = images.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_attack_" + strconv.Itoa(i) + ".png")
				mg = tools.GetEbitenImage(s)
				name = "ma_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				asset[name] = mg
				s, _ = images.ReadFile("resource/man/weapon/" + strconv.Itoa(j) + "_attack_" + strconv.Itoa(i) + ".png")
				mg = tools.GetEbitenImage(s)
				name = "wa_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				assetWea[name] = mg
			}
		}
		wg.Done()
	}()

	go func() {
		for j := 0; j < 8; j++ {
			for i := 0; i < 8; i++ {
				//man
				s, _ := images.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_run_" + strconv.Itoa(i) + ".png")
				mg := tools.GetEbitenImage(s)
				name := "mr_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				asset[name] = mg

				//weapon
				s1, _ := images.ReadFile("resource/man/weapon/" + strconv.Itoa(j) + "_run_" + strconv.Itoa(i) + ".png")
				mg = tools.GetEbitenImage(s1)
				name = "wr_" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				assetWea[name] = mg
			}
		}
		wg.Done()
	}()

	go func() {
		//BG
		s2, _ := images.ReadFile("resource/bg/bg1.png")
		img := tools.GetEbitenImage(s2)
		bgImage = img
		opBg = &ebiten.DrawImageOptions{}
		opBg.Filter = ebiten.FilterLinear
		opBg.GeoM.Translate(-700, -550)
		//player option
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(SCREENWIDTH/2+OFFSETX), float64(SCREENHEIGHT/2+OFFSETY))
		op.GeoM.Scale(0.7, 0.7)
		op.Filter = ebiten.FilterLinear
		//weapon option
		opWea = &ebiten.DrawImageOptions{}
		opWea.GeoM.Translate(float64(SCREENWIDTH/2+WEOFFSETX), float64(SCREENHEIGHT/2+WEOFFSETY))
		opWea.GeoM.Scale(0.7, 0.7)
		op.Filter = ebiten.FilterLinear
		//skill option
		opSkill = &ebiten.DrawImageOptions{}
		opSkill.GeoM.Translate(float64(SCREENWIDTH/2+SKILLOFFSETX), float64(SCREENHEIGHT/2+SKILLOFFSETY))
		opSkill.CompositeMode = ebiten.CompositeModeLighter
		opSkill.GeoM.Scale(0.7, 0.7)
		opSkill.Filter = ebiten.FilterLinear
		//UI
		opUI = &ebiten.DrawImageOptions{}
		opUI.Filter = ebiten.FilterLinear
		opUI.GeoM.Translate(583, 380)
		wg.Done()
	}()
	wg.Wait()
	go func() {
		runtime.GC()
	}()
}

func (g *Game) Update() error {
	g.count++
	if g.player.state != ATTACK {
		g.player.state = IDLE
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		g.player.MouseX = x
		g.player.MouseY = y
		flg = true
	}
	//attack
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > 583 && x < 627 && y > 380 && y < 424 {
			g.player.state = ATTACK
			flg = false
			tools.MapCopy(asset, imgNow, g.player.state, g.player.direction, "man")
			tools.MapCopy(assetWea, imgWeaNow, g.player.state, g.player.direction, "weapon")
			tools.MapCopy(skill_asset, imgSkillNow, g.player.state, g.player.direction, "skill")
		}
	}
	//Calculate direction
	dir := tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY))
	//keyboard controll
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || dir == 6 {
		g.player.state = RUN
		g.player.direction = 6
		opBg.GeoM.Translate(SPEED, 0)
		g.player.x -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || dir == 2 {
		g.player.state = RUN
		g.player.direction = 2
		opBg.GeoM.Translate(-SPEED, 0)
		g.player.x += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || dir == 4 {
		g.player.state = RUN
		g.player.direction = 4
		opBg.GeoM.Translate(0, -SPEED)
		g.player.y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || dir == 0 {
		g.player.state = RUN
		g.player.direction = 0
		opBg.GeoM.Translate(0, SPEED)
		g.player.y -= 2
	}
	//mouse controll
	if dir == 1 && flg {
		g.player.state = RUN
		g.player.direction = 1
		opBg.GeoM.Translate(-SPEED, SPEED)
		g.player.y -= 2
		g.player.x += 2
		flg = false
	}

	if dir == 3 && flg {
		g.player.direction = 3
		g.player.state = RUN
		opBg.GeoM.Translate(-SPEED, -SPEED)
		g.player.y += 2
		g.player.x += 2
		flg = false
	}
	if dir == 5 && flg {
		g.player.direction = 5
		g.player.state = RUN
		opBg.GeoM.Translate(SPEED, -SPEED)
		g.player.y += 2
		g.player.x -= 2
		flg = false
	}
	if dir == 7 && flg {
		g.player.direction = 7
		g.player.state = RUN
		opBg.GeoM.Translate(SPEED, SPEED)
		g.player.y -= 2
		g.player.x -= 2
		flg = false
	}
	//states
	if g.player.state == IDLE {
		frameNums = 6

	} else if g.player.state == ATTACK {
		frameNums = 6
	} else {
		frameNums = 8
	}
	tools.MapCopy(asset, imgNow, g.player.state, g.player.direction, "man")
	tools.MapCopy(assetWea, imgWeaNow, g.player.state, g.player.direction, "weapon")
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error:", r)
		}
	}()
	//draw bg
	screen.DrawImage(bgImage, opBg)
	//draw UI
	screen.DrawImage(UI, opUI)
	//draw player
	screen.DrawImage(imgNow[counts], op)
	//draw wea
	screen.DrawImage(imgWeaNow[counts], opWea)
	//draw skill
	if g.player.state == ATTACK {
		screen.DrawImage(imgSkillNow[counts], opSkill)
	}
	//draw info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nplayer position %d,%d\nmouse position %d,%d\ndir %d",
		int64(ebiten.CurrentFPS()), int64(g.player.x), int64(g.player.y), g.player.MouseX, g.player.MouseY, tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY))))
	//change frame
	if g.count > 5 {
		counts++
		g.count = 0
		if counts >= frameNums {
			counts = 0
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return LAYOUTX, LAYOUTY
}

func main() {
	ebiten.SetWindowSize(SCREENWIDTH*2, SCREENHEIGHT*2)
	ebiten.SetWindowTitle("golang game test")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
