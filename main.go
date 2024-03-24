package main

import (
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Settings struct {
	Fullscreen bool
	XmrigPath  string
	P2poolPath string
	Hostname   string
	ZMQ        int
	RPC        int
}

type Person struct {
	Name    string
	Website string
	Github  string
	Address string
}

type People struct {
	People []Person `json:"people"`
}

var addressToUse string

func main() {
	var settings Settings
	err, settings := load()
	if err != nil {
		fmt.Println(err)
	}
	a := app.NewWithID("com.example.dmvp2p")
	w := a.NewWindow("DMVP2P (Donate Monero Via P2Pool)")

	usernameInput := widget.NewEntry()
	usernameInput.SetPlaceHolder("Enter username...")

	xmrigPath := widget.NewEntry()
	if settings.XmrigPath != "" {
		xmrigPath.SetPlaceHolder(settings.XmrigPath)
	} else {
		xmrigPath.SetPlaceHolder("Path to xmrig executable")
	}

	p2poolPath := widget.NewEntry()
	if settings.P2poolPath != "" {
		p2poolPath.SetPlaceHolder(settings.P2poolPath)
	} else {
		p2poolPath.SetPlaceHolder("Path to p2pool executable")
	}

	paths := container.NewVBox(p2poolPath, xmrigPath, widget.NewButton("Save", func() {
		log.Println("Content was", p2poolPath.Text)
		settings.XmrigPath = xmrigPath.Text
		settings.P2poolPath = p2poolPath.Text
		save(settings)
	}))

	p2poolUrl := widget.NewEntry()
	if settings.Hostname != "" {
		p2poolUrl.SetPlaceHolder(settings.Hostname)
	} else {
		p2poolUrl.SetPlaceHolder("Url for p2pool")
	}

	p2poolRPC := widget.NewEntry()
	if settings.RPC != 0 {
		p2poolRPC.SetPlaceHolder(strconv.Itoa(settings.RPC))
	} else {
		p2poolRPC.SetPlaceHolder("RPC port for p2pool")
	}

	p2poolZMQ := widget.NewEntry()
	if settings.ZMQ != 0 {
		p2poolZMQ.SetPlaceHolder(strconv.Itoa(settings.ZMQ))
	} else {
		p2poolZMQ.SetPlaceHolder("ZMQ port for p2pool")
	}

	p2poolUrlBtn := container.NewVBox(p2poolUrl, p2poolRPC, p2poolZMQ, widget.NewButton("Save", func() {
		settings.Hostname = p2poolUrl.Text
		settings.ZMQ, _ = strconv.Atoi(p2poolZMQ.Text)
		settings.RPC, _ = strconv.Atoi(p2poolRPC.Text)
		save(settings)
		log.Println("Content was", p2poolUrl.Text)
	}))

	radio := widget.NewRadioGroup([]string{}, func(value string) {
		_, addressToUse, _ = findPerson(value)
	})
	radio.Hide()

	result := widget.NewLabel("")

	searchButton := widget.NewButton("Search", func() {
		name, address, website := findPerson(usernameInput.Text)
		if name == "" || address == "" {
			result.SetText("Nothing found")
			radio.Hide() // Hide radio options if nothing is found
		} else {
			radioOptions := []string{name}
			radioOptions = append(radioOptions) // Add some example options
			radio.Options = radioOptions
			radio.Show() // Show radio options if a valid result is found
			firstFive := address[:5]
			lastFive := address[len(address)-4:]
			result.SetText(fmt.Sprintf("Website: %s \nMonero Address: %s...%s\n", website, firstFive, lastFive))
		}
	})

	startMiningButton := widget.NewButton("Start donating/mining", func() {
		err, cmd := startP2pool(addressToUse, settings)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cmd)
		err, cmd = startMining(xmrigPath.Text)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cmd)
	})

	startMiningButton.Importance = widget.HighImportance

	fullScreenButton := widget.NewButton("Fullscreen", func() {
		switch w.FullScreen() {
		case true:
			settings.Fullscreen = false
			save(settings)
			w.SetFullScreen(false)
		default:
			settings.Fullscreen = true
			save(settings)
			w.SetFullScreen(true)
		}

	})

	name := widget.NewLabel("...")
	name.TextStyle = fyne.TextStyle{Bold: true}

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		})
	list.OnSelected = func(id widget.ListItemID) {
		name.Text = data[id]
		name.Refresh()
	}

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), container.NewVBox(usernameInput, searchButton, radio, result, startMiningButton)),
		container.NewTabItemWithIcon("People", theme.ComputerIcon(), container.NewHSplit(list, name)),
		container.NewTabItemWithIcon("CCS", theme.MenuIcon(), widget.NewLabel("Coming soon...")),
		container.NewTabItemWithIcon("Logs", theme.WarningIcon(), widget.NewLabel("g")),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), container.NewVBox(fullScreenButton, paths, p2poolUrlBtn)),
	)
	w.SetIcon(theme.ComputerIcon())
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(900, 750))
	w.ShowAndRun()
}
