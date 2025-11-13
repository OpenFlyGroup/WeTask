import { createFileRoute, redirect } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { Link } from '@tanstack/react-router'
import { getSocket } from '../../realtime/socket'
import { RealtimeEvents } from '../../realtime/events'
import AuthStorage from '@/store/auth'
import { BoardsService } from '@/api/services/boards/boards.service'
import { TeamsService } from '@/api/services/teams/teams.service'
import { UsersService } from '@/api/services/users/users.service'
import { Board } from '@/api/services/boards/boards.interface'

export const Route = createFileRoute('/boards/')({
  beforeLoad: () => {
    if (!AuthStorage.getTokens()) {
      throw redirect({ to: '/auth/login' })
    }
  },
  component: BoardsPage,
})

function BoardsPage() {
  const qc = useQueryClient()
  const { data, isLoading, error } = useQuery({
    queryKey: ['boards'],
    queryFn: () => BoardsService.getBoards(),
  })
  const teamsQ = useQuery({
    queryKey: ['teams'],
    queryFn: () => TeamsService.getTeams(),
  })
  const meQ = useQuery({
    queryKey: ['me'],
    queryFn: () => UsersService.getMe(),
  })
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')

  const createMut = useMutation({
    mutationFn: () => BoardsService.createBoard({ name, description }),
    onSuccess: () => {
      setName('')
      setDescription('')
      void qc.invalidateQueries({ queryKey: ['boards'] })
    },
  })

  useEffect(() => {
    if (!teamsQ.data || !meQ.data) return
    let mounted = true
    void (async () => {
      const socket = await getSocket()
      if (!socket || !mounted) return
      const teamIds = teamsQ.data!.map((t) => t.id)
      teamIds.forEach((teamId) =>
        socket.emit('join:team', { teamId, userId: meQ.data!.id }),
      )
      const onBoardUpdated = () =>
        void qc.invalidateQueries({ queryKey: ['boards'] })
      socket.on(RealtimeEvents.BOARD_UPDATED, onBoardUpdated)
      return () => {
        teamIds.forEach((teamId) => socket.emit('leave:team', { teamId }))
        socket.off(RealtimeEvents.BOARD_UPDATED, onBoardUpdated)
      }
    })()
    return () => {
      mounted = false
    }
  }, [teamsQ.data, meQ.data, qc])

  const deleteMut = useMutation({
    mutationFn: (id: number) => BoardsService.deleteBoard(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['boards'] })
    },
  })

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-2xl font-semibold mb-4">Boards</h1>
      <form
        className="card bg-base-100 shadow mb-6"
        onSubmit={(e) => {
          e.preventDefault()
          createMut.mutate()
        }}
      >
        <div className="card-body">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <label className="form-control md:col-span-1">
              <div className="label">
                <span className="label-text">Board name</span>
              </div>
              <input
                className="input input-bordered"
                placeholder="Board name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </label>
            <label className="form-control md:col-span-2">
              <div className="label">
                <span className="label-text">Description</span>
              </div>
              <input
                className="input input-bordered"
                placeholder="Description (optional)"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </label>
          </div>
          <div className="card-actions justify-end">
            <button className="btn btn-primary" disabled={createMut.isPending}>
              {createMut.isPending ? 'Creating...' : 'Create'}
            </button>
          </div>
        </div>
      </form>

      {isLoading ? <div>Loading...</div> : null}
      {error ? (
        <div className="alert alert-error">{(error as any).message}</div>
      ) : null}
      <ul className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {data?.map((b: Board) => (
          <li key={b.id} className="card bg-base-100 shadow">
            <div className="card-body">
              <h2 className="card-title">
                <Link to="/boards/$boardId" params={{ boardId: String(b.id) }}>
                  {b.name}
                </Link>
              </h2>
              {b.description ? (
                <div className="text-base-content/70">{b.description}</div>
              ) : null}
              <div className="card-actions justify-end">
                <button
                  className="btn btn-error btn-sm"
                  onClick={() => deleteMut.mutate(b.id)}
                  disabled={deleteMut.isPending}
                >
                  Delete
                </button>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}
