import { Link, useRouterState } from '@tanstack/react-router'
import { motion, AnimatePresence } from 'motion/react'
import logo from 'src/assets/logo.svg'
import clsx from 'clsx'
import AuthStorage from '@/shared/store/authStore'
import { disconnectSocket } from '@/shared/api/realtime/socket'
import { useState } from 'react'
import { DoorOpen, Menu, X } from 'lucide-react'

const Header = () => {
  const isAuthed = AuthStorage.isAuthenticated()
  const routerState = useRouterState()
  const currentPath = routerState.location.pathname
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)

  const navItems = [
    { to: '/dashboard', label: 'Dashboard' },
    { to: '/profile', label: 'Profile' },
  ]

  return (
    <motion.nav
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ type: 'spring', stiffness: 100, damping: 15 }}
      className="navbar gap-4 p-4 bg-base-100 shadow-lg z-50 relative"
    >
      <div className="flex-1">
        <div className="flex items-end gap-2">
          <Link to="/">
            <motion.img
              className="h-8 w-auto"
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
            className="text-[0.5rem] opacity-60 hover:opacity-100 transition-opacity"
          >
            V.DEV
          </a>
        </div>
      </div>

      {isAuthed && (
        <div className="hidden lg:flex gap-2">
          {navItems.map(({ to, label }) => {
            const isActive = currentPath.startsWith(to)
            return (
              <motion.div key={to} whileHover={{ scale: 1.05 }}>
                <Link
                  to={to}
                  className={clsx(
                    'btn btn-ghost btn-sm rounded-btn transition-all',
                    isActive && 'btn-active text-primary font-medium',
                  )}
                >
                  {label}
                </Link>
              </motion.div>
            )
          })}
        </div>
      )}

      <div className="hidden md:flex flex-none gap-3">
        {isAuthed ? (
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
        ) : (
          <>
            <motion.div whileHover={{ scale: 1.05 }}>
              <Link to="/signin" className="btn btn-sm btn-outline">
                Login
              </Link>
            </motion.div>
            <motion.div whileHover={{ scale: 1.05 }}>
              <Link to="/signup" className="btn btn-sm btn-primary">
                Register
              </Link>
            </motion.div>
          </>
        )}
      </div>

      <div className="md:hidden flex-none">
        <button
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          className="btn btn-ghost btn-square btn-sm"
          aria-label="Toggle menu"
        >
          {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
      </div>

      <AnimatePresence>
        {mobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.2 }}
            className="absolute top-full left-0 right-0 bg-base-100 shadow-2xl border-t border-base-300 md:hidden z-50"
          >
            <div className="flex flex-col p-4 gap-3">
              {isAuthed ? (
                <>
                  {navItems.map(({ to, label }) => {
                    const isActive = currentPath.startsWith(to)
                    return (
                      <Link
                        key={to}
                        to={to}
                        onClick={() => setMobileMenuOpen(false)}
                        className={clsx(
                          'btn btn-ghost justify-start text-lg',
                          isActive && 'text-primary font-semibold bg-base-200',
                        )}
                      >
                        {label}
                      </Link>
                    )
                  })}

                  <div className="border-t border-base-300 my-2" />

                  <Link
                    to="/signin"
                    onClick={() => {
                      disconnectSocket()
                      AuthStorage.clearTokens()
                      setMobileMenuOpen(false)
                    }}
                    className="btn btn-outline"
                  >
                    Logout
                  </Link>
                </>
              ) : (
                <>
                  <Link
                    to="/signin"
                    onClick={() => setMobileMenuOpen(false)}
                    className="btn btn-outline w-full"
                  >
                    Login
                  </Link>
                  <Link
                    to="/signup"
                    onClick={() => setMobileMenuOpen(false)}
                    className="btn btn-primary w-full"
                  >
                    Register
                  </Link>
                </>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.nav>
  )
}

export default Header
