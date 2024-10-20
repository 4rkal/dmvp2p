package main

import (
	"fmt"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/4rkal/dmvp2p/helpers"
	"github.com/4rkal/dmvp2p/pages"
)

var users []helpers.User

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

	donateButton := widget.NewButton("Start Donating", func() {
		fmt.Printf("donate to: ", user.Name)
	})

	return widget.NewCard(
		user.Name,
		user.Description,
		container.NewVBox(
			container.NewHBox(
				createHyperlink("GitHub", "https://github.com/"+user.GitHub),
				widget.NewLabel("|"),
				createHyperlink("X", "https://x.com/"+user.X),
				widget.NewLabel("|"),
				createHyperlink("Website", user.Website),
			),
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

	scrollContainer.SetMinSize(fyne.NewSize(0, 700))

	userContainer := container.NewVBox(input, scrollContainer)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), container.NewStack(userContainer)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), container.NewVBox(fullScreenButton)),
	)

	w.SetIcon(theme.ComputerIcon())
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(900, 750))
	w.ShowAndRun()
}
