import AuthStorage from '@/shared/store/authStore'
import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth')({
  beforeLoad: () => {
    if (typeof window !== 'undefined') {
      if (!AuthStorage.getTokens()) {
        throw redirect({ to: '/signin' })
      }
    }
  },
  component: () => <Outlet />,
})
