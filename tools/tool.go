package tools

import (
	"bytes"
	"image"
	"log"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

func CaluteDir(x, y, x_tar, y_tar int64) int {
	if x < x_tar && float64(y) == math.Abs(float64(y_tar)) {
		return 2
	}
	if float64(x) == math.Abs(float64(x_tar)) && y > y_tar {
		return 0
	}
	if x < x_tar && y > y_tar {
		return 1
	}
	if x < x_tar && y < y_tar {
		return 3
	}

	if float64(x) == math.Abs(float64(x_tar)) && y < y_tar {
		return 4
	}
	if x > x_tar && y < y_tar {
		return 5
	}
	if x > x_tar && float64(y) == math.Abs(float64(y_tar)) {
		return 6
	}
	if x > x_tar && y > y_tar {
		return 7
	}
	return 0
}

//
func GetEbitenImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func MapCopy(a map[string]*ebiten.Image, b map[int]*ebiten.Image, states, dir int, style string) {
	if style != "skill" {
		index_name := ""
		switch style {
		case "man":
			index_name = "m"
		case "weapon":
			index_name = "w"
		}
		frame := 6
		switch states {
		case 0:
			index_name += "s_"
		case 1:
			index_name += "r_"
			frame = 8
		case 2:
			index_name += "a_"
		}
		for i := 0; i < frame; i++ {
			name := index_name + strconv.Itoa(dir) + "_" + strconv.Itoa(i)
			b[i] = a[name]
		}
	} else {
		for i := 0; i < 6; i++ {
			name := strconv.Itoa(dir) + "_" + strconv.Itoa(i)
			b[i] = a[name]
		}

	}

}
