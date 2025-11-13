export const RealtimeEvents = {
  TASK_CREATED: 'task.created',
  TASK_UPDATED: 'task.updated',
  TASK_DELETED: 'task.deleted',
  BOARD_UPDATED: 'board.updated',
  TEAM_MEMBER_ADDED: 'team.memberAdded',
  TEAM_MEMBER_REMOVED: 'team.memberRemoved',
} as const

export type RealtimeEventKey = keyof typeof RealtimeEvents
