package main

/* Rectangle is a basic type describing a rectangle and providing a method to check whether a point is within its bounds */


/* ================================================================================ Imports */
import (
	"fyne.io/fyne/v2"
)


/* ================================================================================ Public types */
type Rectangle struct {
	TopLeft fyne.Position
	Size    fyne.Size
}


/* ================================================================================ Public functions */
func NewRectangle(topLeft fyne.Position, size fyne.Size) *Rectangle {
	return &Rectangle{ topLeft, size }
}


/* ================================================================================ Public methods */
func (r *Rectangle) Contains(point fyne.Position) bool {
	if (point.X >= r.TopLeft.X) && (point.X < (r.TopLeft.X + r.Size.Width)) && (point.Y >= r.TopLeft.Y) && (point.Y < (r.TopLeft.Y + r.Size.Height)) {
		return true
	}
	return false
}