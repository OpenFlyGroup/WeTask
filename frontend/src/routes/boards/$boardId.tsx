import { createFileRoute, redirect } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { authStorage } from '../../api/http'
import { getBoardById } from '../../api/boards'
import { createColumn, deleteColumn, getColumnsByBoard } from '../../api/columns'
import { addComment, createTask, deleteTask, getComments, getTasksByBoard, moveTask } from '../../api/tasks'
import { useEffect, useMemo, useState } from 'react'
import { getMe } from '../../api/users'
import { getSocket } from '../../realtime/socket'
import { RealtimeEvents } from '../../realtime/events'

export const Route = createFileRoute('/boards/$boardId')({
  beforeLoad: () => {
    if (!authStorage.getTokens()) {
      throw redirect({ to: '/auth/login' })
    }
  },
  component: BoardDetailPage,
})

function BoardDetailPage() {
  const { boardId } = Route.useParams()
  const id = Number(boardId)
  const qc = useQueryClient()

  const boardQ = useQuery({ queryKey: ['board', id], queryFn: ({ signal }) => getBoardById(id, signal) })
  const columnsQ = useQuery({
    queryKey: ['columns', id],
    queryFn: ({ signal }) => getColumnsByBoard(id, signal),
  })
  const tasksQ = useQuery({
    queryKey: ['tasks', id],
    queryFn: ({ signal }) => getTasksByBoard(id, signal),
  })
  const meQ = useQuery({ queryKey: ['me'], queryFn: ({ signal }) => getMe(signal) })

  const tasksByColumn = useMemo(() => {
    const map: Record<number, Array<any>> = {}
    for (const t of tasksQ.data ?? []) {
      if (!map[t.columnId]) map[t.columnId] = []
      map[t.columnId].push(t)
    }
    return map
  }, [tasksQ.data])

  const [newColumn, setNewColumn] = useState('')
  const createColMut = useMutation({
    mutationFn: () => createColumn({ name: newColumn, boardId: id }),
    onSuccess: () => {
      setNewColumn('')
      void qc.invalidateQueries({ queryKey: ['columns', id] })
    },
  })
  const deleteColMut = useMutation({
    mutationFn: (colId: number) => deleteColumn(colId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['columns', id] })
      void qc.invalidateQueries({ queryKey: ['tasks', id] })
    },
  })

  // Realtime: join board room and listen for task updates
  useEffect(() => {
    if (!meQ.data) return
    let mounted = true
    void (async () => {
      const socket = await getSocket()
      if (!socket || !mounted) return
      socket.emit('join:board', { boardId: id, userId: meQ.data!.id })
      const invalidateBoardData = () => {
        void qc.invalidateQueries({ queryKey: ['tasks', id] })
        void qc.invalidateQueries({ queryKey: ['columns', id] })
      }
      socket.on(RealtimeEvents.TASK_CREATED, invalidateBoardData)
      socket.on(RealtimeEvents.TASK_UPDATED, invalidateBoardData)
      socket.on(RealtimeEvents.TASK_DELETED, invalidateBoardData)
      return () => {
        socket.emit('leave:board', { boardId: id })
        socket.off(RealtimeEvents.TASK_CREATED, invalidateBoardData)
        socket.off(RealtimeEvents.TASK_UPDATED, invalidateBoardData)
        socket.off(RealtimeEvents.TASK_DELETED, invalidateBoardData)
      }
    })()
    return () => {
      mounted = false
    }
  }, [id, meQ.data, qc])

  const [newTaskTitle, setNewTaskTitle] = useState<Record<number, string>>({})
  const createTaskMut = useMutation({
    mutationFn: (colId: number) =>
      createTask({ title: newTaskTitle[colId], boardId: id, columnId: colId, description: '' }),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['tasks', id] })
    },
  })
  const deleteTaskMut = useMutation({
    mutationFn: (taskId: number) => deleteTask(taskId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['tasks', id] })
    },
  })
  const moveTaskMut = useMutation({
    mutationFn: ({ taskId, toColumnId }: { taskId: number; toColumnId: number }) =>
      moveTask(taskId, { columnId: toColumnId }),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['tasks', id] })
    },
  })

  return (
    <div>
      <h1 className="text-2xl font-semibold mb-4">{boardQ.data?.name ?? 'Board'}</h1>

      <div className="mb-6">
        <form
          className="card bg-base-100 shadow"
          onSubmit={(e) => {
            e.preventDefault()
            createColMut.mutate()
          }}
        >
          <div className="card-body">
            <div className="grid grid-cols-1 md:grid-cols-6 gap-3">
              <label className="form-control md:col-span-5">
                <div className="label">
                  <span className="label-text">New column name</span>
                </div>
                <input
                  className="input input-bordered"
                  placeholder="New column name"
                  value={newColumn}
                  onChange={(e) => setNewColumn(e.target.value)}
                  required
                />
              </label>
              <div className="md:col-span-1 flex items-end">
                <button className="btn btn-primary w-full" disabled={createColMut.isPending}>
                  Add Column
                </button>
              </div>
            </div>
          </div>
        </form>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {(columnsQ.data ?? []).map((col) => (
          <div key={col.id} className="card bg-base-100 shadow">
            <div className="card-body gap-3">
              <div className="flex items-center justify-between">
                <h2 className="card-title">{col.name}</h2>
                <button
                  className="btn btn-error btn-sm"
                  onClick={() => deleteColMut.mutate(col.id)}
                  disabled={deleteColMut.isPending}
                >
                  Delete
                </button>
              </div>

            <form
              className="flex gap-2"
              onSubmit={(e) => {
                e.preventDefault()
                createTaskMut.mutate(col.id)
              }}
            >
              <input
                className="input input-bordered flex-1"
                placeholder="Task title"
                value={newTaskTitle[col.id] ?? ''}
                onChange={(e) => setNewTaskTitle((s) => ({ ...s, [col.id]: e.target.value }))}
                required
              />
              <button className="btn btn-success" disabled={createTaskMut.isPending}>
                Add
              </button>
            </form>

            <ul className="flex flex-col gap-2">
              {(tasksByColumn[col.id] ?? []).map((t) => (
                <li key={t.id} className="card bg-base-200">
                  <div className="card-body p-3">
                  <div className="flex items-center justify-between">
                    <div className="font-medium">{t.title}</div>
                    <div className="flex items-center gap-2">
                      <select
                        className="select select-bordered select-sm"
                        value={t.columnId}
                        onChange={(e) => moveTaskMut.mutate({ taskId: t.id, toColumnId: Number(e.target.value) })}
                      >
                        {(columnsQ.data ?? []).map((c) => (
                          <option key={c.id} value={c.id}>
                            {c.name}
                          </option>
                        ))}
                      </select>
                      <button
                        className="btn btn-error btn-xs"
                        onClick={() => deleteTaskMut.mutate(t.id)}
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                  <TaskComments taskId={t.id} />
                  </div>
                </li>
              ))}
            </ul>
          </div>
          </div>
        ))}
      </div>
    </div>
  )
}

function TaskComments({ taskId }: { taskId: number }) {
  const qc = useQueryClient()
  const commentsQ = useQuery({ queryKey: ['comments', taskId], queryFn: ({ signal }) => getComments(taskId, signal) })
  const [message, setMessage] = useState('')
  const addMut = useMutation({
    mutationFn: () => addComment(taskId, { message }),
    onSuccess: () => {
      setMessage('')
      void qc.invalidateQueries({ queryKey: ['comments', taskId] })
    },
  })
  return (
    <div className="mt-2">
      <ul className="text-sm flex flex-col gap-1 mb-2">
        {(commentsQ.data ?? []).map((c) => (
          <li key={c.id} className="px-2 py-1 rounded bg-base-100">
            {c.message}
          </li>
        ))}
      </ul>
      <form
        className="flex gap-2"
        onSubmit={(e) => {
          e.preventDefault()
          addMut.mutate()
        }}
      >
        <input
          className="input input-bordered input-sm flex-1"
          placeholder="Add comment"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          required
        />
        <button className="btn btn-primary btn-sm" disabled={addMut.isPending}>
          Add
        </button>
      </form>
    </div>
  )
}


