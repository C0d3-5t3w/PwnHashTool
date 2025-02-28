package gui

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/C0d3-5t3w/PwnHashTool/internal/utils"
)

func Launch(test fyne.App) {
	var mainApp fyne.App
	if test != nil {
		mainApp = test
	} else {
		mainApp = app.New()
	}

	window := mainApp.NewWindow("PwnHashTool")

	pcapInput := widget.NewEntry()
	pcapInput.SetPlaceHolder("Select PCAP file or directory...")
	hashInput := widget.NewEntry()
	hashInput.SetPlaceHolder("Select hash file or directory...")
	wordlistInput := widget.NewEntry()
	wordlistInput.SetPlaceHolder("Select wordlist...")
	potfileInput := widget.NewEntry()
	potfileInput.SetPlaceHolder("Select potfile or directory...")

	selectPcap := widget.NewButton("Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				pcapInput.SetText(uri.Path())
			}
		}, window)
	})

	selectHash := widget.NewButton("Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				hashInput.SetText(uri.Path())
			}
		}, window)
	})

	selectWordlist := widget.NewButton("Browse", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				wordlistInput.SetText(uri.URI().Path())
			}
		}, window)
	})

	selectPotfile := widget.NewButton("Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				potfileInput.SetText(uri.Path())
			}
		}, window)
	})

	status := widget.NewMultiLineEntry()
	status.Disable()
	status.Wrapping = fyne.TextWrapWord
	status.TextStyle = fyne.TextStyle{Monospace: true}

	statusScroll := container.NewScroll(status)
	statusScroll.SetMinSize(fyne.NewSize(380, 100))

	updateStatus := func(text string) {
		status.SetText(text)
		status.CursorRow = len(strings.Split(text, "\n")) - 1
	}

	convertButton := widget.NewButton("Convert PCAP(s)", func() {
		if pcapInput.Text == "" {
			updateStatus("Error: Please select a PCAP file or directory")
			return
		}

		updateStatus("Converting PCAP file(s)...")
		fileInfo, err := os.Stat(pcapInput.Text)
		if err != nil {
			updateStatus("Error: " + err.Error())
			return
		}

		if fileInfo.IsDir() {
			outputs, err := utils.ProcessPcapDirectory(pcapInput.Text, nil)
			if err != nil {
				updateStatus("Error: " + err.Error())
			} else {
				updateStatus(fmt.Sprintf("Converted %d files successfully", len(outputs)))
			}
		} else {
			if output, err := utils.RunHcxPcapngTool(pcapInput.Text, nil); err != nil {
				updateStatus("Error: " + err.Error())
			} else {
				updateStatus("Converted: " + output)
			}
		}
	})

	parseButton := widget.NewButton("Extract Passwords", func() {
		if potfileInput.Text == "" {
			updateStatus("Error: Please select a potfile or directory")
			return
		}

		fileInfo, err := os.Stat(potfileInput.Text)
		if err != nil {
			updateStatus("Error: " + err.Error())
			return
		}

		if fileInfo.IsDir() {
			outputs, err := utils.ProcessPotfileDirectory(potfileInput.Text)
			if err != nil {
				updateStatus("Error: " + err.Error())
			} else {
				updateStatus(fmt.Sprintf("Processed %d potfiles successfully", len(outputs)))
			}
		} else {
			if output, err := utils.ParsePotfile(potfileInput.Text); err != nil {
				updateStatus("Error parsing potfile: " + err.Error())
			} else {
				updateStatus("Passwords extracted to: " + output)
			}
		}
	})

	crackButton := widget.NewButton("Run Hashcat", func() {
		if hashInput.Text == "" {
			updateStatus("Error: Please select a hash file or directory")
			return
		}
		if wordlistInput.Text == "" {
			updateStatus("Error: Please select a wordlist")
			return
		}

		updateStatus("Running hashcat...")
		fileInfo, err := os.Stat(hashInput.Text)
		if err != nil {
			updateStatus("Error: " + err.Error())
			return
		}

		if fileInfo.IsDir() {
			outputs, err := utils.ProcessHashDirectory(hashInput.Text, wordlistInput.Text, nil)
			if err != nil {
				updateStatus("Error: " + err.Error())
				parseButton.Disable()
			} else {
				updateStatus(fmt.Sprintf("Processed %d hash files successfully", len(outputs)))
				parseButton.Enable()
			}
		} else {
			if output, err := utils.RunHashcat(hashInput.Text, wordlistInput.Text, nil); err != nil {
				updateStatus("Error: " + err.Error())
				parseButton.Disable()
			} else {
				updateStatus("Complete: " + output)
				parseButton.Enable()
			}
		}
	})

	content := container.NewVBox(
		widget.NewLabel("1. Convert PCAP:"),
		container.NewBorder(nil, nil, nil, selectPcap, pcapInput),
		convertButton,
		widget.NewLabel("2. Select Hash File:"),
		container.NewBorder(nil, nil, nil, selectHash, hashInput),
		widget.NewLabel("3. Run Hashcat:"),
		container.NewBorder(nil, nil, nil, selectWordlist, wordlistInput),
		crackButton,
		widget.NewLabel("4. Extract Passwords:"),
		container.NewBorder(nil, nil, nil, selectPotfile, potfileInput),
		parseButton,
		widget.NewLabel("Status:"),
		statusScroll,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 500))

	if test == nil {
		window.ShowAndRun()
	} else {
		window.Show()
	}
}
