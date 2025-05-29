package tabs

import (
	"fmt"

	style "github.com/Robotop64/sqlite-tui/internal/style"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"

	bubTxtIn "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgList "github.com/charmbracelet/lipgloss/list"
	lgTree "github.com/charmbracelet/lipgloss/tree"
)

type ProfileTab struct {
	Name                 string
	Profiles             []utils.Profile
	IdxSelected          int
	AddProfile           bool
	AddProfile_textInput bubTxtIn.Model
}

func (b ProfileTab) PostInit() ProfileTab {
	ti := bubTxtIn.New()
	ti.Placeholder = "Enter profile path..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30
	b.AddProfile_textInput = ti
	return b
}

func (b ProfileTab) Init() tea.Cmd {
	return bubTxtIn.Blink
}

func (b ProfileTab) GetName() string {
	return b.Name
}

func (b ProfileTab) View(width, height int) string {
	hint_height := 2
	sidecolumn_width := width / 5
	maincolumn_width := width - sidecolumn_width
	layout_height := height - hint_height

	sidecolumn := style.Box.
		Width(sidecolumn_width - 2*border).
		Height(layout_height - 2*border)
	maincolumn := style.Box.
		Width(maincolumn_width - 2*border).
		Height(layout_height - 2*border)

	tab := sidecolumn.
		Height(1).
		SetString(
			"⏴",
			lipgloss.PlaceHorizontal(lipgloss.Width(sidecolumn.String())-6, lipgloss.Center, b.Name),
			"⏵",
		)

	var (
		sidebar,
		content,
		hints string
	)

	//=Side Column=============
	//-Explorer----------------
	explorer_height := layout_height - lipgloss.Height(tab.String()) - 2*border

	tree := lgTree.New()
	loadProfiles(tree, &b.Profiles)
	tree.ItemStyleFunc(func(_ lgTree.Children, i int) lipgloss.Style {
		if len(b.Profiles) == 0 {
			return style.Normal
		}

		if i == b.IdxSelected {
			return style.Selected
		}
		return style.Normal
	})

	explorer := sidecolumn.
		Padding(0, 1).
		Height(explorer_height)

	explorer = explorer.SetString(tree.String())
	//-------------------------
	sidebar = lipgloss.JoinVertical(
		lipgloss.Top,
		tab.Render(),
		explorer.Render(),
	)

	//=Hints===================
	hints = lipgloss.JoinHorizontal(
		lipgloss.Top,
		"tabs:\nalt+(</>)",
		"│\n│",
		"quit:\nctrl+c, q",
		"│\n│",
		"profiles:\n   ↑/↓",
	)
	//=========================
	//=Content Column==========
	temp_content := maincolumn

	if b.Profiles[b.IdxSelected] != nil {
		title := style.Title.SetString("Profile Properties:").Render()
		list := lgList.New(
			fmt.Sprintf("Name: %s", b.Profiles[b.IdxSelected].GetString("profile.name")),
			"Database:",
			lgList.New(
				fmt.Sprintf("Path: %s", b.Profiles[b.IdxSelected].GetString("database.path")),
				fmt.Sprintf("Type: %s", b.Profiles[b.IdxSelected].GetString("database.type")),
			),
		).String()
		temp_content = temp_content.
			Padding(0, 1).
			SetString(
				lipgloss.JoinVertical(
					lipgloss.Top,
					title,
					list,
				),
			)
	}

	content = temp_content.
		Render()
	//=========================

	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			sidebar,
			content,
		),
		hints,
	)

	if !b.AddProfile {
		return layout
	} else {
		overlay, err := utils.Overlay(
			layout,
			addProfilePrompt(b).Render(),
			utils.Center, utils.Center,
		)
		if err != nil {
			return fmt.Sprintf("Error overlaying popup: %v", err)
		}
		return overlay
	}
}

func (b ProfileTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch b.AddProfile {
		case true:
			switch msg.String() {
			case "esc":
				b.AddProfile = false
				return b, nil
			}
			b.AddProfile_textInput, cmd = b.AddProfile_textInput.Update(msg)
			return b, cmd
		case false:
			switch msg.String() {
			case "up":
				b.IdxSelected = max(b.IdxSelected-1, 0)
				return b, nil
			case "down":
				b.IdxSelected = min(b.IdxSelected+1, len(b.Profiles)-1)
				return b, nil
			case "+":
				b.AddProfile = true
				return b, nil
			default:
				return b, nil
			}
		}
	}
	return b, nil
}

func loadProfiles(t *lgTree.Tree, profiles *[]utils.Profile) {
	for i, profile := range *profiles {
		if profile == nil {
			t.Child(fmt.Sprintf("Faulty Profile!\n@Position %d", i+1))
		} else {
			t.Child(profile.GetString("profile.name"))
		}
	}

	if len(*profiles) == 0 {
		t.Child("...")
	}
}

func addProfilePrompt(b ProfileTab) lipgloss.Style {
	title := style.Title.SetString("Add Profile:").Render()
	msg := style.Normal.SetString("Enter the path to the profile file:\n").Render()

	return style.Box.
		SetString(
			lipgloss.JoinVertical(
				lipgloss.Top,
				title,
				msg,
				b.AddProfile_textInput.View(),
			),
		)
}
