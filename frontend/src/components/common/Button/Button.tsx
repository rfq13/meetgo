import { h, ComponentChildren } from 'preact'
import { cn } from '../../../utils/cn'
import './Button.css'

export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'outline' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  icon?: string
  children: ComponentChildren
  onClick?: (e: MouseEvent) => void
  type?: 'button' | 'submit' | 'reset'
  className?: string
}

export function Button({
  variant = 'primary',
  size = 'md',
  disabled = false,
  loading = false,
  icon,
  children,
  onClick,
  type = 'button',
  className = '',
  ...props
}: ButtonProps) {
  const baseClass = 'btn'
  
  const variantClasses = {
    primary: 'btn-primary',
    secondary: 'btn-secondary',
    outline: 'btn-outline',
    danger: 'btn-danger',
    ghost: 'btn-ghost'
  }
  
  const sizeClasses = {
    sm: 'btn-sm',
    md: 'btn-md',
    lg: 'btn-lg'
  }
  
  const classes = cn(
    baseClass,
    variantClasses[variant],
    sizeClasses[size],
    {
      'btn-disabled': disabled,
      'btn-loading': loading,
    },
    className
  )
  
  return (
    <button
      type={type}
      className={classes}
      onClick={onClick}
      disabled={disabled || loading}
      {...props}
    >
      {loading && (
        <span className="btn-spinner">
          <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        </span>
      )}
      
      {icon && !loading && (
        <span className="btn-icon">{icon}</span>
      )}
      
      <span className="btn-content">{children}</span>
    </button>
  )
}