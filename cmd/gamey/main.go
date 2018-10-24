package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	"image/color"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

// https://coolors.co/ffba49-ed254e-75b9be-a8c256-fffbbd
var (
	colorOrange = &color.RGBA{255, 186, 73, 255}
	colorRed    = &color.RGBA{237, 37, 78, 255}
	colorBlue   = &color.RGBA{117, 185, 190, 255}
	colorGreen  = &color.RGBA{168, 194, 86, 255}
	colorYellow = &color.RGBA{255, 251, 189, 255}
)

const (
	chickenMaxSpeed     = 3.0
	chickenAcceleration = 0.5
)

type chicken struct {
	frames      []pixel.Rect
	spritesheet pixel.Picture
	spriteIndex int
	x           float64
	y           float64
	speed       float64
}

type game struct {
	win          *pixelgl.Window
	txt          *text.Text
	scoreDisplay *text.Text
	intro        bool
	introPos     int
	introText    []*gameText

	menu     bool
	menuPos  int
	menuText *gameText

	inPlay        bool
	score         int
	background    *color.RGBA
	chickens      []*chicken
	spritesheet   pixel.Picture
	chickenFrames []pixel.Rect
}

type gameText struct {
	text       string
	color      *color.RGBA
	background *color.RGBA
}

func (g *game) introNext() *gameText {
	if g.introPos < 0 || g.introPos > len(g.introText) {
		return nil
	}
	return g.introText[g.introPos]
}

func (g *game) menuNext() *gameText {
	return g.menuText
}

func (c *chicken) running() *pixel.Sprite {
	steps := []int{2, 5, 8, 11, 14}
	pixel := pixel.NewSprite(c.spritesheet, c.frames[steps[c.spriteIndex]])
	c.spriteIndex++
	if c.spriteIndex >= len(steps) {
		c.spriteIndex = 0
	}

	return pixel
}

func (g *game) spawnChicken() {
	// get random Y axis position

	y := g.win.Bounds().Center().Y * rand.Float64()
	chick := &chicken{
		frames:      g.chickenFrames,
		spritesheet: g.spritesheet,
		spriteIndex: 0,
		y:           y,
		x:           0,
	}
	g.chickens = append(g.chickens, chick)
	return
}

func (g *game) moveChickens() {
	// firstly sort out any scoring chickens
	for i, c := range g.chickens {
		if c.x >= g.win.Bounds().Max.X {
			g.score += 10
			g.chickens[i] = nil
		}
	}

	// now rebuild the chicken slice without the nils
	filtered := make([]*chicken, 0)
	for _, chicken := range g.chickens {
		if chicken != nil {
			filtered = append(filtered, chicken)
		}
	}

	g.chickens = filtered

	// finally actually move them
	for _, c := range g.chickens {
		c.speed = math.Min(c.speed+chickenAcceleration, chickenMaxSpeed)
		c.x += c.speed
	}
	return
}

func (g *game) drawChickens() {
	for _, c := range g.chickens {
		// draw the chicken
		c.running().Draw(g.win, pixel.IM.Scaled(pixel.ZV, 3).Moved(pixel.V(c.x, c.y*2)))
	}
	return
}

func (g *game) showScore() {
	g.scoreDisplay.Clear()
	fmt.Fprintf(g.scoreDisplay, "Score: %d", g.score)
	g.scoreDisplay.Draw(g.win, pixel.IM.Scaled(g.scoreDisplay.Orig, 2))
	return
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	game := &game{
		intro: true,
		introText: []*gameText{
			&gameText{text: "Chikkin", color: colorRed, background: colorBlue},
			&gameText{text: "A game by Peter Mellett", color: colorYellow, background: colorRed},
			&gameText{text: "For G", color: colorBlue, background: colorYellow},
		},
		menuText: &gameText{text: "Press the big red button", color: colorRed, background: colorGreen},
	}
	cfg := pixelgl.WindowConfig{
		Title:     "Chikkin, for George",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	game.win = win

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	game.txt = text.New(pixel.V(50, 500), atlas)
	game.scoreDisplay = text.New(pixel.V(10, 10), atlas)

	game.loadChickenSprites()

	fps := time.Tick(time.Second / 15)
	for !game.win.Closed() {
		if game.intro {
			runIntro(game)
		}

		if game.menu {
			runMenu(game)
		}

		if game.inPlay {
			// if the spacebar is held down do stuff!!!
			win.Clear(colorGreen)

			if game.win.Pressed(pixelgl.KeySpace) {
				game.spawnChicken()
			}
			game.moveChickens()
			game.drawChickens()
			game.showScore()
		}

		win.Update()
		<-fps
	}
}

func runMenu(game *game) {
	text := game.menuText
	game.txt.Clear()
	game.win.Clear(text.background)
	game.txt.Color = text.color
	fmt.Fprintf(game.txt, text.text)
	game.txt.Draw(game.win, pixel.IM.Scaled(game.txt.Orig, 4))

	if game.win.JustPressed(pixelgl.KeySpace) {
		game.menu = false
		game.inPlay = true
	}
}

func (g *game) loadChickenSprites() {
	spritesheet, err := loadPicture("resources/chick_24x24.png")
	if err != nil {
		panic(err)
	}

	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 24 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 24 {
			g.chickenFrames = append(g.chickenFrames, pixel.R(x, y, x+24, y+24))
		}
	}

	g.spritesheet = spritesheet
}

func runIntro(game *game) {
	text := game.introNext()

	game.txt.Clear()
	if text != nil {
		game.txt.Color = text.color
		game.win.Clear(text.background)
		fmt.Fprintf(game.txt, text.text)
	}
	game.txt.Draw(game.win, pixel.IM.Scaled(game.txt.Orig, 4))
	if game.win.JustPressed(pixelgl.KeySpace) {
		if game.introPos == len(game.introText)-1 {
			game.intro = false
			game.menu = true
			return
		}
		game.introPos++
	}
}

func main() {
	pixelgl.Run(run)
}
