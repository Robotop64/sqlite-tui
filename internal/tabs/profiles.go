package tabs

import (
	"fmt"
	"path/filepath"
	"strings"

	style "github.com/Robotop64/sqlite-tui/internal/style"
	color "github.com/Robotop64/sqlite-tui/internal/style/color"
	utils "github.com/Robotop64/sqlite-tui/internal/utils"

	bubTxtEdit "github.com/charmbracelet/bubbles/textarea"
	bubTxtIn "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgList "github.com/charmbracelet/lipgloss/list"
	yaml "gopkg.in/yaml.v3"
)

type ProfileTab struct {
	name      string
	ElemFocus UiFocus

	IdxFocus    int
	IdxSelected int

	AddProfile  bubTxtIn.Model
	ViewProfile bubTxtEdit.Model
}

type UiFocus int

const (
	None UiFocus = iota
	TxtInput
	TxtEdit
	ProfileList
)

func (b *ProfileTab) GetName() string {
	return b.name
}

func (b *ProfileTab) Init() tea.Cmd {
	if b.ElemFocus == TxtInput {
		return bubTxtEdit.Blink
	}
	if b.ElemFocus == TxtInput {
		return bubTxtIn.Blink
	}

	return nil
}

func (b *ProfileTab) Setup() Tab {
	b.name = "Profiles"
	b.ElemFocus = ProfileList

	txtinput := bubTxtIn.New()
	txtinput.Placeholder = "..."
	txtinput.Width = 256
	txtinput.CharLimit = 256
	b.AddProfile = txtinput

	txtedit := bubTxtEdit.New()
	b.ViewProfile = txtedit

	return b
}

func (b *ProfileTab) Activate() {

	b.IdxSelected = utils.Data.Profiles.LastUsed

	if len(utils.Profiles) > 0 {

		profile := utils.ActiveProfile()
		data, _ := yaml.Marshal(profile)
		b.ViewProfile.SetValue(string(data))
	}
}

func (b *ProfileTab) View(width, height int) string {
	//=Calculations============================================================
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
	//=========================================================================

	//=Left Column=============================================================
	tab_Box := style.Box.
		Width(explorer_size.Width-2).
		Height(1).
		SetString(
			"⏴",
			lipgloss.PlaceHorizontal(explorer_size.Width-6, lipgloss.Center, b.name),
			"⏵",
		).
		Foreground(color.TextHighlight)

	explorer_height := explorer_size.Height - lipgloss.Height(tab_Box.String())

	list_Box := gen_list(b, utils.Dimensions{Width: explorer_size.Width, Height: explorer_height})

	left_Column := lipgloss.JoinVertical(
		lipgloss.Top,
		tab_Box.Render(),
		list_Box.Render(),
	)
	//=========================================================================

	//=Content Column==========================================================
	right_Column := gen_editor(b, content_size).Render()
	//=========================================================================

	//=Hints===================================================================
	hints := lipgloss.JoinHorizontal(
		lipgloss.Top,
		"quit:\nctrl+c|q", "│\n│",
		"save:\nctrl+s", "│\n│",
		"tabs:\nalt+(</>)", "│\n│",
		"profiles:\n   ↑/↓", "│\n│",
		"add:\n +", "│\n│",
		"remove:\n   -", "│\n│",
	)
	//=========================================================================

	//=Layout==================================================================
	layout := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			left_Column,
			right_Column,
		),
		hints,
	)

	if !(b.ElemFocus == TxtInput) {
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
	//=========================================================================
}

func (b *ProfileTab) Update(msg tea.Msg) (Tab, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "alt+left":
			b.ElemFocus = ProfileList
			b.ViewProfile.Blur()
			return b, nil
		case "alt+right":
			b.ElemFocus = TxtEdit
			return b, nil
		case "ctrl+s":
			data := b.ViewProfile.Value()
			if !(len(data) == 0) {
				profile := utils.ActiveProfile()

				if profile == nil {
					return b, nil
				}

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
		}

		switch b.ElemFocus {
		case None:
			return b, nil
		case ProfileList:
			switch msg.String() {
			case "up":
				b.IdxFocus = max(b.IdxFocus-1, 0)
				return b, nil
			case "down":
				b.IdxFocus = min(b.IdxFocus+1, len(utils.Profiles)-1)
				return b, nil
			case "+":
				b.ElemFocus = TxtInput
				return b, nil
			case "-":
				utils.Profiles = utils.RemoveIdx(utils.Profiles, b.IdxFocus)
				utils.Data.Profiles.Paths = utils.RemoveIdx(utils.Data.Profiles.Paths, b.IdxFocus)

				b.IdxFocus = max(b.IdxFocus-1, 0)
				return b, nil
			case "enter":
				utils.Data.Profiles.LastUsed = b.IdxFocus
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
		case TxtInput:
			b.AddProfile.Focus()
			switch msg.String() {
			case "esc":
				b.ElemFocus = ProfileList
				b.AddProfile.Reset()
				return b, nil
			case "enter":
				raw_path := b.AddProfile.Value()
				if raw_path == "" {
					b.ElemFocus = ProfileList
					return b, nil
				}

				var path string
				if !utils.EndsWith(raw_path, "Profile.yaml") {
					path = filepath.Join(raw_path, "Profile.yaml")
				} else {
					path = raw_path
				}

				fileExists := utils.CheckPath(path)
				var profile *utils.Profile
				var err error
				if !fileExists {
					if profile, err = utils.CreateProfile(path); err != nil {
						b.ElemFocus = ProfileList
						b.AddProfile.Reset()
						return b, nil
					}
				} else {
					if profile, err = utils.LoadProfile(path); err != nil {
						b.ElemFocus = ProfileList
						b.AddProfile.Reset()
						return b, nil
					}
				}
				utils.Profiles = append(utils.Profiles, profile)
				utils.Data.Profiles.Paths = append(utils.Data.Profiles.Paths, path)

				b.ElemFocus = ProfileList
				b.AddProfile.Reset()
				return b, nil
			}
			b.AddProfile, cmd = b.AddProfile.Update(msg)
			return b, cmd
		case TxtEdit:
			if len(utils.Profiles) == 0 {
				return b, nil
			}

			switch msg.String() {
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

func gen_list(b *ProfileTab, dim utils.Dimensions) lipgloss.Style {
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

	if b.ElemFocus == ProfileList {
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

func addProfilePrompt(b *ProfileTab, dim utils.Dimensions) lipgloss.Style {
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

	if b.ElemFocus == TxtInput {
		popup = popup.BorderForeground(color.BoxSelected)
	}

	return popup
}

func gen_editor(b *ProfileTab, dims utils.Dimensions) lipgloss.Style {
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

	if b.ElemFocus == TxtEdit {
		view = view.BorderForeground(color.BoxSelected)
	}

	return view
}
