package input

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/style/color"
)

type SectionGroup struct {
	title        string
	sections     []Section
	selectedOpts []string
	sectI        int
}

type Section struct {
	Title   string
	selectI int
	Options []ClientOption
}

type ClientOption struct {
	name     string
	returned string
}

func NewSectionGroup(title string, opts []Section) *SectionGroup {
	sections := make([]Section, 0)
	for _, opt := range opts {
		sections = append(sections, opt)
	}

	return &SectionGroup{
		title:    title,
		sectI:    0,
		sections: sections,
	}
}

func NewSection(name string, opts []ClientOption) (s Section) {
	return Section{
		Title:   name,
		selectI: 0,
		Options: opts,
	}
}

func NewOption(name string, returned string) ClientOption {
	return ClientOption{
		name:     name,
		returned: returned,
	}
}

func (s *SectionGroup) NextSection() []string {
	s.selectedOpts = append(s.selectedOpts, s.sections[s.sectI].CurrentOption().returned)

	s.sectI++
	if s.sectI == len(s.sections) {
		return s.selectedOpts
	}

	return nil
}

func (s *SectionGroup) NextOption() {
	s.sections[s.sectI].NextOption()
}

func (s *SectionGroup) PrevOption() {
	s.sections[s.sectI].PrevOption()
}

func (s *SectionGroup) View() (msg string) {
	msg += s.title + "\n\n"

	for i, sect := range s.sections {
		msg += sect.View(i == s.sectI) + "\n"
	}

	return
}

func (s *Section) CurrentOption() ClientOption {
	return s.Options[s.selectI]
}

func (s *Section) NextOption() {
	s.selectI++
	if s.selectI > len(s.Options)-1 {
		s.selectI = 0
	}
}

func (s *Section) PrevOption() {
	s.selectI--
	if s.selectI < 0 {
		s.selectI = len(s.Options) - 1
	}
}

func (s *Section) View(isActive bool) (msg string) {
	var ind string
	var line string

	msg += s.Title + "\n"

	for i, opt := range s.Options {
		ind = "- "
		if s.selectI == i && isActive {
			ind = "> "
		}

		line = ind + opt.name + "\n"

		msg += line
	}

	if !isActive {
		msg = lipgloss.NewStyle().Foreground(color.InactiveGrey).Render(msg)
	}

	return
}
