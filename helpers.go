package main

/* This file contains general helper functions */


/* ================================================================================ Imports */
import (
	"image/color"
)


/* ================================================================================ Public functions */
func ColorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}	
}


func Round(f float32) float32 {
	return float32(int(f + 0.5))
}