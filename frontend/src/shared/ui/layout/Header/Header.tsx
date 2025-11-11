import { Link, useRouterState } from '@tanstack/react-router'
import { motion } from 'motion/react'
import logo from 'src/assets/logo.svg'
import { authStorage } from '@/api/http'
import { disconnectSocket } from '@/realtime/socket'
import clsx from 'clsx'

const Header = () => {
  const isAuthed = Boolean(authStorage.getTokens())
  const routerState = useRouterState()
  const currentPath = routerState.location.pathname

  const navItems = [
    { to: '/boards', label: 'Boards' },
    { to: '/teams', label: 'Teams' },
    { to: '/profile', label: 'Profile' },
  ]

  return (
    <motion.nav
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ type: 'spring', stiffness: 100, damping: 15 }}
      className="navbar p-4 bg-base-100 shadow"
    >
      <div className="flex-1">
        <div className="flex items-end gap-2">
          <Link to="/">
            <motion.img
              className="h-8"
              src={logo}
              alt="logo"
              whileHover={{ scale: 1.05 }}
              transition={{ duration: 0.2 }}
            />
          </Link>
          <a
            href="https://openflygroup.github.io/enterprise_landing/"
            target="_blank"
            rel="noreferrer"
            className="text-[0.5rem]"
          >
            V.DEV
          </a>
        </div>

        {isAuthed && (
          <div className="hidden md:flex gap-2 ml-2">
            {navItems.map(({ to, label }) => {
              const isActive = currentPath.startsWith(to)
              return (
                <motion.div key={to} whileHover={{ scale: 1.05 }}>
                  <Link
                    to={to}
                    className={clsx(
                      'btn btn-ghost btn-sm transition-colors',
                      isActive && 'btn-active text-primary'
                    )}
                  >
                    {label}
                  </Link>
                </motion.div>
              )
            })}
          </div>
        )}
      </div>

      <div className="flex-none">
        {isAuthed ? (
          <motion.div whileHover={{ scale: 1.05 }}>
            <Link
              to="/auth/login"
              onClick={() => {
                disconnectSocket()
                authStorage.clear()
              }}
              className="btn btn-sm"
            >
              Logout
            </Link>
          </motion.div>
        ) : (
          <div className="flex gap-2">
            <motion.div whileHover={{ scale: 1.05 }}>
              <Link to="/auth/login" className="btn btn-outline btn-sm">
                Login
              </Link>
            </motion.div>
            <motion.div whileHover={{ scale: 1.05 }}>
              <Link to="/auth/register" className="btn btn-primary btn-sm">
                Register
              </Link>
            </motion.div>
          </div>
        )}
      </div>
    </motion.nav>
  )
}

export default Header