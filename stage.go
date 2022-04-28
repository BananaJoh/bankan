package main

/* Stage is a widget type describing a column/category of a board, which contains and manages items */


/* ================================================================================ Imports */
import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)


/* ================================================================================ Public types */
type Stage struct {
	widget.BaseWidget `json:"-"`
	Title  string
	Items  []*Item
}


/* ================================================================================ Private types */
type stageRenderer struct {
	titleLabel      *CustomLabel
	toolbar         *widget.Toolbar
	scrollArea      *container.Scroll
	itemContainer   *fyne.Container
	rightSeparator  *widget.Separator
	bottomSeparator *widget.Separator
	w               *Stage
}


/* ================================================================================ Public functions */
func NewStage(title string) *Stage {
	stage := &Stage{ Title: title }
	stage.ExtendBaseWidget(stage)

	return stage
}


/* ================================================================================ Public methods */
func (w *Stage) ItemIndex(toFind *Item) int {
	for i, item := range w.Items {
		if item == toFind {
			return i
		}
	}
	return -1
}


func (w *Stage) ItemAtPosition(position fyne.Position) *Item {
	for _, item := range w.Items {
		itemRect := Rectangle{ item.Position(), item.Size() }

		if itemRect.Contains(position) {
			return item
		}
	}
	return nil
}


func (w *Stage) AppendItem(title string, tags []Tag, description string, style ItemStyle) {
	item := NewItem(title, tags, description, style)

	w.Items = append(w.Items, item)
	w.Refresh()
}


func (w *Stage) InsertItem(after bool, reference *Item, title string, tags []Tag, description string, style ItemStyle) bool {
	i := w.ItemIndex(reference)
	if i < 0 {
		return false
	}
	if after {
		i++
	}

	w.Items = append(w.Items, nil)
	copy(w.Items[i+1:], w.Items[i:])

	item := NewItem(title, tags, description, style)
	w.Items[i] = item
	w.Refresh()

	return true
}


func (w *Stage) RemoveItem(toRemove *Item) bool {
	i := w.ItemIndex(toRemove)
	if i < 0 {
		return false
	}

	w.Items = append(w.Items[:i], w.Items[i+1:]...)
	w.Refresh()

	return true
}


func (w *Stage) ShowCreateItemDialog() {
	ShowItemDialog("New", "", "", "", ItemStyle{ color.RGBA{ 0, 0, 0, 255 }, color.RGBA{ 255, 255, 153, 255 } },
		func(title, tagEditString, description string, style ItemStyle) {
			w.AppendItem(title, ParseTagEditString(tagEditString), description, style)
		},
	)
}


func (w *Stage) ShowEditStageTitleDialog() {
	ShowEntryDialog("Edit Stage Title", "Title ...", w.Title,
		func(text string) {
			w.Title = text
			w.Refresh()
		},
	)
}


func (w *Stage) ShowRemoveStageConfirmDialog() {
	ShowConfirmDialog("Remove Stage", "This will remove the stage and all contained items from the board.\n\nAre you sure?\n",
		func() {
			board.RemoveStage(w)
		},
	)
}


func (w *Stage) ShowStageMenu() {
	menu := widget.NewPopUpMenu(
		fyne.NewMenu("Stage", 
			fyne.NewMenuItem("Edit Stage Title", w.ShowEditStageTitleDialog),
			fyne.NewMenuItem("Remove Stage",     w.ShowRemoveStageConfirmDialog),
		), window.Canvas(),
	)
	menu.ShowAtPosition(fyne.NewPos(w.Position().X + w.Size().Width - menu.Size().Width - 30, w.Position().Y + menu.Size().Height + 10))
}


func (w *Stage) SetFilterTags(filterTags []Tag) {
	for _, item := range w.Items {
		item.SetFilterTags(filterTags)
	}
}


/* ================================================================================ Public rendering methods */
func (w *Stage) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	titleLabel := NewCustomLabel(fyne.TextAlignLeading, PaintStyle{ color.RGBA{ 255, 255, 255, 255 }, color.RGBA{ 0, 0, 0, 0 }, color.RGBA{ 0, 0, 0, 0 }, 0 }, false, w.Title, theme.TextSubHeadingSize(), fyne.TextStyle{ Italic: true }, Paddings{ 1.0, 1.0, 1.0, 1.0 }, Paddings{ 0.0, 0.0, 0.0, 0.0 })
	toolbar    := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), w.ShowCreateItemDialog),
		widget.NewToolbarAction(theme.MoreVerticalIcon(), w.ShowStageMenu),
	)

	itemContainer := container.NewVBox()
	for _, item := range w.Items {
		itemContainer.Add(item)
	}

	scrollArea := container.NewVScroll(itemContainer)

	return &stageRenderer{ titleLabel, toolbar, scrollArea, itemContainer, widget.NewSeparator(), widget.NewSeparator(), w }
}


func (r stageRenderer) Layout(size fyne.Size) {
	titleSize             := r.titleLabel.MinSize()
	toolbarSize           := r.toolbar.MinSize()
	headerHeight          := fyne.Max(titleSize.Height, toolbarSize.Height)
	rightSeparatorWidth   := theme.Padding() / 2
	bottomSeparatorHeight := theme.Padding() / 2

	r.titleLabel.Resize(fyne.NewSize(size.Width - toolbarSize.Width - theme.Padding(), headerHeight))
	r.titleLabel.Move(fyne.NewPos(0, 0))

	r.toolbar.Resize(fyne.NewSize(toolbarSize.Width, headerHeight))
	r.toolbar.Move(fyne.NewPos(size.Width - toolbarSize.Width - theme.Padding(), 0))

	r.scrollArea.Resize(fyne.NewSize(size.Width - theme.Padding(), size.Height - headerHeight - (2 * theme.Padding())))
	r.scrollArea.Move(fyne.NewPos(0, headerHeight + theme.Padding()))

	r.rightSeparator.Resize(fyne.NewSize(rightSeparatorWidth, size.Height - bottomSeparatorHeight))
	r.rightSeparator.Move(fyne.NewPos(size.Width - rightSeparatorWidth, 0))

	r.bottomSeparator.Resize(fyne.NewSize(size.Width, bottomSeparatorHeight))
	r.bottomSeparator.Move(fyne.NewPos(0, size.Height - bottomSeparatorHeight))
}


func (r stageRenderer) MinSize() fyne.Size {
	titleSize     := r.titleLabel.MinSize()
	toolbarSize   := r.toolbar.MinSize()
	containerSize := r.scrollArea.MinSize()

	minWidth  := fyne.Max(titleSize.Width + toolbarSize.Width, containerSize.Width)
	minHeight := containerSize.Height + fyne.Max(titleSize.Height, toolbarSize.Height)

	return fyne.NewSize(minWidth, minHeight)
}


func (r stageRenderer) Refresh() {
	r.titleLabel.Text = r.w.Title
	r.titleLabel.Refresh()

	for _, item := range r.itemContainer.Objects {
		r.itemContainer.Remove(item)
	}

	for _, item := range r.w.Items {
		r.itemContainer.Add(item)
	}
	r.scrollArea.Refresh()
}


func (r stageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{ r.titleLabel, r.toolbar, r.scrollArea, r.rightSeparator, r.bottomSeparator }
}


func (r stageRenderer) Destroy() {
}