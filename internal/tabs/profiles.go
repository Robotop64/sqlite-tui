package tabs

import (
	"fmt"

	style "github.com/Robotop64/sqlite-tui/internal/style"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgTree "github.com/charmbracelet/lipgloss/tree"
)

type ProfileTab struct {
	Name        string
	Profiles    []utils.Profile
	IdxSelected int
}

func (b ProfileTab) GetName() string {
	return b.Name
}

func (b ProfileTab) View(width, height int) string {
	hint_height := 2
	sidecolumn_width := width / 5
	// maincolumn_width := width - sidecolumn_width
	layout_height := height - hint_height

	sidecolumn := style.Box.
		Width(sidecolumn_width - 2*border).
		Height(layout_height - 2*border)
	// maincolumn := box.
	// Width(maincolumn_width - 2*border).
	// Height(layout_height - 2*border)

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

	//=Side Column========
	//-Explorer------------
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
	//---------------------
	sidebar = lipgloss.JoinVertical(
		lipgloss.Top,
		tab.Render(),
		explorer.Render(),
	)

	//=Hints===============
	hints = lipgloss.JoinHorizontal(
		lipgloss.Top,
		"tabs:\nalt+(</>)",
		"│\n│",
		"quit:\nctrl+c, q",
		"│\n│",
		"profiles:\n   ↑/↓",
	)
	//=====================

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			sidebar,
			content,
		),
		hints,
	)
}

func (b ProfileTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			b.IdxSelected = max(b.IdxSelected-1, 0)
			return b, nil
		case "down":
			b.IdxSelected = min(b.IdxSelected+1, len(b.Profiles)-1)
			return b, nil
		default:
			return b, nil
		}
	}
	return b, nil
}

func loadProfiles(t *lgTree.Tree, profiles *[]utils.Profile) {
	// paths := cfg.GetStringSlice("profiles.paths")

	// t.Root("Profiles:")

	// for i, path := range paths {
	// 	if profile, err := utils.LoadProfile(path); err != nil {
	// 		t.Child(
	// 			"Faulty Profile!",
	// 			lgTree.New().Child(fmt.Sprintf("@Position %d", i+1)),
	// 		)
	// 	} else {
	// 		t.Child(profile.GetString("profile.name"))
	// 	}

	// }

	// if len(paths) == 0 {
	// 	t.Child("...")
	// }

	// return len(paths)

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
