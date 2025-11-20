import { Link, useRouterState } from '@tanstack/react-router'
import {
  DoorOpen,
  LayoutDashboard,
  PanelRightClose,
  Table,
  User,
  Users,
} from 'lucide-react'
import Breadcrumbs from '../Breadcrumbs/Breadcrumbs'
import { IBreadcrumbs } from '@/shared/types/ui/layout/breadcrumbs.interface'
import { motion, AnimatePresence } from 'motion/react'
import OFLogo from '@/assets/of_logo.svg'
import { disconnectSocket } from '@/shared/api/realtime/socket'
import AuthStorage from '@/shared/store/authStore'

const navItems = [
  {
    id: '0',
    icon: <LayoutDashboard className="size-4" />,
    to: '/dashboard',
    label: 'Dashboard',
  },
  {
    id: '1',
    icon: <Table className="size-4" />,
    to: '/boards',
    label: 'Boards',
  },
  { id: '2', icon: <Users className="size-4" />, to: '/teams', label: 'Teams' },
]

const SideBar = ({ children }: { children: React.ReactNode }) => {
  const routerState = useRouterState()
  const currentPath = routerState.location.pathname
  const breadcrumbs: IBreadcrumbs[] = currentPath
    .split('/')
    .filter(Boolean)
    .map((segment, index, arr) => {
      const href = '/' + arr.slice(0, index + 1).join('/')
      return {
        id: String(index),
        title: segment.charAt(0).toUpperCase() + segment.slice(1),
        href,
      }
    })
  return (
    <>
      <div className="drawer lg:drawer-open">
        <input id="SideBar" type="checkbox" className="drawer-toggle" />
        <div className="drawer-content">
          <nav className="navbar w-full bg-base-300">
            <div className="navbar-start">
              <label
                htmlFor="SideBar"
                aria-label="open sidebar"
                className="btn btn-square btn-ghost"
              >
                <PanelRightClose className="size-4" />
              </label>
              <div className="px-4">
                <Breadcrumbs breadcrumbs={breadcrumbs} />
              </div>
            </div>
            <div className="navbar-end gap-2">
              <motion.div whileHover={{ scale: 1.05 }}>
                <Link className="btn p-2 btn-sm btn-circle" to="/profile">
                  <User />
                </Link>
              </motion.div>
              <motion.div whileHover={{ scale: 1.05 }}>
                <Link
                  to="/"
                  onClick={() => {
                    disconnectSocket()
                    AuthStorage.clearTokens()
                  }}
                  className="btn p-2 btn-sm btn-error btn-circle"
                >
                  <DoorOpen />
                </Link>
              </motion.div>
            </div>
          </nav>
          <AnimatePresence mode="wait">
            <motion.main
              key={routerState.location.pathname}
              initial={{ opacity: 0, filter: 'blur(8px)', y: 10 }}
              animate={{ opacity: 1, filter: 'blur(0px)', y: 0 }}
              exit={{ opacity: 0, filter: 'blur(8px)', y: -10 }}
              transition={{ duration: 0.2, ease: 'easeInOut' }}
              className="p-4"
            >
              {children}
            </motion.main>
          </AnimatePresence>
        </div>

        <div className="drawer-side is-drawer-close:overflow-visible">
          <label
            htmlFor="SideBar"
            aria-label="close sidebar"
            className="drawer-overlay"
          ></label>
          <div className="flex min-h-full flex-col items-start bg-base-200 is-drawer-close:w-14 is-drawer-open:w-64">
            {/* Sidebar content here */}
            <Link
              to="/"
              className="flex h-16 w-full items-center justify-center bg-base-100 border-b border-base-300"
            >
              <img className="size-8" src={OFLogo} alt="logo" />
            </Link>
            <ul className="menu gap-2 w-full grow">
              {/* List item */}
              {navItems.map(({ id, icon, to, label }) => {
                const isActive = to === currentPath
                return (
                  <li>
                    <Link
                      key={id}
                      to={to}
                      className={
                        'is-drawer-close:tooltip is-drawer-close:tooltip-right' +
                        (isActive
                          ? ' bg-base-300 text-primary font-semibold'
                          : '')
                      }
                      data-tip={label}
                    >
                      {/* Home icon */}
                      {icon}
                      <span className="is-drawer-close:hidden">{label}</span>
                    </Link>
                  </li>
                )
              })}
            </ul>
          </div>
        </div>
      </div>
    </>
  )
}

export default SideBar
