package main

import (
	"embed"
	"fmt"
	"game/role"
	"game/tools"
	_ "image/png"
	"log"
	"runtime"
	"strconv"

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
var op, opWea, opBg, opUI, opUI1, opSkill *ebiten.DrawImageOptions

var bgImage, UI, UI1 *ebiten.Image

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
	//
	go func() {
		//UI load
		s, _ := images.ReadFile("resource/UI/liehuo.png")
		mgUI := tools.GetEbitenImage(s)
		UI = mgUI
		runtime.GC()
	}()

	go func() {
		//UI load
		s, _ := images.ReadFile("resource/UI/chisha.png")
		mgUI := tools.GetEbitenImage(s)
		UI1 = mgUI
		runtime.GC()
	}()

	go func() {
		//BG
		s2, _ := images.ReadFile("resource/bg/bg1.png")
		img := tools.GetEbitenImage(s2)
		bgImage = img
		opBg = &ebiten.DrawImageOptions{}
		opBg.Filter = ebiten.FilterLinear
		opBg.GeoM.Translate(-700, -550)
		//UI
		opUI = &ebiten.DrawImageOptions{}
		opUI.Filter = ebiten.FilterLinear
		opUI.GeoM.Translate(583, 380)
		opUI1 = &ebiten.DrawImageOptions{}
		opUI1.Filter = ebiten.FilterLinear
		opUI1.GeoM.Translate(620, 330)
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
	//Calculate direction
	dir := tools.CaluteDir(PLAYERCENTERX, PLAYERCENTERY, int64(g.player.MouseX), int64(g.player.MouseY))
	//attack
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		//liehuo
		if x > 583 && x < 627 && y > 380 && y < 424 {
			g.player.SkillName = "liehuo"
			if g.player.State != ATTACK {
				counts = 0
			}
			flg = false
			if g.player.Direction != dir || g.player.State != ATTACK {
				g.player.SetPlayerState(ATTACK, dir)

			}
		}
		if x > 621 && x < 664 && y > 331 && y < 373 {
			g.player.SkillName = "chisha"
			if g.player.State != ATTACK {
				counts = 0
			}
			flg = false
			if g.player.Direction != dir || g.player.State != ATTACK {
				g.player.SetPlayerState(ATTACK, dir)

			}
		}
	}
	//keyboard controll
	// if dir == 6 {
	// 	if g.player.Direction != dir || g.player.State != RUN {
	// 		g.player.SetPlayerState(RUN, dir)
	// 		g.player.SetAnimator(1)
	// 	}
	// 	opBg.GeoM.Translate(SPEED, 0)
	// 	g.player.X -= 2
	// }
	// if dir == 2 {
	// 	if g.player.Direction != dir || g.player.State != RUN {
	// 		g.player.SetPlayerState(RUN, dir)
	// 		g.player.SetAnimator(1)
	// 	}
	// 	opBg.GeoM.Translate(-SPEED, 0)
	// 	g.player.X += 2
	// }
	// if dir == 4 {
	// 	if g.player.Direction != dir || g.player.State != RUN {
	// 		g.player.SetPlayerState(RUN, dir)
	// 		g.player.SetAnimator(1)
	// 	}
	// 	opBg.GeoM.Translate(0, -SPEED)
	// 	g.player.Y += 2
	// }
	// if dir == 0 {
	// 	if g.player.Direction != dir || g.player.State != RUN {
	// 		g.player.SetPlayerState(RUN, dir)
	// 		g.player.SetAnimator(1)
	// 	}
	// 	opBg.GeoM.Translate(0, SPEED)
	// 	g.player.Y -= 2
	// }
	//mouse controll
	if dir == 1 && flg {
		if g.player.Direction != dir || g.player.State != RUN {
			g.player.SetPlayerState(RUN, dir)
		}
		opBg.GeoM.Translate(-SPEED, SPEED)
		g.player.Y -= 2
		g.player.X += 2
		flg = false
	}

	if dir == 3 && flg {
		if g.player.Direction != dir || g.player.State != RUN {
			g.player.SetPlayerState(RUN, dir)

		}
		opBg.GeoM.Translate(-SPEED, -SPEED)
		g.player.Y += 2
		g.player.X += 2
		flg = false
	}
	if dir == 5 && flg {
		if g.player.Direction != dir || g.player.State != RUN {
			g.player.SetPlayerState(RUN, dir)
		}
		opBg.GeoM.Translate(SPEED, -SPEED)
		g.player.Y += 2
		g.player.X -= 2
		flg = false
	}
	if dir == 7 && flg {
		if g.player.Direction != dir || g.player.State != RUN {
			g.player.SetPlayerState(RUN, dir)
		}
		opBg.GeoM.Translate(SPEED, SPEED)
		g.player.Y -= 2
		g.player.X -= 2
		flg = false
	}
	//states
	if g.player.State == IDLE {
		frameNums = 6
		g.player.SetPlayerState(IDLE, dir)

	} else if g.player.State == ATTACK {
		frameNums = 6
	} else {
		frameNums = 8
	}
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
	screen.DrawImage(UI1, opUI1)

	//
	name := ""
	nameSkill := ""
	switch g.player.State {
	case ATTACK:
		if counts >= 6 {
			counts = 0
		}
		name = strconv.Itoa(g.player.Direction) + "_attack_" + strconv.Itoa(counts) + ".png"
		nameSkill = strconv.Itoa(g.player.Direction) + "_skill_" + strconv.Itoa(counts) + ".png"
	case IDLE:
		if counts >= 6 {
			counts = 0
		}
		name = strconv.Itoa(g.player.Direction) + "_stand_" + strconv.Itoa(counts) + ".png"
	default:
		name = strconv.Itoa(g.player.Direction) + "_run_" + strconv.Itoa(counts) + ".png"
	}
	imagess, x, y := g.player.GetAnimator("man", name)
	//draw player
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(SCREENWIDTH/2+OFFSETX+x), float64(SCREENHEIGHT/2+OFFSETY+y))
	op.GeoM.Scale(0.7, 0.7)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(imagess, op)
	//draw wea
	imagess, x, y = g.player.GetAnimator("weapon", name)
	opWea = &ebiten.DrawImageOptions{}
	opWea.GeoM.Translate(float64(SCREENWIDTH/2+WEOFFSETX+x), float64(SCREENHEIGHT/2+WEOFFSETY+y))
	opWea.GeoM.Scale(0.7, 0.7)
	opWea.Filter = ebiten.FilterLinear
	screen.DrawImage(imagess, opWea)
	//draw skill
	if g.player.State == ATTACK {
		imagey, x, y := g.player.GetAnimator("skill", nameSkill)
		//skill option
		opSkill = &ebiten.DrawImageOptions{}
		opSkill.GeoM.Translate(float64(SCREENWIDTH/2+x), float64(SCREENHEIGHT/2+y))
		opSkill.CompositeMode = ebiten.CompositeModeLighter
		opSkill.GeoM.Scale(1.5, 1.5)
		opSkill.Filter = ebiten.FilterLinear
		screen.DrawImage(imagey, opSkill)
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
