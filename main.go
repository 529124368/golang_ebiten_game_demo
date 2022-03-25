package main

import (
	"embed"
	"fmt"
	"game/tools"
	_ "image/png"
	"log"
	"strconv"

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
	PATH          string  = "resource/playerAnm"
	SPEED         float64 = 2
	LAYOUTX       int     = 720
	LAYOUTY       int     = 480
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
	counts    int  = 0
	dir       int  = 0
	frameNums int  = 4
	flg       bool = false
	imgNow    map[int]*ebiten.Image
)
var op, opBg, opUI *ebiten.DrawImageOptions

var asset map[string]map[int]map[int]*ebiten.Image
var sub_asset_1_idle, sub_asset_1_attack, sub_asset_1_run map[int]map[int]*ebiten.Image
var sub_asset_2_idle, sub_asset_2_attack, sub_asset_2_run map[int]*ebiten.Image

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
	//get images
	imgNow = make(map[int]*ebiten.Image)

	//todo
	asset = make(map[string]map[int]map[int]*ebiten.Image)
	sub_asset_1_idle = make(map[int]map[int]*ebiten.Image)
	sub_asset_1_run = make(map[int]map[int]*ebiten.Image)
	sub_asset_1_attack = make(map[int]map[int]*ebiten.Image)

	//UI load
	s, _ := images.ReadFile("resource/UI/attack.png")
	mgUI := tools.GetEbitenImage(s)
	UI = mgUI
	//load idle images
	for j := 0; j < 8; j++ {
		sub_asset_2_idle = make(map[int]*ebiten.Image)
		for i := 0; i < 4; i++ {
			s, _ := images.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_stand_" + strconv.Itoa(i) + ".png")
			mg := tools.GetEbitenImage(s)
			sub_asset_2_idle[i] = mg
		}
		sub_asset_1_idle[j] = sub_asset_2_idle
	}
	//idle image to assign
	asset["idle"] = sub_asset_1_idle
	for j := 0; j < 8; j++ {
		sub_asset_2_run = make(map[int]*ebiten.Image)
		sub_asset_2_attack = make(map[int]*ebiten.Image)
		for i := 0; i < 6; i++ {
			s, _ := images.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_run_" + strconv.Itoa(i) + ".png")
			mg := tools.GetEbitenImage(s)
			sub_asset_2_run[i] = mg
			//
			s, _ = images.ReadFile("resource/man/warrior/" + strconv.Itoa(j) + "_attack_" + strconv.Itoa(i) + ".png")
			mg = tools.GetEbitenImage(s)
			sub_asset_2_attack[i] = mg
		}
		sub_asset_1_run[j] = sub_asset_2_run
		sub_asset_1_attack[j] = sub_asset_2_attack
	}
	//run attack image to assign
	asset["run"] = sub_asset_1_run
	asset["attack"] = sub_asset_1_attack
	//BG
	s, _ = images.ReadFile("resource/bg/bg1.png")
	img := tools.GetEbitenImage(s)
	bgImage = img
	opBg = &ebiten.DrawImageOptions{}
	opBg.Filter = ebiten.FilterLinear
	opBg.GeoM.Translate(-700, -550)
	//player option
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(SCREENWIDTH/2+OFFSETX), float64(SCREENHEIGHT/2+OFFSETY))
	op.GeoM.Scale(0.7, 0.7)
	op.Filter = ebiten.FilterLinear
	//UI
	opUI = &ebiten.DrawImageOptions{}
	opUI.Filter = ebiten.FilterLinear
	opUI.GeoM.Translate(583, 380)
	//copy
	MapCopy(asset["idle"][3], imgNow)
}

func MapCopy(a, b map[int]*ebiten.Image) {
	for k, v := range a {
		b[k] = v
	}
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
			MapCopy(asset["attack"][g.player.direction], imgNow)
		}
	}
	//Calculate direction
	dir = tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY))
	//keyboard controll
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || dir == 6 {
		g.player.state = RUN
		g.player.direction = 6
		opBg.GeoM.Translate(SPEED, 0)
		g.player.x -= 2
		MapCopy(asset["run"][6], imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || dir == 2 {
		g.player.state = RUN
		g.player.direction = 2
		opBg.GeoM.Translate(-SPEED, 0)
		g.player.x += 2
		MapCopy(asset["run"][2], imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || dir == 4 {
		g.player.state = RUN
		g.player.direction = 4
		opBg.GeoM.Translate(0, -SPEED)
		g.player.y += 2
		MapCopy(asset["run"][4], imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || dir == 0 {
		g.player.state = RUN
		g.player.direction = 0
		opBg.GeoM.Translate(0, SPEED)
		g.player.y -= 2
		MapCopy(asset["run"][0], imgNow)
	}
	//mouse controll
	if dir == 1 && flg {
		g.player.state = RUN
		g.player.direction = 1
		opBg.GeoM.Translate(-SPEED, SPEED)
		g.player.y -= 2
		g.player.x += 2
		MapCopy(asset["run"][1], imgNow)
		flg = false
	}
	if dir == 3 && flg {
		g.player.direction = 3
		g.player.state = RUN
		opBg.GeoM.Translate(-SPEED, -SPEED)
		g.player.y += 2
		g.player.x += 2
		MapCopy(asset["run"][3], imgNow)
		flg = false
	}
	if dir == 5 && flg {
		g.player.direction = 5
		g.player.state = RUN
		opBg.GeoM.Translate(SPEED, -SPEED)
		g.player.y += 2
		g.player.x -= 2
		MapCopy(asset["run"][5], imgNow)
		flg = false
	}
	if dir == 7 && flg {
		g.player.direction = 7
		g.player.state = RUN
		opBg.GeoM.Translate(SPEED, SPEED)
		g.player.y -= 2
		g.player.x -= 2
		MapCopy(asset["run"][7], imgNow)
		flg = false
	}
	//states
	if g.player.state == IDLE {
		frameNums = 4
		MapCopy(asset["idle"][g.player.direction], imgNow)

	} else {
		frameNums = 6
	}
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
	//draw palyer
	screen.DrawImage(imgNow[counts], op)
	//draw info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nplayer position %d,%d\nmouse position %d,%d\ndir %d",
		int64(ebiten.CurrentFPS()), int64(g.player.x), int64(g.player.y), g.player.MouseX, g.player.MouseY, tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY))))
	//change frame
	if g.count > frameNums {
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
