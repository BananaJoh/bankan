package main

/* Board is the top-level widget type describing a kanban board, which contains and manages stages */


/* ================================================================================ Imports */
import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)


/* ================================================================================ Public types */
type Board struct {
	widget.BaseWidget                          `json:"-"`
	Name            string
	Stages          []*Stage
	FilterTags      []Tag                      `json:"-"`
	OnFilterChanged func(tagEditString string) `json:"-"`
}


/* ================================================================================ Private types */
type boardRenderer struct {
	stageContainer *fyne.Container
	w              *Board
}


/* ================================================================================ Public functions */
func NewBoard(name string, filterChanged func(tagEditString string)) *Board {
	board := &Board{ Name: name, OnFilterChanged: filterChanged }
	board.ExtendBaseWidget(board)

	return board
}


/* ================================================================================ Public methods */
func (w *Board) Clear() {
	w.Stages = w.Stages[:0]
	w.Refresh()
}


func (w *Board) Data() ([]byte, error) {
	return json.Marshal(w)
}


func (w *Board) Load(data []byte) error {
	if err := json.Unmarshal(data, w); err != nil {
		return err
	}
	w.Refresh()

	return nil
}


func (w *Board) StageIndex(toFind *Stage) int {
	for i, stage := range w.Stages {
		if stage == toFind {
			return i
		}
	}
	return -1
}


func (w *Board) ItemStageIndex(toFind *Item) int {
	for i, stage := range w.Stages {
		if stage.ItemIndex(toFind) >= 0 {
			return i
		}
	}
	return -1
}


func (w *Board) ItemStage(toFind *Item) *Stage {
	i := w.ItemStageIndex(toFind)
	if i < 0 {
		return nil
	}

	return w.Stages[i]
}


func (w *Board) StageAtPosition(position fyne.Position) *Stage {
	for _, stage := range w.Stages {
		stageRect := Rectangle{ stage.Position(), stage.Size() }

		if stageRect.Contains(position) {
			return stage
		}
	}
	return nil
}


func (w *Board) AppendStage(title string) {
	stage := NewStage(title)

	w.Stages = append(w.Stages, stage)
	w.Refresh()
}


func (w *Board) RemoveStage(toRemove *Stage) bool {
	i := w.StageIndex(toRemove)
	if i < 0 {
		return false
	}

	w.Stages = append(w.Stages[:i], w.Stages[i+1:]...)
	w.Refresh()

	return true
}


func (w *Board) RemoveItem(toRemove *Item) {
	for _, stage := range w.Stages {
		if stage.RemoveItem(toRemove) {
			w.Refresh()
			return
		}
	}
}


func (w *Board) ShowCreateStageDialog() {
	ShowEntryDialog("New Stage", "Title ...", "",
		func(text string) {
			w.AppendStage(text)
		},
	)
}


func (w *Board) ApplyTagFilter() {
	for _, stage := range w.Stages {
		stage.SetFilterTags(w.FilterTags)
	}
}


func (w *Board) SetTagFilter(tagEditString string) {
	w.FilterTags = ParseTagEditString(tagEditString)
	w.ApplyTagFilter()
}


func (w *Board) FilterTagIndex(toFind Tag) int {
	for i, filterTag := range w.FilterTags {
		if filterTag.Expression == toFind.Expression {
			return i
		}
	}
	return -1
}


func (w *Board) ToggleFilterTag(tag Tag) {
	i := w.FilterTagIndex(tag)

	if i < 0 {
		w.FilterTags = append(w.FilterTags, tag)
	} else {
		w.FilterTags = append(w.FilterTags[:i], w.FilterTags[i+1:]...)
	}
	w.ApplyTagFilter()

	if w.OnFilterChanged != nil {
		w.OnFilterChanged(ComposeTagEditString(w.FilterTags))
	}
}


/* ================================================================================ Public rendering methods */
func (w *Board) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	// TODO: for small screens
/*
	stageContainer := container.NewAppTabs(
		container.NewTabItem("Tab 1", widget.NewLabel("Hello")),
		container.NewTabItem("Tab 2", widget.NewLabel("World!")),
	)
	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))
	tabs.SetTabLocation(container.TabLocationTop)
*/

	stageContainer := container.NewWithoutLayout()

	stageCount := len(w.Stages)
	if stageCount > 0 {
		stageContainer.Layout = layout.NewGridLayout(stageCount)

		for _, stage := range w.Stages {
			stageContainer.Add(stage)
		}
	}

	return &boardRenderer{ stageContainer, w }
}


func (r boardRenderer) Layout(size fyne.Size) {
	r.stageContainer.Resize(size)
	r.stageContainer.Move(fyne.NewPos(0, 0))
}


func (r boardRenderer) MinSize() fyne.Size {
	containerSize := r.stageContainer.MinSize()

	return fyne.NewSize(containerSize.Width, containerSize.Height)
}


func (r boardRenderer) Refresh() {
	for _, stage := range r.stageContainer.Objects {
		r.stageContainer.Remove(stage)
	}

	r.stageContainer.Layout = layout.NewGridLayout(len(r.w.Stages))

	for _, stage := range r.w.Stages {
		r.stageContainer.Add(stage)
	}
}


func (r boardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{ r.stageContainer }
}


func (r boardRenderer) Destroy() {
}