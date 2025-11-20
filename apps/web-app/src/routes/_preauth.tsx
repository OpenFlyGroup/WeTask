import { createFileRoute, Outlet, useRouterState } from '@tanstack/react-router'
import { motion, AnimatePresence } from 'motion/react'

export const Route = createFileRoute('/_preauth')({
  component: RouteComponent,
})

function RouteComponent() {
  const routerState = useRouterState()
  const currentPath = routerState.location.pathname
  return (
    <AnimatePresence mode="wait">
      <motion.main
        key={currentPath}
        initial={{ opacity: 0, filter: 'blur(8px)', y: 10 }}
        animate={{ opacity: 1, filter: 'blur(0px)', y: 0 }}
        exit={{ opacity: 0, filter: 'blur(8px)', y: -10 }}
        transition={{ duration: 0.2, ease: 'easeInOut' }}
        className="flex-1 mx-auto w-full"
      >
        <Outlet />
      </motion.main>
    </AnimatePresence>
  )
}
