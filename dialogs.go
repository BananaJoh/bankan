package main

/* This file contains helper functions to show different types and combinations of dialogs */


/* ================================================================================ Imports */
import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
)


/* ================================================================================ Public functions */
func ShowConfirmDialog(title, text string, confirmedCallback func()) {
	dialog.ShowConfirm(title, text,
		func(confirmed bool) {
			if confirmed && confirmedCallback != nil {
				confirmedCallback()
			}
		}, window,
	)
}


func ShowEntryDialog(title, placeholder, text string, confirmedCallback func(text string)) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	entry.SetText(text)

	dialogContainer := container.NewVBox(entry, canvas.NewText("", color.Black))

	dialog.ShowCustomConfirm(title, "OK", "Cancel", dialogContainer,
		func(confirmed bool) {
			if confirmed && confirmedCallback != nil {
				confirmedCallback(entry.Text)
			}
		}, window,
	)

	window.Canvas().Focus(entry)
}


func ShowColorPickerDialog(title, message string, preselected color.RGBA, confirmedCallback func(selected color.RGBA)) {
	colorPickerDialog := dialog.NewColorPicker(title, message,
		func(c color.Color) {
			if confirmedCallback != nil {
				confirmedCallback(ColorToRGBA(c))
			}
		}, window,
	)
	colorPickerDialog.Advanced = true
	colorPickerDialog.Show()
	colorPickerDialog.SetColor(preselected)
}


func ShowItemDialog(dialogPrefix, title, tagEditString, description string, style ItemStyle, confirmedCallback func(title, tagEditString, description string, style ItemStyle)) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Title ...")
	titleEntry.SetText(title)
	
	tagsEntry := widget.NewEntry()
	tagsEntry.SetPlaceHolder("Tag1=Value1; Tag2=Value2; ...")
	tagsEntry.SetText(tagEditString)

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Description ...")
	descriptionEntry.SetText(description)

	foregroundColor       := style.Foreground
	foregroundColorButton := widget.NewButtonWithIcon("Foregound", theme.ColorPaletteIcon(),
		func() {
			ShowColorPickerDialog("Choose Foreground Color", "Please choose the color for item text and tag frames.", foregroundColor,
				func(selected color.RGBA) {
					foregroundColor = selected
				},
			)
		},
	)

	backgroundColor       := style.Background
	backgroundColorButton := widget.NewButtonWithIcon("Background", theme.ColorPaletteIcon(),
		func() {
			ShowColorPickerDialog("Choose Background Color", "Please choose the color for the item's background.", backgroundColor,
				func(selected color.RGBA) {
					backgroundColor = selected
				},
			)
		},
	)

	buttonContainer := container.NewGridWithColumns(2, foregroundColorButton, backgroundColorButton)
	dialogContainer := container.NewVBox(titleEntry, tagsEntry, descriptionEntry, buttonContainer, canvas.NewText("", color.Black))

	dialog.ShowCustomConfirm(dialogPrefix + " Item", "OK", "Cancel", dialogContainer,
		func(confirmed bool) {
			if confirmed && confirmedCallback != nil {
				confirmedCallback(titleEntry.Text, tagsEntry.Text, descriptionEntry.Text, ItemStyle{ foregroundColor, backgroundColor })
			}
		}, window,
	)

	window.Canvas().Focus(titleEntry)
}


func ShowFileOpenConfirmDialog(title, text string, defaultFileURI fyne.URI, confirmedCallback func(reader fyne.URIReadCloser)) {
	fileDialog := dialog.NewFileOpen(
		func(reader fyne.URIReadCloser, err error) {
			if reader != nil && err == nil {
				ShowConfirmDialog(title, text,
					func() {
						if confirmedCallback != nil {
							confirmedCallback(reader)
						}
					},
				)
			}
		}, window,
	)

	if defaultFileURI != nil {
		fileDialog.SetLocation(getParentListableURI(defaultFileURI))
	}

	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{ ".json" }))
	fileDialog.Show()
}


func ShowSaveAsDialog(defaultFileURI fyne.URI, confirmedCallback func(writer fyne.URIWriteCloser)) {
	fileDialog := dialog.NewFileSave(
		func(writer fyne.URIWriteCloser, err error) {
			if writer != nil && err == nil && confirmedCallback != nil {
				confirmedCallback(writer)
			}
		}, window,
	)

	if defaultFileURI != nil {
		fileDialog.SetFileName(defaultFileURI.Name())
		fileDialog.SetLocation(getParentListableURI(defaultFileURI))
	} else {
		fileDialog.SetFileName("bankan_board.json")
	}

	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{ ".json" }))
	fileDialog.Show()
}


/* ================================================================================ Private functions */
func getParentListableURI(file fyne.URI) fyne.ListableURI {
	dirURI, err := storage.Parent(file)
	if err != nil {
		return nil
	}

	dirListableURI, err := storage.ListerForURI(dirURI)
	if err != nil {
		return nil
	}

	return dirListableURI
}