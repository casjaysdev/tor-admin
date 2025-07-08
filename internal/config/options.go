// File: internal/config/options.go
// Purpose: Metadata for all known torrc options with UI support

package config

type OptionType string

const (
	TypeString OptionType = "string"
	TypeInt    OptionType = "int"
	TypeBool   OptionType = "bool"
	TypeList   OptionType = "list"
)

type TorOption struct {
	Name        string     `json:"name"`
	Type        OptionType `json:"type"`
	Default     string     `json:"default"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Multiple    bool       `json:"multiple"`
	Deprecated  bool       `json:"deprecated"`
	Advanced    bool       `json:"advanced"`
	InputType   string     `json:"input_type"`
	Choices     []string   `json:"choices,omitempty"`
	Placeholder string     `json:"placeholder,omitempty"`
	Required    bool       `json:"required,omitempty"`
	Resettable  bool       `json:"resettable"`
}

var allTorOptions = []TorOption{
	{
		Name: "SocksPort", Type: TypeInt, Default: "9050", Description: "SOCKS proxy port",
		Category: "Network", InputType: "number", Placeholder: "9050", Required: true, Resettable: true,
	},
	{
		Name: "ControlPort", Type: TypeInt, Default: "9051", Description: "Tor controller port",
		Category: "Network", InputType: "number", Placeholder: "9051", Resettable: true,
	},
	{
		Name: "BandwidthRate", Type: TypeString, Default: "1 MB", Description: "Bandwidth rate limit",
		Category: "Bandwidth", InputType: "text", Placeholder: "5 MB", Resettable: true,
	},
	{
		Name: "BandwidthBurst", Type: TypeString, Default: "1 MB", Description: "Bandwidth burst",
		Category: "Bandwidth", InputType: "text", Placeholder: "10 MB", Resettable: true,
	},
	{
		Name: "HiddenServiceDir", Type: TypeString, Default: "", Description: "Hidden Service directory",
		Category: "Hidden Services", InputType: "text", Placeholder: "/var/lib/tor/hs1", Resettable: true,
	},
	{
		Name: "HiddenServicePort", Type: TypeString, Default: "", Description: "Map virtual port to target address",
		Category: "Hidden Services", InputType: "text", Placeholder: "80 127.0.0.1:8080", Resettable: true,
	},
	{
		Name: "ExitRelay", Type: TypeBool, Default: "0", Description: "Advertise as an exit node",
		Category: "Relay", InputType: "checkbox", Resettable: true,
	},
	{
		Name: "SafeLogging", Type: TypeBool, Default: "1", Description: "Avoid logging sensitive info",
		Category: "Logging", InputType: "checkbox", Resettable: true,
	},
	{
		Name: "Log", Type: TypeString, Default: "notice stdout", Description: "Log level and target",
		Category: "Logging", InputType: "text", Placeholder: "notice stdout", Resettable: true,
	},
}

func GetAllOptions() []TorOption {
	return allTorOptions
}

func GetOptionsByCategory() map[string][]TorOption {
	result := map[string][]TorOption{}
	for _, opt := range allTorOptions {
		result[opt.Category] = append(result[opt.Category], opt)
	}
	return result
}

func GetOption(name string) *TorOption {
	for _, opt := range allTorOptions {
		if opt.Name == name {
			return &opt
		}
	}
	return nil
}
