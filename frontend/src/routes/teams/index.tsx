import { createFileRoute, redirect } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { getSocket } from '../../realtime/socket'
import { RealtimeEvents } from '../../realtime/events'
import AuthStorage from '@/store/auth'
import { TeamsService } from '@/api/services/teams/teams.service'
import { UsersService } from '@/api/services/users/users.service'

export const Route = createFileRoute('/teams/')({
  beforeLoad: () => {
    if (!AuthStorage.getTokens()) {
      throw redirect({ to: '/auth/login' })
    }
  },
  component: TeamsPage,
})

function TeamsPage() {
  const qc = useQueryClient()
  const { data, isLoading, error } = useQuery({
    queryKey: ['teams'],
    queryFn: () => TeamsService.getTeams(),
  })
  const meQ = useQuery({
    queryKey: ['me'],
    queryFn: () => UsersService.getMe(),
  })
  const [name, setName] = useState('')
  const createMut = useMutation({
    mutationFn: () => TeamsService.createTeam({ name }),
    onSuccess: () => {
      setName('')
      void qc.invalidateQueries({ queryKey: ['teams'] })
    },
  })

  useEffect(() => {
    if (!data || !meQ.data) return
    let mounted = true
    void (async () => {
      const socket = await getSocket()
      if (!socket || !mounted) return
      const teamIds = data!.map((t) => t.id)
      teamIds.forEach((teamId) =>
        socket.emit('join:team', { teamId, userId: meQ.data!.id }),
      )
      const refreshTeams = () =>
        void qc.invalidateQueries({ queryKey: ['teams'] })
      socket.on(RealtimeEvents.TEAM_MEMBER_ADDED, refreshTeams)
      socket.on(RealtimeEvents.TEAM_MEMBER_REMOVED, refreshTeams)
      socket.on(
        RealtimeEvents.BOARD_UPDATED,
        () => void qc.invalidateQueries({ queryKey: ['boards'] }),
      )
      return () => {
        teamIds.forEach((teamId) => socket.emit('leave:team', { teamId }))
        socket.off(RealtimeEvents.TEAM_MEMBER_ADDED, refreshTeams)
        socket.off(RealtimeEvents.TEAM_MEMBER_REMOVED, refreshTeams)
        socket.off(RealtimeEvents.BOARD_UPDATED)
      }
    })()
    return () => {
      mounted = false
    }
  }, [data, meQ.data, qc])
  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-2xl font-semibold mb-4">Teams</h1>
      <form
        className="card bg-base-100 shadow mb-6"
        onSubmit={(e) => {
          e.preventDefault()
          createMut.mutate()
        }}
      >
        <div className="card-body">
          <div className="grid grid-cols-1 md:grid-cols-6 gap-3">
            <label className="form-control md:col-span-5">
              <div className="label">
                <span className="label-text">Team name</span>
              </div>
              <input
                className="input input-bordered"
                placeholder="Team name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </label>
            <div className="md:col-span-1 flex items-end">
              <button
                className="btn btn-primary w-full"
                disabled={createMut.isPending}
              >
                {createMut.isPending ? 'Creating...' : 'Create'}
              </button>
            </div>
          </div>
        </div>
      </form>
      {isLoading ? <div>Loading...</div> : null}
      {error ? (
        <div className="alert alert-error">{(error as any).message}</div>
      ) : null}
      <ul className="flex flex-col gap-3">
        {(data ?? []).map((t) => (
          <li key={t.id} className="card bg-base-100 shadow">
            <TeamMembers teamId={t.id} />
          </li>
        ))}
      </ul>
    </div>
  )
}

function TeamMembers({ teamId }: { teamId: number }) {
  const qc = useQueryClient()
  const teamQ = useQuery({
    queryKey: ['team', teamId],
    queryFn: () => TeamsService.getTeamById(teamId),
  })
  const [userId, setUserId] = useState('')
  const addMut = useMutation({
    mutationFn: () =>
      TeamsService.addMember(teamId, { userId: Number(userId) }),
    onSuccess: () => {
      setUserId('')
      void qc.invalidateQueries({ queryKey: ['team', teamId] })
    },
  })
  const removeMut = useMutation({
    mutationFn: (uid: number) => TeamsService.removeMember(teamId, uid),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['team', teamId] })
    },
  })
  return (
    <div className="card-body">
      <div className="flex items-center justify-between mb-2">
        <div className="card-title">
          {teamQ.data?.name ?? `Team #${teamId}`}
        </div>
      </div>
      <form
        className="flex gap-2 mb-3"
        onSubmit={(e) => {
          e.preventDefault()
          addMut.mutate()
        }}
      >
        <input
          className="input input-bordered"
          placeholder="User ID to add"
          type="number"
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
          required
        />
        <button className="btn btn-success" disabled={addMut.isPending}>
          Add Member
        </button>
      </form>
      <div className="text-sm">
        <div className="mb-1 font-medium">Members:</div>
        <ul className="flex flex-col gap-1">
          {(teamQ.data as any)?.members?.map?.((m: any) => (
            <li
              key={m.userId}
              className="flex items-center justify-between bg-base-200 px-2 py-1 rounded"
            >
              <span>
                #{m.userId} {m.name ?? ''} {m.email ?? ''}
              </span>
              <button
                className="btn btn-error btn-xs"
                onClick={() => removeMut.mutate(m.userId)}
                disabled={removeMut.isPending}
              >
                Remove
              </button>
            </li>
          )) ?? <li className="text-gray-400">No members</li>}
        </ul>
      </div>
    </div>
  )
}
