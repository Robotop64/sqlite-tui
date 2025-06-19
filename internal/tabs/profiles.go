package tabs

import (
	"fmt"
	"path/filepath"
	"strings"

	Focus "github.com/Robotop64/sqlite-tui/internal/enums/ui"
	style "github.com/Robotop64/sqlite-tui/internal/style"
	color "github.com/Robotop64/sqlite-tui/internal/style/color"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"
	yaml "gopkg.in/yaml.v3"

	bubTxtEdit "github.com/charmbracelet/bubbles/textarea"
	bubTxtIn "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgList "github.com/charmbracelet/lipgloss/list"
	cfg "github.com/spf13/viper"
)

type ProfileTab struct {
	Name        string
	ElemFocus   Focus.UiFocus
	ChangeFocus bool

	IdxFocus    int
	IdxSelected int

	AddProfile  bubTxtIn.Model
	ViewProfile bubTxtEdit.Model
}

func (b ProfileTab) PostInit() ProfileTab {
	b.IdxSelected = utils.Configs.Profiles.LastUsed

	txtinput := bubTxtIn.New()
	txtinput.Placeholder = "..."
	txtinput.Width = 256
	txtinput.CharLimit = 256
	b.AddProfile = txtinput

	txtedit := bubTxtEdit.New()
	b.ViewProfile = txtedit

	if len(utils.Profiles) > 0 {

		profile := utils.ActiveProfile()
		data, _ := yaml.Marshal(profile)
		b.ViewProfile.SetValue(string(data))
	}

	return b
}

func (b ProfileTab) Init() tea.Cmd {
	if b.ElemFocus == Focus.TxtInput {
		return bubTxtEdit.Blink
	}
	if b.ElemFocus == Focus.TxtInput {
		return bubTxtIn.Blink
	}

	return nil
}

func (b ProfileTab) GetName() string {
	return b.Name
}

func (b ProfileTab) View(width, height int) string {
	hint_height := 2

	popup_size := utils.Dimensions{
		Width:  width * 3 / 5,
		Height: 0,
	}

	explorer_size := utils.Dimensions{
		Width:  width / 5,
		Height: height - hint_height,
	}

	content_size := utils.Dimensions{
		Width:  width - explorer_size.Width,
		Height: height - hint_height,
	}

	sidecolumn := style.Box.
		Width(explorer_size.Width - 2*border).
		Height(explorer_size.Height - 2*border)
	// maincolumn := style.Box.
	// 	Width(content_size.Width - 2*border).
	// 	Height(content_size.Height - 2*border)

	tab := sidecolumn.
		Height(1).
		SetString(
			"⏴",
			lipgloss.PlaceHorizontal(lipgloss.Width(sidecolumn.String())-6, lipgloss.Center, b.Name),
			"⏵",
		).
		Foreground(color.TextHighlight)

	var (
		sidebar,
		hints string
	)

	//=Side Column=============
	//-Explorer----------------
	explorer_height := explorer_size.Height - lipgloss.Height(tab.String())

	//-------------------------
	sidebar = lipgloss.JoinVertical(
		lipgloss.Top,
		tab.Render(),
		explorer(b, utils.Dimensions{Width: explorer_size.Width, Height: explorer_height}).Render(),
	)

	//=Hints===================
	hints = lipgloss.JoinHorizontal(
		lipgloss.Top,
		"quit:\nctrl+c|q",
		"│\n│",
		"save:\nctrl+s",
		"│\n│",
		"tabs:\nalt+(</>)",
		"│\n│",
		"profiles:\n   ↑/↓",
		"│\n│",
		"add:\n +",
		"│\n│",
		"remove:\n   -",
		"│\n│",
	)
	//=========================
	//=Content Column==========

	//=========================
	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			sidebar,
			editor(b, content_size).Render(),
		),
		hints,
	)

	if !(b.ElemFocus == Focus.TxtInput) {
		return layout
	} else {
		overlay, err := utils.Overlay(
			layout,
			addProfilePrompt(b, popup_size).Render(),
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "alt+left":
			b.ElemFocus = Focus.ProfileList
			b.ViewProfile.Blur()
			return b, nil
		case "alt+right":
			b.ElemFocus = Focus.TxtEdit
			return b, nil
		}

		switch b.ElemFocus {
		case Focus.None:
			return b, nil
		case Focus.ProfileList:
			switch msg.String() {
			case "up":
				b.IdxFocus = max(b.IdxFocus-1, 0)
				return b, nil
			case "down":
				b.IdxFocus = min(b.IdxFocus+1, len(utils.Profiles)-1)
				return b, nil
			case "+":
				b.ElemFocus = Focus.TxtInput
				return b, nil
			case "-":
				if len(utils.Profiles) > 0 && b.IdxFocus < len(utils.Profiles) {
					utils.Profiles = append(utils.Profiles[:b.IdxFocus], utils.Profiles[b.IdxFocus+1:]...)
					paths := cfg.GetStringSlice("profiles.paths")
					paths = append(paths[:b.IdxFocus], paths[b.IdxFocus+1:]...)
					cfg.Set("profiles.paths", paths)
				}
				b.IdxFocus = max(b.IdxFocus-1, 0)
				return b, nil
			case "enter":
				utils.Configs.Profiles.LastUsed = b.IdxFocus
				b.IdxSelected = b.IdxFocus
				profile := utils.ActiveProfile()
				data, _ := yaml.Marshal(profile)
				b.ViewProfile.SetValue(string(data))
				return b, nil
			case "c":
				if len(utils.Profiles) > 0 && b.IdxFocus < len(utils.Profiles) {

				}
				return b, nil
			default:
				return b, nil
			}
		case Focus.TxtInput:
			b.AddProfile.Focus()
			switch msg.String() {
			case "esc":
				b.ElemFocus = Focus.ProfileList
				b.AddProfile.Reset()
				return b, nil
			case "enter":
				raw_path := b.AddProfile.Value()
				if raw_path == "" {
					b.ElemFocus = Focus.ProfileList
					return b, nil
				}

				var path string
				if !utils.EndsWith(raw_path, "Profile.yaml") {
					path = filepath.Join(raw_path, "Profile.yaml")
				} else {
					path = raw_path
				}

				fileExists := utils.CheckPath(path)
				var err error
				if !fileExists {
					if _, err = utils.CreateProfile(b.AddProfile.Value()); err != nil {
						b.ElemFocus = Focus.ProfileList
						b.AddProfile.Reset()
						return b, nil
					}
				} else {
					if _, err = utils.LoadProfile(path); err != nil {
						b.ElemFocus = Focus.ProfileList
						b.AddProfile.Reset()
						return b, nil
					}
				}

				b.ElemFocus = Focus.ProfileList
				b.AddProfile.Reset()
				return b, nil
			}
			b.AddProfile, cmd = b.AddProfile.Update(msg)
			return b, cmd
		case Focus.TxtEdit:
			if len(utils.Profiles) == 0 {
				return b, nil
			}

			switch msg.String() {
			case "ctrl+s":
				data := b.ViewProfile.Value()
				if !(len(data) == 0) {
					profile := utils.ActiveProfile()

					if err := yaml.Unmarshal([]byte(data), profile); err != nil {
						fmt.Println("Error unmarshalling profile data:", err)
						return b, nil
					}
					if err := utils.SaveProfile(profile, profile.Path); err != nil {
						fmt.Println("Error saving profile:", err)
						return b, nil
					}
				}
				return b, nil
			default:
				if !b.ViewProfile.Focused() {
					cmd = b.ViewProfile.Focus()
					cmds = append(cmds, cmd)
				}
			}

			b.ViewProfile, cmd = b.ViewProfile.Update(msg)
			cmds = append(cmds, cmd)
			return b, tea.Batch(cmds...)
		}
	}
	return b, nil
}

func explorer(b ProfileTab, dim utils.Dimensions) lipgloss.Style {
	list := lgList.New()
	names := utils.Map(utils.Profiles, func(i int, p *utils.Profile) string {
		if p == nil {
			return "Faulty Profile!"
		} else {
			return p.Name
		}
	})

	view := style.Box.
		Padding(0, 1).
		Width(dim.Width - 2).
		Height(dim.Height - 2)

	if b.ElemFocus == Focus.ProfileList {
		view = view.BorderForeground(color.BoxSelected)
	}

	list.Items(names)

	list.Enumerator(func(l lgList.Items, i int) string {
		if i == b.IdxSelected {
			return ">"
		}
		return "•"
	})

	list.ItemStyleFunc(func(_ lgList.Items, i int) lipgloss.Style {
		if len(utils.Profiles) == 0 {
			return style.Normal
		}

		if i == b.IdxFocus {
			return style.Selected
		}
		return style.Normal
	})

	if len(names) == 0 {
		return view.SetString("...")
	}

	return view.SetString(list.String())
}

func addProfilePrompt(b ProfileTab, dim utils.Dimensions) lipgloss.Style {
	title := style.Title.SetString("Add Profile:").Render()
	msg := style.Normal.SetString("Enter the path to the profile file or directory:\n").Render()

	label_cancel, key_cancel := " Cancel ", "esc"
	label_confirm, key_confirm := " Confirm ", "enter"

	button_cancel := style.Button.SetString(label_cancel).Render() + "\n" +
		style.Normal.SetString(fmt.Sprintf(" (%s)", key_cancel)).Render()
	button_confirm := style.Button.SetString(label_confirm).Render() + "\n" +
		style.Normal.SetString(fmt.Sprintf(" (%s)", key_confirm)).Render()

	buffer_length := dim.Width - 4 - lipgloss.Width(button_cancel) - lipgloss.Width(button_confirm)

	hints := lipgloss.JoinHorizontal(
		lipgloss.Top,
		button_cancel,
		strings.Repeat(" ", buffer_length),
		button_confirm,
	)

	dim.Width = dim.Width - 2

	popup := style.Box.
		Padding(0, 1).
		Width(dim.Width).
		SetString(
			lipgloss.JoinVertical(
				lipgloss.Top,
				title,
				msg,
				b.AddProfile.View(),
				"",
				hints,
			),
		)

	if b.ElemFocus == Focus.TxtInput {
		popup = popup.BorderForeground(color.BoxSelected)
	}

	return popup
}

func editor(b ProfileTab, dims utils.Dimensions) lipgloss.Style {
	b.ViewProfile.SetWidth(dims.Width - 2)
	b.ViewProfile.SetHeight(dims.Height - 3)

	b.ViewProfile.Prompt = ""

	view := style.Box.
		Padding(0, 1).
		Width(dims.Width - 2).
		Height(dims.Height - 2).
		SetString(
			lipgloss.JoinVertical(
				lipgloss.Top,
				style.Title.SetString("Viewer/Editor:").Render(),
				utils.Ifelse(len(utils.Profiles) > 0, b.ViewProfile.View(), "").(string),
			),
		)

	if b.ElemFocus == Focus.TxtEdit {
		view = view.BorderForeground(color.BoxSelected)
	}

	return view
}
