import {
  HeadContent,
  Outlet,
  Scripts,
  createRootRouteWithContext,
  useRouterState,
} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'
import { motion, AnimatePresence } from 'motion/react'
import { Toaster } from 'react-hot-toast'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import appCss from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'

import Footer from '@/shared/ui/layout/Footer/Footer'
import Header from '@/shared/ui/layout/Header/Header'

interface MyRouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        name: 'viewport',
        content: 'viewport-fit=cover',
      },
      { title: 'WeTask' },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: appCss,
      },
    ],
  }),

  shellComponent: RootDocument,
})

const urlsWithHeader = ['/signin', '/signup', '/']

function RootDocument({ children }: { children: React.ReactNode }) {
  const routerState = useRouterState()
  const currentPath = routerState.location.pathname

  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body className="flex flex-col min-h-screen">
        {urlsWithHeader.includes(currentPath) ? <Header /> : null}

        {children ?? <Outlet />}

        <Footer />

        <TanStackDevtools
          config={{
            position: 'bottom-right',
          }}
          plugins={[
            {
              name: 'Tanstack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
            TanStackQueryDevtools,
          ]}
        />
        <Toaster position="bottom-right" reverseOrder={false} />
        <Scripts />
      </body>
    </html>
  )
}
