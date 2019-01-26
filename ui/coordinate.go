package ui

type Coordinate struct {
	TopLeftXrel, TopLeftYrel, BottomRightXrel, BottomRightYrel int
	TopLeftXabs, TopLeftYabs, BottomRightXabs, BottomRightYabs int
}

func (c *Coordinate) Scale(maxX, maxY int) {
	if c.TopLeftXrel != 0 {
		c.TopLeftXabs = maxX + c.TopLeftXrel
	}
	if c.TopLeftYrel != 0 {
		c.TopLeftYabs = maxY + c.TopLeftYrel
	}
	if c.BottomRightXrel != 0 {
		c.BottomRightXabs = maxX + c.BottomRightXrel
	}
	if c.BottomRightYrel != 0 {
		c.BottomRightYabs = maxY + c.BottomRightYrel
	}
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
