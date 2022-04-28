package main

/* CustomLabel is a widget type to display text with a background and using a customizable style */


/* ================================================================================ Imports */
import (
	"strings"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
)


/* ================================================================================ Public types */
type PaintStyle struct {
	Foreground, Background, Stroke color.RGBA
	StrokeWidth                    float32
}


type Paddings struct {
	Top, Bottom, Left, Right float32
}


type CustomLabel struct {
	widget.BaseWidget
	Alignment                        fyne.TextAlign
	Style                            PaintStyle
	LineWrapping                     bool
	Text                             string
	TextSize                         float32
	TextStyle                        fyne.TextStyle
	BackgroundPaddings, TextPaddings Paddings
}


/* ================================================================================ Private types */
type customLabelRenderer struct {
	background   *canvas.Rectangle
	textCanvases *[]*canvas.Text
	w            *CustomLabel
}


/* ================================================================================ Public functions */
func CalculatePaddings(paddingMultipliers, textPaddingOffsets Paddings) (backgroundPaddings, textPaddings Paddings) {
	backgroundPaddings = Paddings{
		theme.Padding() * paddingMultipliers.Top,  theme.Padding() * paddingMultipliers.Bottom,
		theme.Padding() * paddingMultipliers.Left, theme.Padding() * paddingMultipliers.Right,
	}

	textPaddings = Paddings{
		backgroundPaddings.Top  + textPaddingOffsets.Top,  backgroundPaddings.Bottom + textPaddingOffsets.Bottom,
		backgroundPaddings.Left + textPaddingOffsets.Left, backgroundPaddings.Right  + textPaddingOffsets.Right,
	}

	return backgroundPaddings, textPaddings
}


func NewCustomLabel(alignment fyne.TextAlign, style PaintStyle, lineWrapping bool, text string, textSize float32, textStyle fyne.TextStyle, paddingMultipliers, textPaddingOffsets Paddings) *CustomLabel {
	backgroundPaddings, textPaddings := CalculatePaddings(paddingMultipliers, textPaddingOffsets)

	customLabel := &CustomLabel{ Alignment: alignment, Style: style, LineWrapping: lineWrapping, Text: text, TextSize: textSize, TextStyle: textStyle, BackgroundPaddings: backgroundPaddings, TextPaddings: textPaddings }
	customLabel.ExtendBaseWidget(customLabel)

	return customLabel
}


/* ================================================================================ Private methods */
func wrapLine(line string, textSize float32, textStyle fyne.TextStyle, maxWidth float32) []string {
	/* Ensure there are at least two characters for splitting */
	if len(line) < 2 {
		return []string{ line }
	}

	/* Return if the whole line already fits into the available width */
	if fyne.MeasureText(line, textSize, textStyle).Width <= maxWidth {
		return []string{ line }
	}

	/* Ensure that at least the first character fits into the available width */
	if fyne.MeasureText(line[:1], textSize, textStyle).Width > maxWidth {
		return []string{ line }
	}

	/* Split the line, first by spaces (into words) */
	sep := " "
	for {
		parts     := strings.Split(line, sep)
		partCount := len(parts)

		/* Ensure that there are at least two parts */
		if partCount < 2 && sep != "" {
			/* Otherwise retry with splitting the line into characters instead of words */
			sep = ""
			continue
		}

		/* Find the maximum fitting sequence of parts using binary search */
		wrapIndex       := (partCount + 1) / 2
		wrapIndexStep   := 1
		rangeUpperIndex := partCount - 1
		rangeLowerIndex := 0

innerLoop:
		for {
			candidate := strings.Join(parts[:wrapIndex], sep)
			width     := fyne.MeasureText(candidate, textSize, textStyle).Width

			switch {
				case wrapIndexStep > 0 && width < maxWidth:
					/* The candidate sequence fits but is still converging, so use the current
					   wrap index as lower bound and continue in the center of the new range */
					rangeLowerIndex = wrapIndex

					/* Floor the increment step to avoid oscillating between the maximum fitting and the minimum non-fitting
					   candidate sequence if the index difference is just one - stay at the smaller/fitting one instead */
					wrapIndexStep = (rangeUpperIndex - rangeLowerIndex) / 2
					wrapIndex    += wrapIndexStep

				case wrapIndexStep > 0 && width > maxWidth:
					/* The candidate sequence does not fit but is still converging, so use the current
					   wrap index as upper bound and continue in the center of the new range */
					rangeUpperIndex = wrapIndex

					/* Ceil the decrement step to always be able to go back to a smaller/fitting
					   candidate sequence on overshoot, even if the index difference is just one */
					wrapIndexStep = (rangeUpperIndex - rangeLowerIndex + 1) / 2
					wrapIndex    -= wrapIndexStep

				case wrapIndex > 0 && width <= maxWidth:
					/* The candidate sequence fits and got stable at a length greater than zero, so
					   return it appended by the result of this function for the rest of the line */
					result     := []string{ candidate }
					rest       := strings.Join(parts[wrapIndex:], sep)
					restResult := wrapLine(rest, textSize, textStyle, maxWidth)
					return append(result, restResult...)

				case wrapIndex < 2 && sep != "":
					/* The candidate sequence got stable at a length of zero or one words (and the remaining
					   word does not fit), so retry with splitting the line into characters instead of words */
					sep = ""
					break innerLoop

				case wrapIndex < 2:
					/* The candidate sequence got stable at a length of zero or one characters (and the remaining character
					   does not fit), so just return it appended by the result of this function for the rest of the line */
					result     := []string{ candidate }
					rest       := strings.Join(parts[wrapIndex:], sep)
					restResult := wrapLine(rest, textSize, textStyle, maxWidth)
					return append(result, restResult...)
			}
		}
	}
}


func (w *CustomLabel) wrapLines(lines []string) []string {
	if w.LineWrapping {
		linesWrapped := make([]string, 0, len(lines))
		maxWidth     := w.Size().Width - w.TextPaddings.Left - w.TextPaddings.Right

		for _, line := range lines {
			subLines    := wrapLine(line, w.TextSize, w.TextStyle, maxWidth)
			linesWrapped = append(linesWrapped, subLines...)
		}

		return linesWrapped
	} else {
		return lines
	}
}


/* ================================================================================ Public rendering methods */
func (w *CustomLabel) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	background            := canvas.NewRectangle(w.Style.Background)
	background.StrokeColor = w.Style.Stroke
	background.StrokeWidth = w.Style.StrokeWidth

	lines        := strings.Split(w.Text, "\n")
	linesWrapped := w.wrapLines(lines)
	textCanvases := make([]*canvas.Text, len(linesWrapped))

	for i, line := range linesWrapped {
//		textCanvas := &canvas.Text{ Alignment: w.Alignment, Color: w.Color, Text: line, TextSize: w.TextSize, TextStyle: w.TextStyle }
		textCanvas := canvas.NewText(line, w.Style.Foreground)
		textCanvas.Alignment = w.Alignment
		textCanvas.TextSize  = w.TextSize
		textCanvas.TextStyle = w.TextStyle

		textCanvases[i] = textCanvas
	}

	return &customLabelRenderer{ background, &textCanvases, w }
}


func (r customLabelRenderer) Layout(size fyne.Size) {
	r.Refresh()

	backgroundSize := fyne.NewSize(size.Width - r.w.BackgroundPaddings.Left - r.w.BackgroundPaddings.Right, size.Height - r.w.BackgroundPaddings.Top - r.w.BackgroundPaddings.Bottom)
	r.background.Resize(backgroundSize)
	r.background.Move(fyne.NewPos(r.w.BackgroundPaddings.Left, r.w.BackgroundPaddings.Top))

	textSize     := fyne.NewSize(size.Width - r.w.TextPaddings.Left - r.w.TextPaddings.Right, size.Height - r.w.TextPaddings.Top - r.w.TextPaddings.Bottom)
	lineCount    := len(*r.textCanvases)
	lineHeight   := textSize.Height / float32(lineCount)
	heightOffset := r.w.TextPaddings.Top

	for _, textCanvas := range *r.textCanvases {
		textCanvas.Resize(fyne.NewSize(textSize.Width, lineHeight))
		textCanvas.Move(fyne.NewPos(r.w.TextPaddings.Left, heightOffset))
		heightOffset += lineHeight
	}
}


func (r customLabelRenderer) MinSize() fyne.Size {
	lines        := strings.Split(r.w.Text, "\n")
	linesWrapped := r.w.wrapLines(lines)
	maxLineWidth := float32(0)
	blockHeight  := float32(0)

	for _, line := range linesWrapped {
		lineSize    := fyne.MeasureText(line, r.w.TextSize, r.w.TextStyle)
		maxLineWidth = fyne.Max(maxLineWidth, lineSize.Width)
		blockHeight += lineSize.Height
	}

	return fyne.NewSize(maxLineWidth + r.w.TextPaddings.Left + r.w.TextPaddings.Right, blockHeight + r.w.TextPaddings.Top + r.w.TextPaddings.Bottom)
}


func (r customLabelRenderer) Refresh() {
	r.background.FillColor   = r.w.Style.Background
	r.background.StrokeColor = r.w.Style.Stroke
	r.background.Refresh()

	lines         := strings.Split(r.w.Text, "\n")
	linesWrapped  := r.w.wrapLines(lines)
	canvasesCount := len(*r.textCanvases)

	for i, line := range linesWrapped {
		var textCanvas *canvas.Text

		if i < canvasesCount {
			textCanvas       = (*r.textCanvases)[i]
			textCanvas.Text  = line
			textCanvas.Color = r.w.Style.Foreground
		} else {
			textCanvas = canvas.NewText(line, r.w.Style.Foreground)
		}

		textCanvas.Alignment = r.w.Alignment
		textCanvas.TextSize  = r.w.TextSize
		textCanvas.TextStyle = r.w.TextStyle

		if i < canvasesCount {
			(*r.textCanvases)[i] = textCanvas
			textCanvas.Refresh()
		} else {
			*r.textCanvases = append(*r.textCanvases, textCanvas)
		}
	}

	*r.textCanvases = (*r.textCanvases)[:len(linesWrapped)]
}


func (r customLabelRenderer) Objects() []fyne.CanvasObject {
	objectCount := len(*r.textCanvases) + 1
	objects     := make([]fyne.CanvasObject, objectCount)
	objects[0]   = r.background

	for i, textCanvas := range *r.textCanvases {
		objects[i + 1] = textCanvas
	}

	return objects
}


func (r customLabelRenderer) Destroy() {
}