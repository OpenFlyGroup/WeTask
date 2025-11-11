import { http } from './http'

export type Team = {
  id: number
  name: string
  ownerId: number
}

export type TeamMember = {
  id: number
  userId: number
  teamId: number
  email?: string
  name?: string
}

export type CreateTeamDto = {
  name: string
}

export type AddMemberDto = {
  userId: number
}

export function getTeams(signal?: AbortSignal) {
  return http<Team[]>('/teams', { auth: true, signal })
}

export function getTeamById(id: number, signal?: AbortSignal) {
  return http<Team>(`/teams/${id}`, { auth: true, signal })
}

export function createTeam(data: CreateTeamDto) {
  return http<Team>('/teams', { method: 'POST', body: data, auth: true })
}

export function addMember(teamId: number, data: AddMemberDto) {
  return http<TeamMember>(`/teams/${teamId}/members`, {
    method: 'POST',
    body: data,
    auth: true,
  })
}

export function removeMember(teamId: number, userId: number) {
  return http<{ success: true }>(`/teams/${teamId}/members/${userId}`, {
    method: 'DELETE',
    auth: true,
  })
}


