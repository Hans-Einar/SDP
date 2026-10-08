package values

import ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"

func OptionSet(mode string) []ui.ChoiceOption {
	options := []ui.ChoiceOption{{ID: "alpha", Label: "Repeated", Enabled: true}, {ID: "beta", Label: "Repeated", Enabled: true}, {ID: "blocked", Label: "Disabled", Enabled: false}}
	switch mode {
	case "removed":
		return []ui.ChoiceOption{options[1], options[2], {ID: "gamma", Label: "Gamma", Enabled: true}}
	case "disabled":
		options[0].Enabled = false
	case "reordered":
		options[0], options[1] = options[1], options[0]
	}
	return options
}
func (f *Fixture) Choices() map[string][]ui.ChoiceOption {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := map[string][]ui.ChoiceOption{}
	for path, options := range f.choices {
		out[path] = append([]ui.ChoiceOption{}, options...)
	}
	return out
}
func (f *Fixture) storeChoices(path string, options []ui.ChoiceOption) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.choices[path] = append([]ui.ChoiceOption{}, options...)
}
