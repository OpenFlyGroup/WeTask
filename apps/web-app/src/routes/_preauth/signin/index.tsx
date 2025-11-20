import { AuthService } from '@/shared/api/services/auth/auth.service'
import { useForm } from '@tanstack/react-form'
import { useMutation } from '@tanstack/react-query'
import {
  createFileRoute,
  Link,
  redirect,
  useNavigate,
} from '@tanstack/react-router'
import { Eye, EyeClosed } from 'lucide-react'
import { useState } from 'react'
import { isBrowser, motion } from 'motion/react'
import toast from 'react-hot-toast'

import logo from 'src/assets/logo.svg'
import AuthStorage from '@/shared/store/authStore'

export const Route = createFileRoute('/_preauth/signin/')({
  beforeLoad: () => {
    if (!isBrowser) return
    if (AuthStorage.getTokens()) {
      throw redirect({ to: '/dashboard' })
    }
  },
  component: SigninPage,
})

interface SigninFormValues {
  email: string
  password: string
  remember: boolean
}

export function SigninPage() {
  const navigate = useNavigate()
  const [showPassword, setShowPassword] = useState<boolean>(false)

  const mutation = useMutation({
    mutationFn: (values: SigninFormValues) => AuthService.signIn(values),
    onSuccess: () => {
      toast.success('Logged in successfully!')
      navigate({ to: '/boards' })
    },
    onError: (err: any) => {
      toast.error(err?.message ?? 'Login failed. Please try again.')
    },
  })

  const form = useForm({
    defaultValues: {
      email: '',
      password: '',
      remember: false,
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
          <h1 className="text-2xl font-semibold">Welcome back</h1>
          <p className="text-sm text-muted-foreground text-center">
            Sign in to continue to <span className="font-medium">WeTask</span>
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
                  onBlur: ({ value }) => {
                    if (!value) return 'Email is required'
                    return undefined
                  },
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

              <form.Field
                name="password"
                validators={{
                  onChange: ({ value }) =>
                    !value
                      ? 'Password is required'
                      : value.length < 6
                        ? 'Password must be at least 6 characters'
                        : undefined,
                  onBlur: ({ value }) => {
                    if (!value) return 'Password is required'
                    if (value.length < 6)
                      return 'Password must be at least 6 characters'
                    return undefined
                  },
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
                      autoComplete="current-password"
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
                      className="absolute right-2 top-9 btn btn-ghost btn-xs btn-circle h-8 w-8"
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

              <div className="flex items-center justify-between">
                <form.Field
                  name="remember"
                  validators={{
                    onChange: () => undefined,
                  }}
                >
                  {(field) => (
                    <label className="flex items-center gap-2 cursor-pointer select-none">
                      <input
                        type="checkbox"
                        checked={field.state.value ?? false}
                        onChange={(e) => field.handleChange(e.target.checked)}
                        className="checkbox checkbox-sm"
                      />
                      <span className="text-sm">Remember me</span>
                    </label>
                  )}
                </form.Field>

                <Link to="/" className="text-sm link link-hover">
                  Forgot password?
                </Link>
              </div>

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
                        Signing in...
                      </>
                    ) : (
                      'Sign In'
                    )}
                  </motion.button>
                )}
              </form.Subscribe>
            </form>

            <div className="text-center text-sm mt-6">
              Don't have an account?{' '}
              <Link to="/signup" className="link link-primary font-medium">
                Sign Up
              </Link>
            </div>

            <div className="mt-4 text-xs text-center text-muted-foreground">
              By signing in you agree to our{' '}
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
