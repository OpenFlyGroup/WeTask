import { UrlBuilder } from '@relatecom/utils'
import {
  AddMemberDto,
  CreateTeamDto,
  Team,
  TeamMember,
} from './teams.interface'
import { instance } from '../../instance'

const PATH = '/teams'
const { buildUrl } = new UrlBuilder(PATH)

export const TeamsService = {
  async getTeams() {
    return (await instance.get<Team[]>(buildUrl(''))).data
  },

  async getTeamById(id: number) {
    return (await instance.get<Team>(buildUrl(`/${id}`))).data
  },

  async createTeam(data: CreateTeamDto) {
    return (await instance.post<Team>(buildUrl(''), data)).data
  },

  async addMember(teamId: number, data: AddMemberDto) {
    return (
      await instance.post<TeamMember>(buildUrl(`/${teamId}/members`), data)
    ).data
  },

  async removeMember(teamId: number, userId: number) {
    return instance.delete<{ success: true }>(
      buildUrl(`/${teamId}/members/${userId}`),
    )
  },
}
