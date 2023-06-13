package painter

import (
	"image"
	"image/color"
	"image/draw"

	//"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/shiny/screen"
)

type mockOperation struct {
	mock.Mock
}

func TestLoop_Post(t *testing.T) {
	screen := new(mockScreen)
	texture := new(mockTexture)
	receiver := new(testReceiver)
	tx := image.Pt(800, 800)
	loop := Loop{
		Receiver: receiver,
	}

	screen.On("NewTexture", tx).Return(texture, nil)
	receiver.On("Update", texture).Return()

	loop.Start(screen)

	operationOne := new(mockOperation)

	texture.On("Bounds").Return(image.Rectangle{})
	operationOne.On("Do", texture).Return(false)

	assert.Empty(t, loop.mq.ops)
	time.Sleep(1 * time.Second)
	assert.Empty(t, loop.mq.ops)

	operationOne.AssertCalled(t, "Do", texture)
	receiver.AssertCalled(t, "Update", texture)
	screen.AssertCalled(t, "NewTexture", image.Pt(800, 800))
}

func logOp(t *testing.T, msg string, op OperationFunc) OperationFunc {
	return func(tx screen.Texture) {
		t.Log(msg)
		op(tx)
	}
}

type testReceiver struct {
	mock.Mock
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.Called(t)
}

type mockScreen struct {
	mock.Mock
}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	args := m.Called(size)
	return args.Get(0).(screen.Texture), args.Error(1)
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

type mockTexture struct {
	mock.Mock
}

func (m *mockTexture) Release() {
	m.Called()
}

func (m *mockTexture) Size() image.Point {
	args := m.Called()
	return args.Get(0).(image.Point)
}

func (m *mockTexture) Bounds() image.Rectangle {
	args := m.Called()
	return args.Get(0).(image.Rectangle)
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	m.Called(dp, src, sr)
}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Called(dr, src, op)
}
