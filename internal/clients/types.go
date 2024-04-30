package clients

type Agent struct {
	Name                 string            `json:"name"`
	Uuid                 string            `json:"uuid"`
	Email                string            `json:"email"`
	Image                string            `json:"image"`
	Links                []string          `json:"links"`
	Owners               []string          `json:"owners"`
	Runners              []string          `json:"runners"`
	Secrets              []string          `json:"secrets"`
	Starters             []string          `json:"starters"`
	Metadata             *Metadata         `json:"metadata"`
	LlmModel             string            `json:"llm_model"`
	Description          string            `json:"description"`
	Integrations         []string          `json:"integrations"`
	Organization         string            `json:"organization"`
	AllowedUsers         []string          `json:"allowed_users"`
	AllowedGroups        []string          `json:"allowed_groups"`
	AiInstructions       string            `json:"ai_instructions"`
	EnvironmentVariables map[string]string `json:"environment_variables"`
}

type Runner struct {
	Path string `json:"-"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

type Metadata struct {
	CreatedAt       string `json:"created_at"`
	LastUpdated     string `json:"last_updated"`
	UserCreated     string `json:"user_created"`
	UserLastUpdated string `json:"user_last_updated"`
}

type User struct {
	Roles      []any    `json:"roles"`
	UUID       string   `json:"uuid"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Image      string   `json:"image"`
	Groups     []string `json:"groups"`
	Status     bool     `json:"status"`
	CreateAt   string   `json:"create_at"`
	UserStatus string   `json:"user_status"`
	InviteLink string   `json:"invite_link"`
}

type Group struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Roles       []string `json:"roles"`
	System      bool     `json:"system"`
	CreateAt    string   `json:"create_at"`
	Description string   `json:"description"`
}
