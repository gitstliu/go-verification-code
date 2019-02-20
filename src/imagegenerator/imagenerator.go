package imagegenerator

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"

	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"syscommon"

	"github.com/gitstliu/log4go"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Position struct {
	X int
	Y int
}

type NoiceRunePosition struct {
	KeyWord        string
	ConfusionWords []rune
	Position       []Position
	WordIndex      []int
	ColumnCount    int
	LineCount      int
}

var currFont *truetype.Font

func LoadFont(fontfile string) error {
	file, openFileErr := os.Open(fontfile)
	//	fontBytes, err := ioutil.ReadFile(fontfile)
	if openFileErr != nil {
		log4go.Error(openFileErr)
		return openFileErr
	}

	fontBytes, fileReadAllErr := ioutil.ReadAll(file)

	if fileReadAllErr != nil {
		log4go.Error(fileReadAllErr)
		return fileReadAllErr
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log4go.Error(err)
		return err
	}

	currFont = f

	return nil
}

func CreateCodeImage(height int, width int, dpi float64, fontSize float64, hinting string, spacing float64, text []string) (*image.RGBA, error) {
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	//	if wonb {
	//		fg, bg = image.White, image.Black
	//		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	//	}
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(currFont)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(fontSize)>>6))
	for _, s := range text {
		_, err := c.DrawString(s, pt)
		if err != nil {
			log4go.Error(err)
			return nil, err
		}
		pt.Y += c.PointToFixed(fontSize * spacing)
	}

	return rgba, nil
}

func CreateCodeImageByPosition(xCount int, yCount int, stepCell int, dpi float64, fontSize float64, hinting string, spacing float64, position *NoiceRunePosition, fgValue string, bgValue string) (*image.RGBA, error) {
	log4go.Debug("fgValue = ", fgValue)
	log4go.Debug("bgValue = ", bgValue)

	fgTemp, fgTempErr := strconv.ParseUint(fgValue, 16, 32)
	bgTemp, bgTempErr := strconv.ParseUint(bgValue, 16, 32)

	//	fgTemp, fgTempErr := strconv.ParseUint("FFFFFFFF", 16, 32)
	//	bgTemp, bgTempErr := strconv.ParseUint("000000FF", 16, 32)

	log4go.Debug("fgTemp = ", fgTemp)
	log4go.Debug("bgTemp = ", bgTemp)

	if fgTempErr != nil {
		log4go.Debug(fgTempErr)
		return nil, errors.New("Parameter fg format invalid!! The value must look like FFFFFFFF ")
	}

	if bgTempErr != nil {
		log4go.Debug(bgTempErr)
		return nil, errors.New("Parameter bg format invalid!! The value must look like FFFFFFFF ")
	}

	fgByteMeta := syscommon.UInt32ToBytes(uint32(fgTemp))
	bgByteMeta := syscommon.UInt32ToBytes(uint32(bgTemp))

	log4go.Debug("fgByteMeta = ", fgByteMeta)
	log4go.Debug("bgByteMeta = ", bgByteMeta)

	//	fg, bg := image.Black, image.Transparent
	fg := image.NewUniform(color.RGBA{uint8(fgByteMeta[0]), uint8(fgByteMeta[1]), uint8(fgByteMeta[2]), uint8(fgByteMeta[3])})
	bg := image.NewUniform(color.RGBA{uint8(bgByteMeta[0]), uint8(bgByteMeta[1]), uint8(bgByteMeta[2]), uint8(bgByteMeta[3])})
	//	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	//	if wonb {
	//		fg, bg = image.White, image.Black
	//		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	//	}
	//	bg.RGBA()
	rgba := image.NewRGBA(image.Rect(0, 0, xCount*stepCell, yCount*stepCell))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(currFont)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	//	for i := 0; i < 200; i++ {
	//		rgba.Set(10, 10+i, ruler)
	//		rgba.Set(10+i, 10, ruler)
	//	}

	// Draw the text.

	log4go.Debug("Positions = %v", position.Position)
	log4go.Debug("len(position.ConfusionWords) = %v", len(position.ConfusionWords))

	for index, pos := range position.Position {
		offset := rand.Intn(stepCell / 2)
		currPt := freetype.Pt(pos.Y*stepCell+offset, pos.X*stepCell+offset+int(c.PointToFixed(fontSize)>>6))
		_, err := c.DrawString(string(position.ConfusionWords[index]), currPt)
		if err != nil {
			log4go.Error(err)
			return nil, err
		}
		//		pt.Y += c.PointToFixed(fontSize * spacing)
	}
	log4go.Debug("ReturnRGBA")
	return rgba, nil
}

func ImageToPngBytes(rgba *image.RGBA) ([]byte, error) {
	buf := &bytes.Buffer{}

	err := png.Encode(buf, rgba)

	if err != nil {
		return nil, err
	}

	buffer := buf.Bytes()
	log4go.Debug("len(buffer) = %v", len(buffer))

	return buffer, nil
}
