package main

import (
	"fmt"
	"image/color"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/4rkal/dmvp2p/helpers"
	"github.com/4rkal/dmvp2p/pages"
)

var users []helpers.User

var settings pages.Settings

func init() {
	users = helpers.LoadUsers()
}

func createUserCard(user helpers.User) fyne.CanvasObject {
	createHyperlink := func(label string, link string) *widget.Hyperlink {
		parsedURL, err := url.Parse(link)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return nil
		}
		return widget.NewHyperlink(label, parsedURL)
	}

	mining := false

	var donateButton *widget.Button

	donateButton = widget.NewButton("Start Donating", func() {
		if !mining {
			fmt.Println("Starting mining for:", user.Name)
			err := helpers.StartMining("p2pmd.xmrvsbeast.com", user.Address, settings.XmrigPath, settings.P2poolPath)
			if err != nil {
				fmt.Println("Error starting mining:", err)
				return
			}
			donateButton.SetText("Stop Donating")
			donateButton.Importance = widget.DangerImportance
			mining = true
		} else {
			fmt.Println("Stopping mining for:", user.Name)
			err := helpers.StopMining()
			if err != nil {
				fmt.Println("Error stopping mining:", err)
				return
			}
			donateButton.SetText("Start Donating")
			donateButton.Importance = widget.HighImportance
			mining = false
		}
	})
	donateButton.Importance = widget.HighImportance

	var socialLinks []fyne.CanvasObject

	if user.GitHub != "" {
		socialLinks = append(socialLinks, createHyperlink("GitHub", "https://github.com/"+user.GitHub))
		socialLinks = append(socialLinks, widget.NewLabel("|"))
	}

	if user.X != "" {
		socialLinks = append(socialLinks, createHyperlink("X", "https://x.com/"+user.X))
		socialLinks = append(socialLinks, widget.NewLabel("|"))
	}

	if user.Website != "" {
		socialLinks = append(socialLinks, createHyperlink("Website", user.Website))
	}

	if len(socialLinks) > 0 {
		if label, ok := socialLinks[len(socialLinks)-1].(*widget.Label); ok && label.Text == "|" {
			socialLinks = socialLinks[:len(socialLinks)-1]
		}
	}

	return widget.NewCard(
		user.Name,
		user.Description,
		container.NewVBox(
			container.NewHBox(socialLinks...),
			donateButton,
		),
	)
}

func filterUsers(searchTerm string) []helpers.User {
	var filteredUsers []helpers.User
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(searchTerm)) {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

func init() {
	var err error
	settings, err = pages.LoadSettings()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	a := app.NewWithID("com.example.dmvp2p")
	w := a.NewWindow("DMVP2P (Donate Monero Via P2Pool)")

	fullScreenButton := pages.NewFullscreenButton(w)

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter user...")

	userList := container.NewVBox()

	updateUserList := func() {
		userList.Objects = nil

		searchTerm := input.Text

		filteredUsers := filterUsers(searchTerm)

		for _, user := range filteredUsers {
			userList.Add(createUserCard(user))
		}

		userList.Refresh()
	}

	input.OnChanged = func(content string) {
		updateUserList()
	}

	updateUserList()

	scrollContainer := container.NewScroll(userList)

	text := canvas.NewText("DMVP2P", color.White)
	text.TextSize = 50
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle = fyne.TextStyle{Bold: true}

	text2 := canvas.NewText("Donate Monero Via P2Pool", color.White)
	text2.Alignment = fyne.TextAlignCenter

	text3 := canvas.NewText("Mining Statistics", color.White)
	text3.TextSize = 30
	text3.Alignment = fyne.TextAlignCenter
	text3.TextStyle = fyne.TextStyle{Bold: true}

	text4 := canvas.NewText("Settings", color.White)
	text4.TextSize = 30
	text4.Alignment = fyne.TextAlignCenter
	text4.TextStyle = fyne.TextStyle{Bold: true}

	emtpy_line := widget.NewLabel("")

	p2pool_label := widget.NewLabel("P2Pool Path:")
	infoLabel := widget.NewLabel("No file selected.")
	if settings.P2poolPath != "" {
		infoLabel.SetText("Selected file: " + settings.P2poolPath)
	}

	selectFileButton := widget.NewButton("Select File", func() {
		helpers.SelectFileWithDialog(infoLabel, &settings.P2poolPath)
	})

	xmrig_label := widget.NewLabel("XMRig Path:")
	infoLabel2 := widget.NewLabel("No file selected.")
	if settings.XmrigPath != "" {
		infoLabel2.SetText("Selected file: " + settings.XmrigPath)
	}

	selectFileButton2 := widget.NewButton("Select File", func() {
		helpers.SelectFileWithDialog(infoLabel2, &settings.XmrigPath)
	})

	saveSettings := widget.NewButton("Save Settings", func() {
		pages.SaveSettings(settings)
	})

	scrollContainer.SetMinSize(fyne.NewSize(0, 700))

	userContainer := container.NewVBox(text, text2, input, scrollContainer)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), container.NewStack(userContainer)),
		container.NewTabItemWithIcon("Mining", theme.BrokenImageIcon(), container.NewVBox(text3)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), container.NewVBox(text4, fullScreenButton, emtpy_line, p2pool_label, infoLabel, selectFileButton, xmrig_label, infoLabel2, selectFileButton2, emtpy_line, saveSettings)),
	)

	w.SetIcon(theme.ComputerIcon())
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(900, 750))
	w.ShowAndRun()
}
