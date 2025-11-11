import { createFileRoute, redirect } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { authStorage } from '../api/http'
import { getMe, updateUser } from '../api/users'
import { useState, useEffect } from 'react'

export const Route = createFileRoute('/profile')({
  beforeLoad: () => {
    if (!authStorage.getTokens()) {
      throw redirect({ to: '/auth/login' })
    }
  },
  component: ProfilePage,
})

function ProfilePage() {
  const qc = useQueryClient()
  const { data, isLoading, error } = useQuery({ queryKey: ['me'], queryFn: ({ signal }) => getMe(signal) })
  const [name, setName] = useState('')
  useEffect(() => {
    if (data?.name) setName(data.name)
  }, [data])

  const updateMut = useMutation({
    mutationFn: () => updateUser(data!.id, { name }),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ['me'] })
    },
  })

  if (isLoading) return <div className="loading loading-spinner loading-md" />
  if (error) return <div className="alert alert-error">{(error as any).message}</div>

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


