package main

import (
	"embed"
	"fmt"
	"game/role"
	"game/tools"
	_ "image/png"
	"log"
	"runtime"

	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//config
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

var game *Game

//Game
type Game struct {
	count  int
	player *role.Player
	//monster *role.Monster
}

var (
	counts    int  = 0
	frameNums int  = 4
	flg       bool = false
)
var op, opWea, opBg, opUI, opSkill *ebiten.DrawImageOptions

var bgImage, UI *ebiten.Image

//go:embed resource
var images embed.FS

//factory
func NewGame() *Game {
	r := role.NewPlayer(float64(SCREENWIDTH/2), float64(SCREENHEIGHT/2), IDLE, 0, 0, 0, &images)
	gameStart := &Game{
		count:  0,
		player: r,
	}
	return gameStart
}

func init() {
	//game init
	game = NewGame()
	//palyer init
	go func() {
		game.player.LoadImages()
		runtime.GC()
	}()
	//UI load
	s, _ := images.ReadFile("resource/UI/attack.png")
	mgUI := tools.GetEbitenImage(s)
	UI = mgUI
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
		runtime.GC()
	}()
}

func (g *Game) Update() error {
	g.count++
	if g.player.State != ATTACK {
		g.player.State = IDLE
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
			if g.player.State != ATTACK {
				counts = 0
			}
			g.player.State = ATTACK
			flg = false
			g.player.SetAnimator(0)
		}
	}
	//Calculate direction
	dir := tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY))
	//keyboard controll
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || dir == 6 {
		g.player.SetPlayerState(RUN, 6)
		opBg.GeoM.Translate(SPEED, 0)
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || dir == 2 {
		g.player.SetPlayerState(RUN, 2)
		opBg.GeoM.Translate(-SPEED, 0)
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || dir == 4 {
		g.player.SetPlayerState(RUN, 4)
		opBg.GeoM.Translate(0, -SPEED)
		g.player.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || dir == 0 {
		g.player.SetPlayerState(RUN, 0)
		opBg.GeoM.Translate(0, SPEED)
		g.player.Y -= 2
	}
	//mouse controll
	if dir == 1 && flg {
		g.player.SetPlayerState(RUN, 1)
		opBg.GeoM.Translate(-SPEED, SPEED)
		g.player.Y -= 2
		g.player.X += 2
		flg = false
	}

	if dir == 3 && flg {
		g.player.SetPlayerState(RUN, 3)
		opBg.GeoM.Translate(-SPEED, -SPEED)
		g.player.Y += 2
		g.player.X += 2
		flg = false
	}
	if dir == 5 && flg {
		g.player.SetPlayerState(RUN, 5)
		opBg.GeoM.Translate(SPEED, -SPEED)
		g.player.Y += 2
		g.player.X -= 2
		flg = false
	}
	if dir == 7 && flg {
		g.player.SetPlayerState(RUN, 7)
		opBg.GeoM.Translate(SPEED, SPEED)
		g.player.Y -= 2
		g.player.X -= 2
		flg = false
	}
	//states
	if g.player.State == IDLE {
		frameNums = 6

	} else if g.player.State == ATTACK {
		frameNums = 6
	} else {
		frameNums = 8
	}
	g.player.SetAnimator(1)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("has error is :", r)
		}
	}()
	//draw bg
	screen.DrawImage(bgImage, opBg)
	//draw UI
	screen.DrawImage(UI, opUI)
	//draw player
	screen.DrawImage(g.player.ImgNow[counts], op)
	//draw wea
	screen.DrawImage(g.player.ImgWeaNow[counts], opWea)
	//draw skill
	if g.player.State == ATTACK {
		screen.DrawImage(g.player.ImgSkillNow[counts], opSkill)
	}
	//draw info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %d\nplayer position %d,%d\nmouse position %d,%d\ndir %d",
		int64(ebiten.CurrentFPS()), int64(g.player.X), int64(g.player.Y), g.player.MouseX, g.player.MouseY, tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY))))
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
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
