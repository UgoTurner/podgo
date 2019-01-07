package ui

type Coordinate struct {
	TopLeftXrel, TopLeftYrel, BottomRightXrel, BottomRightYrel float32
	TopLeftXabs, TopLeftYabs, BottomRightXabs, BottomRightYabs int
}

func (c *Coordinate) Scale(maxX, maxY int) {
	c.TopLeftXabs = int(c.TopLeftXrel * float32(maxX))
	c.TopLeftYabs = int(c.TopLeftYrel * float32(maxY))
	c.BottomRightXabs = int(c.BottomRightXrel * float32(maxX))
	c.BottomRightYabs = int(c.BottomRightYrel * float32(maxY))
}

func (c *Coordinate) SetVisibility(isHidden bool) {
	var i = 1
	if isHidden {
		i = -1
	}
	c.TopLeftXabs = c.TopLeftXabs * i
	c.TopLeftYabs = c.TopLeftYabs * i
	c.BottomRightXabs = c.BottomRightXabs * i
	c.BottomRightYabs = c.BottomRightYabs * i
}
