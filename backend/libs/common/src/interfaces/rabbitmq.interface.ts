export interface RabbitMQPattern {
  cmd: string;
}

export interface RPCResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  statusCode?: number;
}

export enum RabbitMQPatterns {
  // Auth patterns
  AUTH_REGISTER = 'auth.register',
  AUTH_LOGIN = 'auth.login',
  AUTH_REFRESH = 'auth.refresh',
  AUTH_VALIDATE = 'auth.validate',

  // Users patterns
  USERS_GET_BY_ID = 'users.getById',
  USERS_GET_BY_EMAIL = 'users.getByEmail',
  USERS_UPDATE = 'users.update',
  USERS_GET_ME = 'users.getMe',

  // Teams patterns
  TEAMS_CREATE = 'teams.create',
  TEAMS_GET_ALL = 'teams.getAll',
  TEAMS_GET_BY_ID = 'teams.getById',
  TEAMS_ADD_MEMBER = 'teams.addMember',
  TEAMS_REMOVE_MEMBER = 'teams.removeMember',
  TEAMS_GET_USER_TEAMS = 'teams.getUserTeams',

  // Boards patterns
  BOARDS_CREATE = 'boards.create',
  BOARDS_GET_ALL = 'boards.getAll',
  BOARDS_GET_BY_ID = 'boards.getById',
  BOARDS_UPDATE = 'boards.update',
  BOARDS_DELETE = 'boards.delete',
  BOARDS_GET_BY_TEAM = 'boards.getByTeam',

  // Columns patterns
  COLUMNS_CREATE = 'columns.create',
  COLUMNS_GET_BY_BOARD = 'columns.getByBoard',
  COLUMNS_UPDATE = 'columns.update',
  COLUMNS_DELETE = 'columns.delete',

  // Tasks patterns
  TASKS_CREATE = 'tasks.create',
  TASKS_GET_BY_ID = 'tasks.getById',
  TASKS_GET_BY_BOARD = 'tasks.getByBoard',
  TASKS_UPDATE = 'tasks.update',
  TASKS_DELETE = 'tasks.delete',
  TASKS_MOVE = 'tasks.move',
  TASKS_ADD_COMMENT = 'tasks.addComment',
  TASKS_GET_COMMENTS = 'tasks.getComments',
}

export enum RabbitMQEvents {
  TASK_CREATED = 'task.created',
  TASK_UPDATED = 'task.updated',
  TASK_DELETED = 'task.deleted',
  BOARD_UPDATED = 'board.updated',
  TEAM_MEMBER_ADDED = 'team.memberAdded',
  TEAM_MEMBER_REMOVED = 'team.memberRemoved',
}
