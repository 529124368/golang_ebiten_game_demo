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
var imgs_idel0, imgs_idel0_5, imgs_idel1, imgs_idel1_5, imgs_idel2, imgs_idel2_5, imgs_idel3, imgs_idel3_5 map[int]*ebiten.Image
var imgs_run0, imgs_run0_5, imgs_run1, imgs_run1_5, imgs_run2, imgs_run2_5, imgs_run3, imgs_run3_5 map[int]*ebiten.Image
var imgs_atc0, imgs_atc0_5, imgs_atc1, imgs_atc1_5, imgs_atc2, imgs_atc2_5, imgs_atc3, imgs_atc3_5 map[int]*ebiten.Image
var op, opBg, opUI *ebiten.DrawImageOptions

//BG UI
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
	//idle
	imgs_idel0 = make(map[int]*ebiten.Image)
	imgs_idel0_5 = make(map[int]*ebiten.Image)
	imgs_idel1 = make(map[int]*ebiten.Image)
	imgs_idel1_5 = make(map[int]*ebiten.Image)
	imgs_idel2 = make(map[int]*ebiten.Image)
	imgs_idel2_5 = make(map[int]*ebiten.Image)
	imgs_idel3 = make(map[int]*ebiten.Image)
	imgs_idel3_5 = make(map[int]*ebiten.Image)
	//run
	imgs_run0 = make(map[int]*ebiten.Image)
	imgs_run0_5 = make(map[int]*ebiten.Image)
	imgs_run1 = make(map[int]*ebiten.Image)
	imgs_run1_5 = make(map[int]*ebiten.Image)
	imgs_run2 = make(map[int]*ebiten.Image)
	imgs_run2_5 = make(map[int]*ebiten.Image)
	imgs_run3 = make(map[int]*ebiten.Image)
	imgs_run3_5 = make(map[int]*ebiten.Image)
	//attack
	imgs_atc0 = make(map[int]*ebiten.Image)
	imgs_atc0_5 = make(map[int]*ebiten.Image)
	imgs_atc1 = make(map[int]*ebiten.Image)
	imgs_atc1_5 = make(map[int]*ebiten.Image)
	imgs_atc2 = make(map[int]*ebiten.Image)
	imgs_atc2_5 = make(map[int]*ebiten.Image)
	imgs_atc3 = make(map[int]*ebiten.Image)
	imgs_atc3_5 = make(map[int]*ebiten.Image)

	//UI load
	s, _ := images.ReadFile("resource/UI/attack.png")
	mgUI := tools.GetEbitenImage(s)
	UI = mgUI
	//load images
	for i := 0; i < 4; i++ {
		s, _ := images.ReadFile(PATH + "/idle0/" + strconv.Itoa(i) + ".png")
		mg0 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/idle0.5/" + strconv.Itoa(i) + ".png")
		mg0_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/idle1/" + strconv.Itoa(i) + ".png")
		mg1 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/idle1.5/" + strconv.Itoa(i) + ".png")
		mg1_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/idle2/" + strconv.Itoa(i) + ".png")
		mg2 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/idle2.5/" + strconv.Itoa(i) + ".png")
		mg2_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/idle3/" + strconv.Itoa(i) + ".png")
		mg3 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/idle3.5/" + strconv.Itoa(i) + ".png")
		mg3_5 := tools.GetEbitenImage(s)
		imgs_idel0[i] = mg0
		imgs_idel0_5[i] = mg0_5
		imgs_idel1[i] = mg1
		imgs_idel1_5[i] = mg1_5
		imgs_idel2[i] = mg2
		imgs_idel2_5[i] = mg2_5
		imgs_idel3[i] = mg3
		imgs_idel3_5[i] = mg3_5
	}
	for i := 0; i < 6; i++ {
		s, _ := images.ReadFile(PATH + "/run0/" + strconv.Itoa(i) + ".png")
		img0 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/run0.5/" + strconv.Itoa(i) + ".png")
		img0_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/run1/" + strconv.Itoa(i) + ".png")
		img1 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/run1.5/" + strconv.Itoa(i) + ".png")
		img1_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/run2/" + strconv.Itoa(i) + ".png")
		img2 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/run2.5/" + strconv.Itoa(i) + ".png")
		img2_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/run3/" + strconv.Itoa(i) + ".png")
		img3 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/run3.5/" + strconv.Itoa(i) + ".png")
		img3_5 := tools.GetEbitenImage(s)
		//
		s, _ = images.ReadFile(PATH + "/attack0/" + strconv.Itoa(i) + ".png")
		imgs0 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/attack0.5/" + strconv.Itoa(i) + ".png")
		imgs0_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/attack1/" + strconv.Itoa(i) + ".png")
		imgs1 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/attack1.5/" + strconv.Itoa(i) + ".png")
		imgs1_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/attack2/" + strconv.Itoa(i) + ".png")
		imgs2 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/attack2.5/" + strconv.Itoa(i) + ".png")
		imgs2_5 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/attack3/" + strconv.Itoa(i) + ".png")
		imgs3 := tools.GetEbitenImage(s)
		s, _ = images.ReadFile(PATH + "/attack3.5/" + strconv.Itoa(i) + ".png")
		imgs3_5 := tools.GetEbitenImage(s)
		imgs_run0[i] = img0
		imgs_run0_5[i] = img0_5
		imgs_run1[i] = img1
		imgs_run1_5[i] = img1_5
		imgs_run2[i] = img2
		imgs_run2_5[i] = img2_5
		imgs_run3[i] = img3
		imgs_run3_5[i] = img3_5
		//
		imgs_atc0[i] = imgs0
		imgs_atc0_5[i] = imgs0_5
		imgs_atc1[i] = imgs1
		imgs_atc1_5[i] = imgs1_5
		imgs_atc2[i] = imgs2
		imgs_atc2_5[i] = imgs2_5
		imgs_atc3[i] = imgs3
		imgs_atc3_5[i] = imgs3_5
	}
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
	MapCopy(imgs_run0, imgNow)
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
			switch g.player.direction {
			case 0:
				MapCopy(imgs_atc0, imgNow)
			case 1:
				MapCopy(imgs_atc0_5, imgNow)
			case 2:
				MapCopy(imgs_atc1, imgNow)
			case 3:
				MapCopy(imgs_atc1_5, imgNow)
			case 4:
				MapCopy(imgs_atc2, imgNow)
			case 5:
				MapCopy(imgs_atc2_5, imgNow)
			case 6:
				MapCopy(imgs_atc3, imgNow)
			case 7:
				MapCopy(imgs_atc3_5, imgNow)
			}

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
		MapCopy(imgs_run3, imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || dir == 2 {
		g.player.state = RUN
		g.player.direction = 2
		opBg.GeoM.Translate(-SPEED, 0)
		g.player.x += 2
		MapCopy(imgs_run1, imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || dir == 4 {
		g.player.state = RUN
		g.player.direction = 4
		opBg.GeoM.Translate(0, -SPEED)
		g.player.y += 2
		MapCopy(imgs_run2, imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || dir == 0 {
		g.player.state = RUN
		g.player.direction = 0
		opBg.GeoM.Translate(0, SPEED)
		g.player.y -= 2
		MapCopy(imgs_run0, imgNow)
	}
	//mouse controll
	if dir == 1 && flg {
		g.player.state = RUN
		g.player.direction = 1
		opBg.GeoM.Translate(-SPEED, SPEED)
		g.player.y -= 2
		g.player.x += 2
		MapCopy(imgs_run0_5, imgNow)
		flg = false
	}
	if dir == 3 && flg {
		g.player.direction = 3
		g.player.state = RUN
		opBg.GeoM.Translate(-SPEED, -SPEED)
		g.player.y += 2
		g.player.x += 2
		MapCopy(imgs_run1_5, imgNow)
		flg = false
	}
	if dir == 5 && flg {
		g.player.direction = 5
		g.player.state = RUN
		opBg.GeoM.Translate(SPEED, -SPEED)
		g.player.y += 2
		g.player.x -= 2
		MapCopy(imgs_run2_5, imgNow)
		flg = false
	}
	if dir == 7 && flg {
		g.player.direction = 7
		g.player.state = RUN
		opBg.GeoM.Translate(SPEED, SPEED)
		g.player.y -= 2
		g.player.x -= 2
		MapCopy(imgs_run3_5, imgNow)
		flg = false
	}
	//states
	if g.player.state == IDLE {
		frameNums = 4
		if g.player.direction == 0 {
			MapCopy(imgs_idel0, imgNow)
		}
		if g.player.direction == 1 {
			MapCopy(imgs_idel0_5, imgNow)
		}
		if g.player.direction == 2 {
			MapCopy(imgs_idel1, imgNow)
		}
		if g.player.direction == 3 {
			MapCopy(imgs_idel1_5, imgNow)
		}
		if g.player.direction == 4 {
			MapCopy(imgs_idel2, imgNow)
		}
		if g.player.direction == 5 {
			MapCopy(imgs_idel2_5, imgNow)
		}
		if g.player.direction == 6 {
			MapCopy(imgs_idel3, imgNow)
		}
		if g.player.direction == 7 {
			MapCopy(imgs_idel3_5, imgNow)
		}

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
	//draw images
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
