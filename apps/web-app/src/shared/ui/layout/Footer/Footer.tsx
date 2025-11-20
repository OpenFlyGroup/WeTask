import logo from 'src/assets/openfly_h_logo.svg'
import { motion } from 'motion/react'

const Footer = () => {
  return (
    <motion.footer
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.2 }}
      className="flex flex-col md:flex-row items-center justify-between p-4 bg-base-100 text-base-content shadow-lg"
    >
      <aside>
        <a
          href="https://openflygroup.github.io/enterprise_landing/"
          target="_blank"
        >
          <img src={logo} className="" alt="OpenFly" />
        </a>
        <p>Tools for a conscious life</p>
      </aside>
      <div>
        <p>© 2025 OpenFly. All rights reserved.</p>
      </div>
    </motion.footer>
  )
}

export default Footer
