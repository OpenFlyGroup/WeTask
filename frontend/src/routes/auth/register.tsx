import { createFileRoute, redirect, Link } from '@tanstack/react-router'
import { useNavigate } from '@tanstack/react-router'
import { useMutation } from '@tanstack/react-query'
import { useState } from 'react'
import { motion } from 'motion/react'
import logo from 'src/assets/logo.svg'
import { Eye, EyeClosed } from 'lucide-react'
import toast from 'react-hot-toast'
import { useForm } from '@tanstack/react-form'
import AuthStorage from '@/store/auth'
import { AuthService } from '@/api/services/auth/auth.service'

export const Route = createFileRoute('/auth/register')({
  beforeLoad: () => {
    if (AuthStorage.getTokens()) {
      throw redirect({ to: '/boards' })
    }
  },
  component: RegisterPage,
})

interface RegisterFormValues {
  name: string
  email: string
  password: string
}

export function RegisterPage() {
  const navigate = useNavigate()
  const [showPassword, setShowPassword] = useState(false)

  const mutation = useMutation({
    mutationFn: (values: RegisterFormValues) => AuthService.signUp(values),
    onSuccess: () => {
      toast.success('Account created successfully!')
      navigate({ to: '/boards' })
    },
    onError: (err: any) => {
      toast.error(err?.message ?? 'Registration failed. Please try again.')
    },
  })

  const form = useForm({
    defaultValues: {
      name: '',
      email: '',
      password: '',
    },
    onSubmit: async ({ value }) => {
      mutation.mutate(value)
    },
  })

  return (
    <div className="min-h-[70vh] flex items-center justify-center py-12 px-4">
      <div className="w-full max-w-md">
        <div className="flex flex-col items-center gap-4 mb-6">
          <img src={logo} alt="WeTask logo" className="h-10" />
          <h1 className="text-2xl font-semibold">Create your account</h1>
          <p className="text-sm text-muted-foreground text-center">
            Join <span className="font-medium">WeTask</span> and start
            organizing your work better
          </p>
        </div>

        <div className="card bg-base-100 shadow-lg border border-base-200">
          <div className="card-body p-6">
            <form
              onSubmit={(e) => {
                e.preventDefault()
                e.stopPropagation()
                form.handleSubmit()
              }}
              className="flex flex-col gap-4"
            >
              {/* Full Name Field */}
              <form.Field
                name="name"
                validators={{
                  onChange: ({ value }) =>
                    !value
                      ? 'Full name is required'
                      : value.trim().length < 2
                        ? 'Name must be at least 2 characters'
                        : undefined,
                  onBlur: ({ value }) =>
                    !value
                      ? 'Full name is required'
                      : value.trim().length < 2
                        ? 'Name must be at least 2 characters'
                        : undefined,
                }}
              >
                {(field) => (
                  <label className="form-control">
                    <div className="label">
                      <span className="label-text">Full name</span>
                    </div>
                    <input
                      id={field.name}
                      value={field.state.value ?? ''}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      className={`input input-bordered w-full ${
                        field.state.meta.errors.length ? 'input-error' : ''
                      }`}
                      placeholder="John Doe"
                      autoComplete="name"
                      aria-invalid={field.state.meta.errors.length > 0}
                      aria-describedby={
                        field.state.meta.errors.length > 0
                          ? `${field.name}-error`
                          : undefined
                      }
                    />
                    {field.state.meta.errors.length > 0 && (
                      <div id={`${field.name}-error`} className="label">
                        <span className="label-text-alt text-error">
                          {field.state.meta.errors.join(', ')}
                        </span>
                      </div>
                    )}
                  </label>
                )}
              </form.Field>

              {/* Email Field */}
              <form.Field
                name="email"
                validators={{
                  onChange: ({ value }) =>
                    !value
                      ? 'Email is required'
                      : !/^\S+@\S+\.\S+$/.test(value)
                        ? 'Please enter a valid email'
                        : undefined,
                  onBlur: ({ value }) =>
                    !value
                      ? 'Email is required'
                      : !/^\S+@\S+\.\S+$/.test(value)
                        ? 'Please enter a valid email'
                        : undefined,
                }}
              >
                {(field) => (
                  <label className="form-control">
                    <div className="label">
                      <span className="label-text">Email</span>
                    </div>
                    <input
                      id={field.name}
                      value={field.state.value ?? ''}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      className={`input input-bordered w-full ${
                        field.state.meta.errors.length ? 'input-error' : ''
                      }`}
                      type="email"
                      placeholder="you@example.com"
                      autoComplete="email"
                      aria-invalid={field.state.meta.errors.length > 0}
                      aria-describedby={
                        field.state.meta.errors.length > 0
                          ? `${field.name}-error`
                          : undefined
                      }
                    />
                    {field.state.meta.errors.length > 0 && (
                      <div id={`${field.name}-error`} className="label">
                        <span className="label-text-alt text-error">
                          {field.state.meta.errors.join(', ')}
                        </span>
                      </div>
                    )}
                  </label>
                )}
              </form.Field>

              {/* Password Field */}
              <form.Field
                name="password"
                validators={{
                  onChange: ({ value }) =>
                    !value
                      ? 'Password is required'
                      : value.length < 6
                        ? 'Password must be at least 6 characters'
                        : undefined,
                  onBlur: ({ value }) =>
                    !value
                      ? 'Password is required'
                      : value.length < 6
                        ? 'Password must be at least 6 characters'
                        : undefined,
                }}
              >
                {(field) => (
                  <label className="form-control relative">
                    <div className="label">
                      <span className="label-text">Password</span>
                    </div>
                    <input
                      id={field.name}
                      value={field.state.value ?? ''}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      className={`input input-bordered w-full pr-10 ${
                        field.state.meta.errors.length ? 'input-error' : ''
                      }`}
                      type={showPassword ? 'text' : 'password'}
                      placeholder="Enter your password"
                      autoComplete="new-password"
                      aria-invalid={field.state.meta.errors.length > 0}
                      aria-describedby={
                        field.state.meta.errors.length > 0
                          ? `${field.name}-error`
                          : undefined
                      }
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword((prev) => !prev)}
                      className="absolute right-2 z-10 top-7 btn btn-ghost btn-xs btn-circle h-8 w-8"
                      aria-label={
                        showPassword ? 'Hide password' : 'Show password'
                      }
                    >
                      {showPassword ? (
                        <EyeClosed className="h-4 w-4" />
                      ) : (
                        <Eye className="h-4 w-4" />
                      )}
                    </button>
                    {field.state.meta.errors.length > 0 && (
                      <div id={`${field.name}-error`} className="label">
                        <span className="label-text-alt text-error">
                          {field.state.meta.errors.join(', ')}
                        </span>
                      </div>
                    )}
                  </label>
                )}
              </form.Field>

              {/* Submit Button */}
              <form.Subscribe
                selector={(state) => [state.canSubmit, state.isSubmitting]}
              >
                {([canSubmit, isSubmitting]) => (
                  <motion.button
                    whileHover={
                      canSubmit && !isSubmitting && !mutation.isPending
                        ? { scale: 1.02 }
                        : {}
                    }
                    whileTap={
                      canSubmit && !isSubmitting && !mutation.isPending
                        ? { scale: 0.98 }
                        : {}
                    }
                    type="submit"
                    disabled={!canSubmit || isSubmitting || mutation.isPending}
                    className="btn btn-primary mt-2 w-full"
                  >
                    {isSubmitting || mutation.isPending ? (
                      <>
                        <span className="loading loading-spinner loading-xs mr-2"></span>
                        Creating account...
                      </>
                    ) : (
                      'Create Account'
                    )}
                  </motion.button>
                )}
              </form.Subscribe>
            </form>

            <div className="text-center text-sm mt-6">
              Already have an account?{' '}
              <Link to="/auth/login" className="link link-primary font-medium">
                Sign in
              </Link>
            </div>

            <div className="mt-4 text-xs text-center text-muted-foreground">
              By signing up you agree to our{' '}
              <a
                className="link"
                href="/terms"
                target="_blank"
                rel="noopener noreferrer"
              >
                Terms
              </a>{' '}
              and{' '}
              <a
                className="link"
                href="/privacy"
                target="_blank"
                rel="noopener noreferrer"
              >
                Privacy Policy
              </a>
              .
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
