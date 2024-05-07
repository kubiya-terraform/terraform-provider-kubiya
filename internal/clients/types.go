package clients

type Runner struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type user struct {
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

type group struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Roles       []string `json:"roles"`
	System      bool     `json:"system"`
	CreateAt    string   `json:"create_at"`
	Description string   `json:"description"`
}

type secret struct {
	Name        string `json:"secret_name"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
}

type integration struct {
	Id   string `json:"uuid"`
	Name string `json:"name"`
}

type state struct {
	//self         *user
	users        []*user
	agents       []*agent
	groups       []*group
	runners      []*runner
	secrets      []*secret
	webhooks     []*webhook
	integrations []*integration
}
