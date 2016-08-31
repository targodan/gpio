package gpio

import (
	"fmt"
	"os"
	"strconv"
)

func exportGPIO(p Pin) error {
	export, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open gpio export file for writing: %s", err.Error())
	}
	defer export.Close()
	export.Write([]byte(strconv.Itoa(int(p.Number))))
	return nil
}

func unexportGPIO(p Pin) error {
	export, err := os.OpenFile("/sys/class/gpio/unexport", os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open gpio unexport file for writing: %s", err.Error())
	}
	defer export.Close()
	export.Write([]byte(strconv.Itoa(int(p.Number))))
	return nil
}

func setDirection(p Pin, d Direction, initialValue State) error {
	dir, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", p.Number), os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open gpio %d direction file for writing: %s", p.Number, err.Error())
	}
	defer dir.Close()

	switch {
	case d == DirectionIn:
		dir.Write([]byte("in"))
	case d == DirectionOut && initialValue == 0:
		dir.Write([]byte("low"))
	case d == DirectionOut && initialValue == 1:
		dir.Write([]byte("high"))
	default:
		return fmt.Errorf("setDirection called with invalid direction or initialValue, %d, %d", d, initialValue)
	}
	return nil
}

func setEdgeTrigger(p Pin, e Edge) error {
	edge, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/edge", p.Number), os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open gpio %d edge file for writing\n", p.Number)
	}
	defer edge.Close()

	switch e {
	case EdgeNone:
		edge.Write([]byte("none"))
	case EdgeRising:
		edge.Write([]byte("rising"))
	case EdgeFalling:
		edge.Write([]byte("falling"))
	case EdgeBoth:
		edge.Write([]byte("both"))
	default:
		return fmt.Errorf("setEdgeTrigger called with invalid edge %d", e)
	}
	return nil
}

func openPin(p *Pin, write bool) error {
	flags := os.O_RDONLY
	if write {
		flags = os.O_RDWR
	}
	f, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.Number), flags, 0600)
	if err != nil {
		return fmt.Errorf("failed to open gpio %d value file for reading\n", p.Number)
	}
	p.f = f
	return nil
}

func readPin(p Pin) (val State, err error) {
	file := p.f
	file.Seek(0, 0)
	buf := make([]byte, 1)
	_, err = file.Read(buf)
	if err != nil {
		return 0, err
	}
	c := buf[0]
	switch c {
	case '0':
		return 0, nil
	case '1':
		return 1, nil
	default:
		return 0, fmt.Errorf("read inconsistent value in pinfile, %c", c)
	}
}

func writePin(p Pin, v State) error {
	var buf []byte
	switch v {
	case 0:
		buf = []byte{'0'}
	case 1:
		buf = []byte{'1'}
	default:
		return fmt.Errorf("invalid output value %d", v)
	}
	_, err := p.f.Write(buf)
	return err
}
