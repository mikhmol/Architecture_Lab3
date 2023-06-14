package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циелі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправленя останнього разу у Receiver

	Mq      MessageQueue
	done    chan struct{}
	stopped bool
}

var size = image.Pt(800, 800)

func (loop *Loop) Start(s screen.Screen) {
	loop.next, _ = s.NewTexture(size)
	loop.prev, _ = s.NewTexture(size)

	loop.done = make(chan struct{})

	go func() {
		for !loop.stopped || !loop.Mq.isEmpty() {
			op := loop.Mq.Pull()
			update := op.Do(loop.next)
			if update {
				loop.Receiver.Update(loop.next)
				loop.next, loop.prev = loop.prev, loop.next
			}
		}
		close(loop.done)
	}()
}

func (loop *Loop) Post(op Operation) {
	loop.Mq.Push(op)
}

func (loop *Loop) StopAndWait() {
	loop.Post(OperationFunc(func(t screen.Texture) {
		loop.stopped = true
	}))
	loop.stopped = true
	<-loop.done
}

type MessageQueue struct {
	Operations []Operation
	mu         sync.Mutex
	blocked    chan struct{}
}

func (Mq *MessageQueue) Push(op Operation) {
	Mq.mu.Lock()
	defer Mq.mu.Unlock()

	Mq.Operations = append(Mq.Operations, op)

	if Mq.blocked != nil {
		close(Mq.blocked)
		Mq.blocked = nil
	}
}

func (Mq *MessageQueue) Pull() Operation {
	Mq.mu.Lock()
	defer Mq.mu.Unlock()

	for len(Mq.Operations) == 0 {
		Mq.blocked = make(chan struct{})
		Mq.mu.Unlock()
		<-Mq.blocked
		Mq.mu.Lock()
	}

	op := Mq.Operations[0]
	Mq.Operations[0] = nil
	Mq.Operations = Mq.Operations[1:]
	return op
}

func (Mq *MessageQueue) isEmpty() bool {
	Mq.mu.Lock()
	defer Mq.mu.Unlock()

	return len(Mq.Operations) == 0
}
