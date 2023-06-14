package test

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"github.com/mikhmol/Architecture_Lab3/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post(t *testing.T) {
	screen := new(screenMock)
	texture := new(textureMock)
	receiver := new(receiverMock)
	tx := image.Pt(800, 800)
	loop := painter.Loop{
		Receiver: receiver,
	}

	screen.On("NewTexture", tx).Return(texture, nil)
	receiver.On("Update", texture).Return()

	loop.Start(screen)

	operationOne := new(operationMock)
	operationTwo := new(operationMock)
	operationThree := new(operationMock)

	texture.On("Bounds").Return(image.Rectangle{})
	operationOne.On("Do", texture).Return(false)
	operationTwo.On("Do", texture).Return(true)
	operationThree.On("Do", texture).Return(true)

	assert.Empty(t, loop.Mq.Operations)
	loop.Post(operationOne)
	loop.Post(operationTwo)
	loop.Post(operationThree)
	time.Sleep(1 * time.Second)
	assert.Empty(t, loop.Mq.Operations)

	operationOne.AssertCalled(t, "Do", texture)
	operationTwo.AssertCalled(t, "Do", texture)
	operationThree.AssertCalled(t, "Do", texture)
	receiver.AssertCalled(t, "Update", texture)
	screen.AssertCalled(t, "NewTexture", image.Pt(800, 800))
}

func TestMessageQueue_Push(t *testing.T) {
	Mq := &painter.MessageQueue{}

	operationOne := &operationQueueMock{}
	Mq.Push(operationOne)
	if len(Mq.Operations) != 1 {
		t.Errorf("Expected queue length to be 1, but got %d", len(Mq.Operations))
	}
	if !reflect.DeepEqual(operationOne, Mq.Operations[0]) {
		t.Error("Expected pushed operation to be in the queue")
	}

	operationTwo := &operationQueueMock{}
	Mq.Push(operationTwo)
	if len(Mq.Operations) != 2 {
		t.Errorf("Expected queue length to be 2, but got %d", len(Mq.Operations))
	}
	if !reflect.DeepEqual(operationTwo, Mq.Operations[0]) {
		t.Error("Expected pushed operation to be in the queue")
	}
}

func TestMessageQueue_Pull(t *testing.T) {
	Mq := &painter.MessageQueue{}

	operationOne := &operationQueueMock{}
	go func() {
		time.Sleep(50 * time.Millisecond)
		Mq.Push(operationOne)
	}()
	start := time.Now()
	op := Mq.Pull()
	elapsed := time.Since(start)

	if !reflect.DeepEqual(op, operationOne) {
		t.Errorf("Expected pulled operation to be the same as the pushed")
	}
	if elapsed < 50*time.Millisecond {
		t.Errorf("Expected Pull to block when pulling from an empty queue")
	}
	if len(Mq.Operations) != 0 {
		t.Errorf("Expected queue length to be 0, but got %d", len(Mq.Operations))
	}
	operationTwo := &operationQueueMock{}
	operationThree := &operationQueueMock{}
	Mq.Push(operationTwo)
	Mq.Push(operationThree)
	op = Mq.Pull()
	if len(Mq.Operations) != 1 {
		t.Errorf("Expected queue length to be 1, but got %d", len(Mq.Operations))
	}
	if !reflect.DeepEqual(op, operationTwo) {
		t.Error("Expected pulled operation to be the first pushed operation")
	}
}

type receiverMock struct {
	mock.Mock
}

func (rm *receiverMock) Update(t screen.Texture) {
	rm.Called(t)
}

type screenMock struct {
	mock.Mock
}

func (sm *screenMock) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (sm *screenMock) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

func (sm *screenMock) NewTexture(size image.Point) (screen.Texture, error) {
	args := sm.Called(size)
	return args.Get(0).(screen.Texture), args.Error(1)
}

type textureMock struct {
	mock.Mock
}

func (tm *textureMock) Release() {
	tm.Called()
}

func (tm *textureMock) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	tm.Called(dp, src, sr)
}

func (tm *textureMock) Bounds() image.Rectangle {
	args := tm.Called()
	return args.Get(0).(image.Rectangle)
}

func (tm *textureMock) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	tm.Called(dr, src, op)
}

func (tm *textureMock) Size() image.Point {
	args := tm.Called()
	return args.Get(0).(image.Point)
}

type operationMock struct {
	mock.Mock
}

func (om *operationMock) Do(t screen.Texture) bool {
	args := om.Called(t)
	return args.Bool(0)
}

type operationQueueMock struct{}

func (m *operationQueueMock) Do(t screen.Texture) (ready bool) {
	return false
}
