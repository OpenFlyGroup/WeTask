import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, useEffect } from 'react'
import { UsersService } from '@/shared/api/services/users/users.service'

export const Route = createFileRoute('/_auth/profile/')({
  component: ProfilePage,
})

function ProfilePage() {
  const qc = useQueryClient()
  const { data, isLoading, error } = useQuery({
    queryKey: ['me'],
    queryFn: () => UsersService.getMe(),
  })
  const [name, setName] = useState('')
  useEffect(() => {
    if (data?.name) setName(data.name)
  }, [data])

  const updateMut = useMutation({
    mutationFn: () => UsersService.updateUser(Number(data!.id), { name }),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['me'] })
    },
  })

  if (isLoading) return <div className="loading loading-spinner loading-md" />
  if (error)
    return <div className="alert alert-error">{(error as any).message}</div>

  return (
    <div className="max-w-lg">
      <h1 className="text-2xl font-semibold mb-4">Profile</h1>
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <div>
            <div className="label">
              <span className="label-text">Email</span>
            </div>
            <div className="font-medium">{data?.email}</div>
          </div>
          <form
            className="flex gap-2"
            onSubmit={(e) => {
              e.preventDefault()
              updateMut.mutate()
            }}
          >
            <input
              className="input input-bordered flex-1"
              placeholder="Your name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
            <button className="btn btn-primary" disabled={updateMut.isPending}>
              Save
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}
