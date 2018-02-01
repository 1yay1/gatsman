package gatsman

import "fmt"

//City struct to store the x and y values
type City struct {
	X float64
	Y float64
}

//String method for City struct
func (c *City) String() string {
	return fmt.Sprintf("[%v:%v]", c.X, c.Y)
}
