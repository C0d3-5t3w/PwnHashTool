package gui

import (
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
	pcapInput.SetPlaceHolder("Select PCAP file...")
	hashInput := widget.NewEntry()
	hashInput.SetPlaceHolder("Select hash file...")
	wordlistInput := widget.NewEntry()
	wordlistInput.SetPlaceHolder("Select wordlist...")
	potfileInput := widget.NewEntry()
	potfileInput.SetPlaceHolder("Select potfile...")

	selectPcap := widget.NewButton("Browse", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				pcapInput.SetText(uri.URI().Path())
			}
		}, window)
	})

	selectHash := widget.NewButton("Browse", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				hashInput.SetText(uri.URI().Path())
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
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				potfileInput.SetText(uri.URI().Path())
			}
		}, window)
	})

	status := widget.NewLabel("")

	convertButton := widget.NewButton("Convert PCAP", func() {
		if pcapInput.Text == "" {
			status.SetText("Error: Please select a PCAP file")
			return
		}

		status.SetText("Converting PCAP...")
		if output, err := utils.RunHcxPcapngTool(pcapInput.Text, nil); err != nil {
			status.SetText("Error: " + err.Error())
		} else {
			status.SetText("Converted: " + output)
		}
	})

	parseButton := widget.NewButton("Extract Passwords", func() {
		if potfileInput.Text == "" {
			status.SetText("Error: Please select a potfile")
			return
		}

		if output, err := utils.ParsePotfile(potfileInput.Text); err != nil {
			status.SetText("Error parsing potfile: " + err.Error())
		} else {
			status.SetText("Passwords extracted to: " + output)
		}
	})

	crackButton := widget.NewButton("Run Hashcat", func() {
		if hashInput.Text == "" {
			status.SetText("Error: Please select a hash file")
			return
		}
		if wordlistInput.Text == "" {
			status.SetText("Error: Please select a wordlist")
			return
		}

		status.SetText("Running hashcat...")
		if output, err := utils.RunHashcat(hashInput.Text, wordlistInput.Text, nil); err != nil {
			status.SetText("Error: " + err.Error())
			parseButton.Disable()
		} else {
			status.SetText("Complete: " + output)
			parseButton.Enable()
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
		status,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 300))

	if test == nil {
		window.ShowAndRun()
	} else {
		window.Show()
	}
}

// Author: C0d3-5t3w