//go:generate stringer -type=enums,Direction,Edge,State,PullMode

package gpio

// Just for nicer output
type enums uint

const (
	dummy enums = iota
)

type Direction uint

const (
	DirectionIn Direction = iota
	DirectionOut
)

type Edge uint

const (
	EdgeNone Edge = iota
	EdgeRising
	EdgeFalling
	EdgeBoth
)

type State uint

const (
	StateHigh State = iota
	StateLow
)

type PullMode uint

const (
	PullNone PullMode = iota
	PullUp
	PullDown
)
