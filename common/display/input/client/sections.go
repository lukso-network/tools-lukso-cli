package input

type SectionGroup struct {
	sections     []Section
	selectedOpts []string
	sectI        int
}

type Section struct {
	Name    string
	selectI int
	Options []ClientOption
}

type ClientOption struct {
	name     string
	returned string
}

func NewSectionGroup(opts []Section) (grp *SectionGroup) {
	for _, opt := range opts {
		grp.sections = append(grp.sections, opt)
	}

	grp.sectI = 0

	return
}

func NewSection(name string, opts []ClientOption) (s Section) {
	return Section{
		Name:    name,
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

func (s SectionGroup) NextSection() []string {
	s.sectI++
	if s.sectI == len(s.sections) {
		return s.selectedOpts
	}

	return nil
}

func (s SectionGroup) View() (msg string) {
	for i, sect := range s.sections {
		msg += sect.View(i == s.sectI)
	}

	return
}

func (s Section) View(isActive bool) (msg string) {
}
