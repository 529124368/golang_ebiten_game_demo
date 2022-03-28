package tools

import (
	"bytes"
	"image"
	"log"
	"math"

	"github.com/fzipp/texturepacker"
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

//read images from bytes
func GetEbitenImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

//read images from plist
func GetImageFromPlist(s []byte, json []byte) (*texturepacker.SpriteSheet, *image.NRGBA) {
	sheet, err := texturepacker.SheetFromData(json, texturepacker.FormatJSONHash{})

	img, _, err := image.Decode(bytes.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	//sheetImage := imageToRGBA(img)
	sheetImage := img.(*image.NRGBA)
	return sheet, sheetImage
}
