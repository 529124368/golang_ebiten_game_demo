package main

import (
	"fmt"

	_ "image/png"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	IDLE         int = 0
	RUN          int = 1
	ATTACK       int = 2
	screenWidth  int = 320
	screenHeight int = 240
	offsetX      int = 150
	offsetY      int = 80
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

var counts int = 0
var frameNums int = 4
var imgs_idel0, imgs_idel1, imgs_idel2, imgs_idel3, imgs_run0, imgs_run1, imgs_run2, imgs_run3, imgNow map[int]*ebiten.Image
var op, opBg *ebiten.DrawImageOptions

//BG
var bgImage *ebiten.Image

//factory
func NewGame() *Game {
	gameStart := &Game{
		count: 0,
		player: &Player{
			x:         float64(screenWidth/2 + offsetX),
			y:         float64(screenHeight/2 + offsetY),
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
	imgs_idel1 = make(map[int]*ebiten.Image)
	imgs_idel2 = make(map[int]*ebiten.Image)
	imgs_idel3 = make(map[int]*ebiten.Image)
	//run
	imgs_run0 = make(map[int]*ebiten.Image)
	imgs_run1 = make(map[int]*ebiten.Image)
	imgs_run2 = make(map[int]*ebiten.Image)
	imgs_run3 = make(map[int]*ebiten.Image)
	for i := 0; i < 4; i++ {
		mg, _, err := ebitenutil.NewImageFromFile("resource/idle0/" + strconv.Itoa(i) + ".png")
		mg1, _, err := ebitenutil.NewImageFromFile("resource/idle1/" + strconv.Itoa(i) + ".png")
		mg2, _, err := ebitenutil.NewImageFromFile("resource/idle2/" + strconv.Itoa(i) + ".png")
		mg3, _, err := ebitenutil.NewImageFromFile("resource/idle3/" + strconv.Itoa(i) + ".png")
		if err != nil {
			log.Fatal(err)
		}
		imgs_idel0[i] = mg
		imgs_idel1[i] = mg1
		imgs_idel2[i] = mg2
		imgs_idel3[i] = mg3
	}
	for i := 0; i < 6; i++ {
		img, _, err := ebitenutil.NewImageFromFile("resource/run0/" + strconv.Itoa(i) + ".png")
		img1, _, err := ebitenutil.NewImageFromFile("resource/run1/" + strconv.Itoa(i) + ".png")
		img2, _, err := ebitenutil.NewImageFromFile("resource/run2/" + strconv.Itoa(i) + ".png")
		img3, _, err := ebitenutil.NewImageFromFile("resource/run3/" + strconv.Itoa(i) + ".png")
		if err != nil {
			log.Fatal(err)
		}
		imgs_run0[i] = img
		imgs_run1[i] = img1
		imgs_run2[i] = img2
		imgs_run3[i] = img3
	}
	//BG
	img, _, err := ebitenutil.NewImageFromFile("resource/bg/bg1.png")
	if err != nil {
		log.Fatal(err)
	}
	bgImage = img
	opBg = &ebiten.DrawImageOptions{}
	opBg.Filter = ebiten.FilterLinear
	opBg.GeoM.Translate(-700, -550)
	opBg.GeoM.Scale(0.5, 0.5)
	//image option
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screenWidth/2+offsetX), float64(screenHeight/2+offsetY))
	fmt.Println(float64(screenWidth/2+offsetX), float64(screenHeight/2+offsetY))
	op.GeoM.Scale(0.4, 0.4)
	op.Filter = ebiten.FilterLinear
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
	g.player.state = IDLE
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.player.MouseX = x
		g.player.MouseY = y
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.state = RUN
		g.player.direction = 3
		opBg.GeoM.Translate(1, 0)
		g.player.x -= 1
		MapCopy(imgs_run3, imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.state = RUN
		g.player.direction = 1
		opBg.GeoM.Translate(-1, 0)
		g.player.x += 1
		MapCopy(imgs_run1, imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.state = RUN
		g.player.direction = 2
		opBg.GeoM.Translate(0, -1)
		g.player.y += 1
		MapCopy(imgs_run2, imgNow)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.state = RUN
		g.player.direction = 0
		opBg.GeoM.Translate(0, 1)
		g.player.y -= 1
		MapCopy(imgs_run0, imgNow)
	}
	if g.player.state == IDLE {
		frameNums = 4
		if g.player.direction == 0 {
			MapCopy(imgs_idel0, imgNow)
		}
		if g.player.direction == 1 {
			MapCopy(imgs_idel1, imgNow)
		}
		if g.player.direction == 2 {
			MapCopy(imgs_idel2, imgNow)
		}
		if g.player.direction == 3 {
			MapCopy(imgs_idel3, imgNow)
		}

	} else {
		frameNums = 6
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("#################", r)
		}
	}()
	//draw bg
	screen.DrawImage(bgImage, opBg)
	//draw images
	screen.DrawImage(imgNow[counts], op)
	//draw info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS:%d\nplayer positon %d,%d\nmouse position %d,%d", int64(ebiten.CurrentFPS()), int64(g.player.x), int64(g.player.y), g.player.MouseX, g.player.MouseY))
	if g.count > frameNums {
		counts++
		g.count = 0
		if counts >= frameNums {
			counts = 0
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	//ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("golang game test")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
