import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
  Outlet,
  useRouterState,
} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'
import { motion, AnimatePresence } from 'motion/react'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'
import appCss from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import Footer from '@/shared/ui/layout/Footer/Footer'
import Header from '@/shared/ui/layout/Header/Header'

interface MyRouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { title: 'WeTask' },
    ],
    links: [{ rel: 'stylesheet', href: appCss }],
  }),
  shellComponent: RootDocument,
})

function RootDocument({ children }: { children: React.ReactNode }) {
  const routerState = useRouterState()

  const canUseLocalStorage = (() => {
    try {
      if (typeof window === 'undefined') return false
      const k = '__wetask_ls_check__'
      window.localStorage?.setItem(k, '1')
      window.localStorage?.removeItem(k)
      return true
    } catch {
      return false
    }
  })()

  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body className="min-h-dvh bg-base-200 flex flex-col overflow-x-hidden">
        {/* Header component */}
        <Header />

        {/* Page transitions */}
        <AnimatePresence mode="wait">
          <motion.main
            key={routerState.location.pathname}
            initial={{ opacity: 0, filter: 'blur(8px)', y: 10 }}
            animate={{ opacity: 1, filter: 'blur(0px)', y: 0 }}
            exit={{ opacity: 0, filter: 'blur(8px)', y: -10 }}
            transition={{ duration: 0.35, ease: 'easeInOut' }}
            className="flex-1 p-4 max-w-7xl mx-auto w-full"
          >
            {children ?? <Outlet />}
          </motion.main>
        </AnimatePresence>

        <Footer />

        {canUseLocalStorage && (
          <TanStackDevtools
            config={{ position: 'bottom-right' }}
            plugins={[
              {
                name: 'Tanstack Router',
                render: <TanStackRouterDevtoolsPanel />,
              },
              TanStackQueryDevtools,
            ]}
          />
        )}

        <Toaster position="bottom-right" reverseOrder={false} />
        <Scripts />
      </body>
    </html>
  )
}
