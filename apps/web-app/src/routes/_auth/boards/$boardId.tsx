import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useEffect, useMemo, useState } from 'react'
import { BoardsService } from '@/shared/api/services/boards/boards.service'
import { ColumnsService } from '@/shared/api/services/columns/columns.service'
import { TasksService } from '@/shared/api/services/tasks/tasks.service'
import { UsersService } from '@/shared/api/services/users/users.service'
import { Task } from '@/shared/api/services/tasks/tasks.interface'
import { getSocket } from '@/shared/api/realtime/socket'
import { RealtimeEvents } from '@/shared/api/realtime/events'

export const Route = createFileRoute('/_auth/boards/$boardId')({
  component: BoardDetailPage,
})

function BoardDetailPage() {
  const { boardId } = Route.useParams()
  const id = Number(boardId)
  const qc = useQueryClient()

  const boardQ = useQuery({
    queryKey: ['board', id],
    queryFn: () => BoardsService.getBoardById(id),
  })

  const columnsQ = useQuery({
    queryKey: ['columns', id],
    queryFn: () => ColumnsService.getColumnsByBoard(id),
  })

  const tasksQ = useQuery({
    queryKey: ['tasks', id],
    queryFn: () => TasksService.getTasksByBoard(id),
  })

  const meQ = useQuery({
    queryKey: ['me'],
    queryFn: () => UsersService.getMe(),
  })

  const tasksByColumn = useMemo(() => {
    const map: Record<number, Task[]> = {}
    for (const task of tasksQ.data ?? []) {
      if (!map[task.columnId]) map[task.columnId] = []
      map[task.columnId].push(task)
    }
    return map
  }, [tasksQ.data])

  const [newColumn, setNewColumn] = useState<string>('')
  const createColMut = useMutation({
    mutationFn: () =>
      ColumnsService.createColumn({ name: newColumn, boardId: id }),
    onSuccess: () => {
      setNewColumn('')
      void qc.invalidateQueries({ queryKey: ['columns', id] })
    },
  })

  const deleteColMut = useMutation({
    mutationFn: (colId: number) => ColumnsService.deleteColumn(colId),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['columns', id] })
      void qc.invalidateQueries({ queryKey: ['tasks', id] })
    },
  })

  useEffect(() => {
    if (!meQ.data) return

    let mounted = true
    let cleanup: (() => void) | undefined

    void (async () => {
      const socket = await getSocket()
      if (!socket || !mounted) return

      socket.emit('join:board', { boardId: id, userId: meQ.data.id })

      const invalidateBoardData = () => {
        void qc.invalidateQueries({ queryKey: ['tasks', id] })
        void qc.invalidateQueries({ queryKey: ['columns', id] })
      }

      socket.on(RealtimeEvents.TASK_CREATED, invalidateBoardData)
      socket.on(RealtimeEvents.TASK_UPDATED, invalidateBoardData)
      socket.on(RealtimeEvents.TASK_DELETED, invalidateBoardData)

      cleanup = () => {
        socket.emit('leave:board', { boardId: id })
        socket.off(RealtimeEvents.TASK_CREATED, invalidateBoardData)
        socket.off(RealtimeEvents.TASK_UPDATED, invalidateBoardData)
        socket.off(RealtimeEvents.TASK_DELETED, invalidateBoardData)
      }
    })()

    return () => {
      mounted = false
      cleanup?.()
    }
  }, [id, meQ.data?.id, qc])

  const [newTaskTitle, setNewTaskTitle] = useState<Record<number, string>>({})

  const createTaskMut = useMutation({
    mutationFn: (colId: number) => {
      const title = newTaskTitle[colId]
      if (!title?.trim()) throw new Error('Title is required')
      return TasksService.createTask({
        title,
        boardId: id,
        columnId: colId,
        description: '',
      })
    },
    onSuccess: (_, colId) => {
      setNewTaskTitle((s) => ({ ...s, [colId]: '' }))
      void qc.invalidateQueries({ queryKey: ['tasks', id] })
    },
  })

  const deleteTaskMut = useMutation({
    mutationFn: (taskId: number) => TasksService.deleteTask(taskId),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['tasks', id] }),
  })

  const moveTaskMut = useMutation({
    mutationFn: ({
      taskId,
      toColumnId,
    }: {
      taskId: number
      toColumnId: number
    }) => TasksService.moveTask(taskId, { columnId: toColumnId }),
    onSuccess: () => void qc.invalidateQueries({ queryKey: ['tasks', id] }),
  })

  return (
    <div className="p-4">
      <h1 className="text-2xl font-semibold mb-4">
        {boardQ.data?.name ?? 'Loading...'}
      </h1>

      <div className="mb-6">
        <form
          className="card bg-base-100 shadow"
          onSubmit={(e) => {
            e.preventDefault()
            if (newColumn.trim()) createColMut.mutate()
          }}
        >
          <div className="card-body">
            <div className="grid grid-cols-1 md:grid-cols-6 gap-3">
              <fieldset className="fieldset md:col-span-5">
                <legend className="fieldset-legend">New column name</legend>
                <input
                  className="input input-bordered"
                  placeholder="Enter column name"
                  value={newColumn}
                  onChange={(e) => setNewColumn(e.target.value)}
                  required
                />
              </fieldset>
              <div className="md:col-span-1 flex items-end">
                <button
                  className="btn btn-primary w-full"
                  type="submit"
                  disabled={createColMut.isPending || !newColumn.trim()}
                >
                  {createColMut.isPending ? 'Adding...' : 'Add Column'}
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
                  const title = newTaskTitle[col.id]?.trim()
                  if (title) createTaskMut.mutate(col.id)
                }}
              >
                <input
                  className="input input-bordered flex-1"
                  placeholder="Task title"
                  value={newTaskTitle[col.id] ?? ''}
                  onChange={(e) =>
                    setNewTaskTitle((s) => ({ ...s, [col.id]: e.target.value }))
                  }
                  required
                />
                <button
                  className="btn btn-success"
                  type="submit"
                  disabled={createTaskMut.isPending}
                >
                  Add
                </button>
              </form>

              <ul className="flex flex-col gap-2">
                {(tasksByColumn[col.id] ?? []).map((task) => (
                  <li key={task.id} className="card bg-base-200">
                    <div className="card-body p-3">
                      <div className="flex items-center justify-between">
                        <div className="font-medium">{task.title}</div>
                        <div className="flex items-center gap-2">
                          <select
                            className="select select-bordered select-sm"
                            value={task.columnId}
                            onChange={(e) => {
                              const toColumnId = parseInt(e.target.value, 10)
                              if (!isNaN(toColumnId)) {
                                moveTaskMut.mutate({
                                  taskId: task.id,
                                  toColumnId,
                                })
                              }
                            }}
                          >
                            {(columnsQ.data ?? []).map((c) => (
                              <option key={c.id} value={c.id}>
                                {c.name}
                              </option>
                            ))}
                          </select>
                          <button
                            className="btn btn-error btn-xs"
                            onClick={() => deleteTaskMut.mutate(task.id)}
                            disabled={deleteTaskMut.isPending}
                          >
                            X
                          </button>
                        </div>
                      </div>
                      <TaskComments taskId={task.id} />
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
  const commentsQ = useQuery({
    queryKey: ['comments', taskId],
    queryFn: () => TasksService.getComments(taskId),
  })

  const [message, setMessage] = useState('')

  const addMut = useMutation({
    mutationFn: () => {
      if (!message.trim()) throw new Error('Message is required')
      return TasksService.addComment(taskId, { message })
    },
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
          if (message.trim()) addMut.mutate()
        }}
      >
        <input
          className="input input-bordered input-sm flex-1"
          placeholder="Add comment"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          required
        />
        <button
          className="btn btn-primary btn-sm"
          type="submit"
          disabled={addMut.isPending || !message.trim()}
        >
          {addMut.isPending ? '...' : 'Add'}
        </button>
      </form>
    </div>
  )
}
