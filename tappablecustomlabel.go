package main

/* TappableCustomLabel is a type extending the CustomLabel to be tappable */


/* ================================================================================ Imports */
import (
	"fyne.io/fyne/v2"
)


/* ================================================================================ Public types */
type TappableCustomLabel struct {
	CustomLabel
	OnTapped func()
}


/* ================================================================================ Public functions */
func NewTappableCustomLabel(alignment fyne.TextAlign, style PaintStyle, lineWrapping bool, text string, textSize float32, textStyle fyne.TextStyle, paddingMultipliers, textPaddingOffsets Paddings, tapped func()) *TappableCustomLabel {
	backgroundPaddings, textPaddings := CalculatePaddings(paddingMultipliers, textPaddingOffsets)

	customLabel         := CustomLabel{ Alignment: alignment, Style: style, LineWrapping: lineWrapping, Text: text, TextSize: textSize, TextStyle: textStyle, BackgroundPaddings: backgroundPaddings, TextPaddings: textPaddings }
	tappableCustomLabel := &TappableCustomLabel{ customLabel, tapped }
	tappableCustomLabel.ExtendBaseWidget(tappableCustomLabel)

	return tappableCustomLabel
}


/* ================================================================================ Public methods */
func (w *TappableCustomLabel) Tapped(event *fyne.PointEvent) {
	if w.OnTapped != nil {
		w.OnTapped()
	}
}