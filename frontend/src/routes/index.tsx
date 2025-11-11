import { createFileRoute, Link, redirect } from '@tanstack/react-router'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getBoards } from '../api/boards'
import { getTeams } from '../api/teams'
import { authStorage } from '../api/http'
import { useEffect } from 'react'
import { getMe } from '../api/users'
import { getSocket } from '../realtime/socket'
import { RealtimeEvents } from '../realtime/events'

export const Route = createFileRoute('/')({
  beforeLoad: () => {
    if (!authStorage.getTokens()) {
      throw redirect({ to: '/auth/login' })
    }
  },
  component: Dashboard,
})

function Dashboard() {
  const qc = useQueryClient()
  const boardsQ = useQuery({ queryKey: ['boards'], queryFn: ({ signal }) => getBoards(signal) })
  const teamsQ = useQuery({ queryKey: ['teams'], queryFn: ({ signal }) => getTeams(signal) })
  const meQ = useQuery({ queryKey: ['me'], queryFn: ({ signal }) => getMe(signal) })

  // Join all team rooms and update boards/teams on board updates or team membership changes
  useEffect(() => {
    if (!teamsQ.data || !meQ.data) return
    let mounted = true
    void (async () => {
      const socket = await getSocket()
      if (!socket || !mounted) return
      const teamIds = teamsQ.data!.map((t) => t.id)
      teamIds.forEach((teamId) => socket.emit('join:team', { teamId, userId: meQ.data!.id }))
      const refreshBoards = () => void qc.invalidateQueries({ queryKey: ['boards'] })
      const refreshTeams = () => void qc.invalidateQueries({ queryKey: ['teams'] })
      socket.on(RealtimeEvents.BOARD_UPDATED, refreshBoards)
      socket.on(RealtimeEvents.TEAM_MEMBER_ADDED, refreshTeams)
      socket.on(RealtimeEvents.TEAM_MEMBER_REMOVED, refreshTeams)
      return () => {
        teamIds.forEach((teamId) => socket.emit('leave:team', { teamId }))
        socket.off(RealtimeEvents.BOARD_UPDATED, refreshBoards)
        socket.off(RealtimeEvents.TEAM_MEMBER_ADDED, refreshTeams)
        socket.off(RealtimeEvents.TEAM_MEMBER_REMOVED, refreshTeams)
      }
    })()
    return () => {
      mounted = false
    }
  }, [teamsQ.data, meQ.data, qc])
  return (
    <div className="max-w-5xl mx-auto">
      <h1 className="text-2xl font-semibold mb-6">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <section className="card bg-base-100 shadow">
          <div className="card-body">
            <div className="flex items-center justify-between mb-2">
              <h2 className="card-title">Boards</h2>
              <Link to="/boards" className="btn btn-ghost btn-sm">
              View all
            </Link>
          </div>
          <ul className="menu bg-base-200 rounded-box">
            {(boardsQ.data ?? []).slice(0, 5).map((b) => (
              <li key={b.id}>
                <Link to="/boards/$boardId" params={{ boardId: String(b.id) }} className="justify-start">
                  {b.name}
                </Link>
              </li>
            ))}
            {!boardsQ.data?.length ? <li className="disabled">No boards yet</li> : null}
          </ul>
        </div>
      </section>
        <section className="card bg-base-100 shadow">
          <div className="card-body">
            <div className="flex items-center justify-between mb-2">
              <h2 className="card-title">Teams</h2>
              <Link to="/teams" className="btn btn-ghost btn-sm">
              View all
            </Link>
            </div>
          <ul className="menu bg-base-200 rounded-box">
            {(teamsQ.data ?? []).slice(0, 5).map((t) => (
              <li key={t.id}>
                {t.name}
              </li>
          ))}
            {!teamsQ.data?.length ? <li className="disabled">No teams yet</li> : null}
          </ul>
        </div>
      </section>
      </div>
    </div>
  )
}
