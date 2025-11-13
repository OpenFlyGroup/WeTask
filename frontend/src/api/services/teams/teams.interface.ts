export interface Team {
  id: number
  name: string
  ownerId: number
}

export interface TeamMember {
  id: number
  userId: number
  teamId: number
  email?: string
  name?: string
}

export interface CreateTeamDto {
  name: string
}

export interface AddMemberDto {
  userId: number
}
