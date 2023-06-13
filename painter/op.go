package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

type BgRectangle struct {
	X1, Y1, X2, Y2 int
}

func (bg *BgRectangle) Do(t screen.Texture) bool {
	t.Fill(image.Rect(bg.X1, bg.Y1, bg.X2, bg.Y2), color.Black, screen.Src)
	return false
}

type Figure struct {
	X, Y int
}

func (f *Figure) Do(t screen.Texture) bool {
	height, width := 400, 150
	yellow := color.RGBA{255, 255, 0, 255}
	t.Fill(image.Rect(f.X-width/2, f.Y-height/2, f.X+width/2, f.Y+height/2), yellow, draw.Src)
	t.Fill(image.Rect(f.X-height/2, f.Y-width/2, f.X+height/2, f.Y+width/2), yellow, draw.Src)

	return false
}

type Move struct {
	X, Y    int
	Figures []*Figure
}

func (m *Move) Do(t screen.Texture) bool {
	for _, figure := range m.Figures {
		figure.X += m.X
		figure.Y += m.Y
	}
	return false
}

func ResetScreen(t screen.Texture) {
	t.Fill(t.Bounds(), color.Black, draw.Src)
}
