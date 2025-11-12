package common

// RabbitMQ patterns for RPC communication
const (
	// Auth patterns
	AuthRegister = "auth.register"
	AuthLogin    = "auth.login"
	AuthRefresh  = "auth.refresh"
	AuthValidate = "auth.validate"

	// Users patterns
	UsersGetByID    = "users.getById"
	UsersGetByEmail = "users.getByEmail"
	UsersUpdate     = "users.update"
	UsersGetMe      = "users.getMe"

	// Teams patterns
	TeamsCreate      = "teams.create"
	TeamsGetAll      = "teams.getAll"
	TeamsGetByID     = "teams.getById"
	TeamsAddMember   = "teams.addMember"
	TeamsRemoveMember = "teams.removeMember"
	TeamsGetUserTeams = "teams.getUserTeams"

	// Boards patterns
	BoardsCreate   = "boards.create"
	BoardsGetAll   = "boards.getAll"
	BoardsGetByID  = "boards.getById"
	BoardsUpdate   = "boards.update"
	BoardsDelete   = "boards.delete"
	BoardsGetByTeam = "boards.getByTeam"

	// Columns patterns
	ColumnsCreate    = "columns.create"
	ColumnsGetByBoard = "columns.getByBoard"
	ColumnsUpdate    = "columns.update"
	ColumnsDelete    = "columns.delete"

	// Tasks patterns
	TasksCreate      = "tasks.create"
	TasksGetByID     = "tasks.getById"
	TasksGetByBoard  = "tasks.getByBoard"
	TasksUpdate      = "tasks.update"
	TasksDelete      = "tasks.delete"
	TasksMove        = "tasks.move"
	TasksAddComment  = "tasks.addComment"
	TasksGetComments = "tasks.getComments"
)

// RabbitMQ events for pub/sub
const (
	TaskCreated      = "task.created"
	TaskUpdated      = "task.updated"
	TaskDeleted      = "task.deleted"
	BoardUpdated     = "board.updated"
	TeamMemberAdded  = "team.memberAdded"
	TeamMemberRemoved = "team.memberRemoved"
)

// RPCResponse represents a standard RPC response
type RPCResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
}

