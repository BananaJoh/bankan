package main

/* Item is a draggable widget type describing an expandable entry inside a stage, which holds the actual task information */


/* ================================================================================ Imports */
import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
)


/* ================================================================================ Public types */
type ItemStyle struct {
	Foreground, Background color.RGBA
}


type Item struct {
	widget.BaseWidget               `json:"-"`
	Title             string
	Description       string
	Tags              []Tag
	Style             ItemStyle
	Expanded          bool
	dragActive        bool          `json:"-"`
	dragStartPosition fyne.Position `json:"-"`
	dragEndPosition   fyne.Position `json:"-"`
}


/* ================================================================================ Private types */
type itemRenderer struct {
	background        *canvas.Rectangle
	titleLabel        *TappableCustomLabel
	toolbarBackground *canvas.Circle
	toolbar           *widget.Toolbar
	tagLabels         *[]*TappableCustomLabel
	descriptionLabel  *TappableCustomLabel
	w                 *Item
}


/* ================================================================================ Public functions */
func NewItem(title string, tags []Tag, description string, style ItemStyle) *Item {
	item := &Item{ Title: title, Tags: tags, Description: description, Style: style, Expanded: false }
	item.ExtendBaseWidget(item)

	return item
}


/* ================================================================================ Public methods */
func (w *Item) NewTagLabel(tag Tag) *TappableCustomLabel {
	return NewTappableCustomLabel(fyne.TextAlignCenter, PaintStyle{ w.Style.Background, w.Style.Foreground, color.RGBA{ 0, 0, 0, 0 }, 1 }, false, tag.DisplayString(), theme.CaptionTextSize(), fyne.TextStyle{ Italic: true }, Paddings{ 0.0, 1.0, 1.0, 0.5 }, Paddings{ 0.0, 0.0, 2.0, 2.0 },
		func() {
			board.ToggleFilterTag(tag)
		},
	)
}


func (w *Item) ShowEditItemDialog() {
	ShowItemDialog("Edit", w.Title, ComposeTagEditString(w.Tags), w.Description, w.Style,
		func(title, tagEditString, description string, style ItemStyle) {
			w.Title       = title
			w.Tags        = ParseTagEditString(tagEditString)
			w.Description = description
			w.Style       = style
			w.Refresh()
		},
	)
}


func (w *Item) ShowRemoveItemConfirmDialog() {
	ShowConfirmDialog("Remove Item", "This will remove the item from the board.\n\nAre you sure?\n",
		func() {
			board.RemoveItem(w)
		},
	)
}


func (w *Item) ShowItemMenu() {
	menu := widget.NewPopUpMenu(
		fyne.NewMenu("Item", 
			fyne.NewMenuItem("Edit Item", w.ShowEditItemDialog),
			fyne.NewMenuItem("Remove Item", w.ShowRemoveItemConfirmDialog),
		), window.Canvas(),
	)

	stage := board.ItemStage(w)
	if stage == nil {
		return
	}

	menu.ShowAtPosition(fyne.NewPos(stage.Position().X + w.Position().X + w.Size().Width + - menu.Size().Width - 18, stage.Position().Y + w.Position().Y + menu.Size().Height + 38))
}


func (w *Item) SetFilterTags(filterTags []Tag) {
	match := false

	for _, filterTag := range filterTags {
		for _, tag := range w.Tags {
			if tag == filterTag {
				match = true
				break
			}
		}
	}

	if match || len(filterTags) < 1 {
		w.Show()
	} else {
		w.Hide()
	}
}


func (w *Item) ToggleExpanded() {
	w.Expanded = !w.Expanded
	w.Refresh()
}


func (w *Item) Dragged(event *fyne.DragEvent) {
	if !w.dragActive {
		w.dragActive        = true
		w.dragStartPosition = event.Position
	}
	w.dragEndPosition = event.Position
}


func (w *Item) DragEnd() {
	w.dragActive = false

	itemRect := Rectangle{ fyne.NewPos(w.Position().X, 0), w.Size() }
	if !itemRect.Contains(w.dragStartPosition) || itemRect.Contains(w.dragEndPosition) {
		return
	}

	sourceStage := board.ItemStage(w)
	if sourceStage == nil {
		return
	}

	boardRelativeEndPosition := fyne.NewPos(sourceStage.Position().X + w.Position().X + w.dragEndPosition.X, sourceStage.Position().Y + w.Position().Y + w.dragEndPosition.Y)
	targetStage              := board.StageAtPosition(boardRelativeEndPosition)
	if targetStage == nil {
		return
	}

	targetStageRelativeEndPosition := fyne.NewPos(boardRelativeEndPosition.X - targetStage.Position().X, boardRelativeEndPosition.Y - targetStage.Position().Y)
	targetItem                     := targetStage.ItemAtPosition(targetStageRelativeEndPosition)

	if targetItem != nil {
		targetItemRelativeEndY := targetStageRelativeEndPosition.Y - targetItem.Position().Y
		targetItemHeightMidY   := (targetItem.Size().Height / 2)
		if targetItemRelativeEndY < targetItemHeightMidY {
			targetStage.InsertItem(false, targetItem, w.Title, w.Tags, w.Description, w.Style)
		} else {
			targetStage.InsertItem(true, targetItem, w.Title, w.Tags, w.Description, w.Style)
		}
	} else {
		targetStage.AppendItem(w.Title, w.Tags, w.Description, w.Style)
	}

	sourceStage.RemoveItem(w)
}


/* ================================================================================ Public rendering methods */
func (w *Item) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	background := canvas.NewRectangle(w.Style.Background)
	titleLabel := NewTappableCustomLabel(fyne.TextAlignLeading, PaintStyle{ w.Style.Foreground, color.RGBA{ 0, 0, 0, 0 }, color.RGBA{ 0, 0, 0, 0 }, 0 }, true, w.Title, theme.TextSize(), fyne.TextStyle{ Bold: true }, Paddings{ 0.0, 0.25, 1.0, 0.0 }, Paddings{ 0.0, 0.0, 0.0, 0.0 }, w.ToggleExpanded)
	toolbarBackground := canvas.NewCircle(color.RGBA{ 0, 0, 0, 127 })
	toolbar           := widget.NewToolbar(widget.NewToolbarAction(theme.MoreVerticalIcon(), w.ShowItemMenu))
	tagLabels         := make([]*TappableCustomLabel, len(w.Tags))

	for i, tag := range w.Tags {
		tagLabels[i] = w.NewTagLabel(tag)
	}

	descriptionLabel := NewTappableCustomLabel(fyne.TextAlignLeading, PaintStyle{ w.Style.Foreground, color.RGBA{ 0, 0, 0, 0 }, color.RGBA{ 0, 0, 0, 0 }, 0 }, true, w.Description, theme.TextSize(), fyne.TextStyle{ Monospace: true }, Paddings{ 0.0, 1.0, 1.0, 0.5 }, Paddings{ 0.0, 0.0, 0.0, 0.0 }, w.ToggleExpanded)

	if !w.Expanded {
		descriptionLabel.Hide()
	}

	return &itemRenderer{ background, titleLabel, toolbarBackground, toolbar, &tagLabels, descriptionLabel, w }
}


func (r itemRenderer) Layout(size fyne.Size) {
	headerHeight             := Round(r.titleLabel.MinSize().Height)
	toolbarHeight            := fyne.MeasureText(r.w.Title, theme.TextSize(), fyne.TextStyle{ Bold: true }).Height
	toolbarWidth             := r.toolbar.MinSize().Width
	toolbarBackgroundPadding := theme.Padding() / 2
	toolbarBackgroundOffset  := toolbarBackgroundPadding / 2

	r.background.Resize(size)
	r.background.Move(fyne.NewPos(0, 0))

	r.titleLabel.Resize(fyne.NewSize(size.Width - toolbarWidth, headerHeight))
	r.titleLabel.Move(fyne.NewPos(0, 0))

	r.toolbarBackground.Resize(fyne.NewSize(toolbarWidth - toolbarBackgroundPadding, toolbarHeight - toolbarBackgroundPadding))
	r.toolbarBackground.Move(fyne.NewPos(size.Width - toolbarWidth + toolbarBackgroundOffset, ((headerHeight - toolbarHeight) / 2) + toolbarBackgroundOffset))

	r.toolbar.Resize(fyne.NewSize(toolbarWidth, toolbarHeight))
	r.toolbar.Move(fyne.NewPos(size.Width - toolbarWidth, (headerHeight - toolbarHeight) / 2))

	tagsLineWidth     := float32(0)
	tagsBlockHeight   := float32(0)
	tagsLineMaxHeight := float32(0)

	for _, tagLabel := range *r.tagLabels {
		tagSize := tagLabel.MinSize()

		if tagsLineWidth > 0 && (tagsLineWidth + tagSize.Width) > size.Width {
			tagsBlockHeight  += tagsLineMaxHeight
			tagsLineWidth     = 0
			tagsLineMaxHeight = 0
		}

		tagLabel.Resize(tagSize)
		tagLabel.Move(fyne.NewPos(tagsLineWidth, headerHeight + tagsBlockHeight))

		tagsLineWidth    += tagSize.Width
		tagsLineMaxHeight = fyne.Max(tagsLineMaxHeight, tagSize.Height)
	}
	tagsBlockHeight += tagsLineMaxHeight

	r.descriptionLabel.Resize(fyne.NewSize(size.Width, size.Height - headerHeight - tagsBlockHeight))
	r.descriptionLabel.Move(fyne.NewPos(0, headerHeight + tagsBlockHeight))
}


func (r itemRenderer) MinSize() fyne.Size {
	titleSize    := r.titleLabel.MinSize()
	headerHeight := Round(titleSize.Height)
	toolbarWidth := r.toolbar.MinSize().Width

	maxWidth          := r.w.Size().Width
	tagsLineWidth     := float32(0)
	tagsBlockHeight   := float32(0)
	tagsLineMaxWidth  := float32(0)
	tagsLineMaxHeight := float32(0)

	for _, tagLabel := range *r.tagLabels {
		tagSize := tagLabel.MinSize()

		if tagsLineWidth > 0 && (tagsLineWidth + tagSize.Width) > maxWidth {
			tagsBlockHeight  += tagsLineMaxHeight
			tagsLineWidth     = 0
			tagsLineMaxHeight = 0
		}
		tagsLineWidth    += tagSize.Width
		tagsLineMaxWidth  = fyne.Max(tagsLineMaxWidth, tagsLineWidth)
		tagsLineMaxHeight = fyne.Max(tagsLineMaxHeight, tagSize.Height)
	}
	tagsBlockHeight += tagsLineMaxHeight

	descriptionSize := r.descriptionLabel.MinSize()
	if !r.w.Expanded {
		descriptionSize.Height = 0
	}

	minWidth  := fyne.Max(tagsLineMaxWidth, fyne.Max(titleSize.Width + toolbarWidth, descriptionSize.Width))
	minHeight := headerHeight + tagsBlockHeight + descriptionSize.Height

	return fyne.NewSize(minWidth, Round(minHeight))
}


func (r itemRenderer) Refresh() {
	r.background.FillColor = r.w.Style.Background
	r.background.Refresh()

	r.titleLabel.Style.Foreground = r.w.Style.Foreground
	r.titleLabel.Text             = r.w.Title
	r.titleLabel.Refresh()

	tagLabelsCount := len(*r.tagLabels)

	for i, tag := range r.w.Tags {
		if i < tagLabelsCount {
			(*r.tagLabels)[i].Style.Foreground = r.w.Style.Background
			(*r.tagLabels)[i].Style.Background = r.w.Style.Foreground
			(*r.tagLabels)[i].Text             = tag.DisplayString()
			(*r.tagLabels)[i].Refresh()
		} else {
			*r.tagLabels = append(*r.tagLabels, r.w.NewTagLabel(tag))
		}
	}

	*r.tagLabels = (*r.tagLabels)[:len(r.w.Tags)]

	r.descriptionLabel.Style.Foreground = r.w.Style.Foreground
	r.descriptionLabel.Text             = r.w.Description
	r.descriptionLabel.Refresh()

	if r.w.Expanded {
		r.descriptionLabel.Show()
	} else {
		r.descriptionLabel.Hide()
	}
}


func (r itemRenderer) Objects() []fyne.CanvasObject {
	objectCount := len(*r.tagLabels) + 5
	objects     := make([]fyne.CanvasObject, objectCount)
	objects[0]   = r.background
	objects[1]   = r.titleLabel
	objects[2]   = r.toolbarBackground
	objects[3]   = r.toolbar

	for i, tagLabel := range *r.tagLabels {
		objects[i + 4] = tagLabel
	}

	objects[objectCount - 1] = r.descriptionLabel

	return objects
}


func (r itemRenderer) Destroy() {
}