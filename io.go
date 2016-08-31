package gpio

import (
	"errors"
	"os"
	"time"
)

// Pin represents a single pin, which can be used either for reading or writing
type Pin struct {
	Number    uint
	direction Direction
	f         *os.File
}

// NewInput opens the given pin number for reading. The number provided should be the pin number known by the kernel
func NewPin(p uint) (*Pin, error) {
	pin := &Pin{
		Number: p,
	}
	err := exportGPIO(*pin)
	time.Sleep(10 * time.Millisecond)
	err = pin.SetOutput()
	err = pin.SetLow()
	return pin, err
}

// Close releases the resources related to Pin
func (p *Pin) Close() {
	p.f.Close()
}

func (p *Pin) SetOutput() error {
	p.direction = DirectionOut
	val, _ := p.Read()
	return setDirection(*p, DirectionOut, val)
}

func (p *Pin) SetInput() error {
	p.direction = DirectionIn
	val, _ := p.Read()
	return setDirection(*p, DirectionIn, val)
}

// Read returns the value read at the pin as reported by the kernel. This should only be used for input pins
func (p Pin) Read() (value State, err error) {
	return readPin(p)
}

// High sets the value of an output pin to logic high
func (p *Pin) SetHigh() error {
	if p.direction != DirectionOut {
		return errors.New("pin is not configured for output")
	}
	return writePin(*p, 1)
}

// Low sets the value of an output pin to logic low
func (p *Pin) SetLow() error {
	if p.direction != DirectionOut {
		return errors.New("pin is not configured for output")
	}
	return writePin(*p, 0)
}
