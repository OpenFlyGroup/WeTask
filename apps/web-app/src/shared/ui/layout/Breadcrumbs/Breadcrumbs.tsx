import { IBreadcrumbs } from '@/shared/types/ui/layout/breadcrumbs.interface'
import { Link } from '@tanstack/react-router'
import { Home } from 'lucide-react'
import { FC } from 'react'

interface IBreadcrumbsProps {
  className?: string
  breadcrumbs?: IBreadcrumbs[]
}

const Breadcrumbs: FC<IBreadcrumbsProps> = ({ className, breadcrumbs }) => {
  return (
    <div className={`breadcrumbs text-sm ${className || ''}`}>
      <ul>
        <li>
          <Link to="/dashboard">
            <Home className="size-4" />
          </Link>
        </li>
        {breadcrumbs &&
          breadcrumbs.map((breadcrumb) => (
            <li key={breadcrumb.id}>
              <Link to={breadcrumb.href}>
                {breadcrumb.icon && breadcrumb.icon}
                {breadcrumb.title}
              </Link>
            </li>
          ))}
      </ul>
    </div>
  )
}

export default Breadcrumbs
