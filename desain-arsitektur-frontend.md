# Desain Arsitektur Frontend (Preact) dengan Tema Hijau dan Abu-abu

## 1. Arsitektur Umum Frontend

### 1.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              Preact Application                                │
└─────────────────────────────────────────────────────────────────────────────────┘
                                         │
                    ┌────────────────────┼────────────────────┐
                    │                    │                    │
┌───────────────────▼───────────┐ ┌─────▼────────────────┐ ┌─▼───────────────────┐
│   Components Layer           │ │   State Management   │ │   Services Layer     │
│                              │ │                      │ │                      │
│  ┌─────────────────────────┐ │ │  ┌─────────────────┐  │ │  ┌─────────────────┐ │
│  │ Layout Components       │ │ │  │ Zustand Store   │  │ │  │ API Service     │ │
│  │                         │ │ │  │                 │  │ │  │                 │ │
│  │ - Header                │ │ │  │ - User Store    │  │ │  │ - Auth Service  │ │
│  │ - Sidebar               │ │ │  │ - Room Store    │  │ │  │ - Room Service  │ │
│  │ - Footer                │ │ │  │ - UI Store      │  │ │  │ - User Service  │ │
│  └─────────────────────────┘ │ │  └─────────────────┘  │ │  └─────────────────┘ │
│                             │ │                      │ │                      │
│  ┌─────────────────────────┐ │ │  ┌─────────────────┐  │ │  ┌─────────────────┐ │
│  │ Page Components         │ │ │  │ Context API     │  │ │  │ WebSocket       │ │
│  │                         │ │ │  │                 │  │ │  │ Service         │ │
│  │ - LoginPage             │ │ │  │ - Theme Context │  │ │  │                 │ │
│  │ - DashboardPage         │ │ │  │ - Auth Context  │  │ │  └─────────────────┘ │
│  │ - MeetingPage           │ │ │  └─────────────────┘  │ │                      │
│  │ - ProfilePage           │ │ │                      │ │  ┌─────────────────┐ │
│  └─────────────────────────┘ │ └───────────────────────┘ │  │ WebRTC Service  │ │
│                             │                          │ │  │                 │ │
│  ┌─────────────────────────┐ │                          │ │  │ - Peer Manager │ │
│  │ UI Components           │ │                          │ │  │ - Media Handler │ │
│  │                         │ │                          │ │  └─────────────────┘ │
│  │ - VideoPlayer           │ │                          │ │                      │
│  │ - AudioControls         │ │                          │ │  ┌─────────────────┐ │
│  │ - ParticipantList       │ │                          │ │  │ Storage Service │ │
│  │ - ChatBox               │ │                          │ │  │                 │ │
│  │ - SettingsPanel         │ │                          │ │  │ - Local Storage │ │
│  └─────────────────────────┘ │                          │ │  │ - Session Mgmt  │ │
│                             │                          │ │  └─────────────────┘ │
│  ┌─────────────────────────┐ │                          │ │                      │
│  │ Form Components         │ │                          │ │  ┌─────────────────┐ │
│  │                         │ │                          │ │  │ Utility         │ │
│  │ - LoginForm             │ │                          │ │  │                 │ │
│  │ - RegisterForm          │ │                          │ │  │ - Date Utils    │ │
│  │ - CreateRoomForm        │ │                          │ │  │ - Validator     │ │
│  │ - ProfileForm           │ │                          │ │  │ - Formatters    │ │
│  └─────────────────────────┘ │                          │ │  └─────────────────┘ │
└─────────────────────────────┘ └──────────────────────────┘ └──────────────────────┘
```

### 1.2 Teknologi Stack Frontend

#### 1.2.1 Core Framework
- **UI Framework**: Preact (lightweight alternative to React)
- **State Management**: Zustand (simple, fast, and scalable state management)
- **Routing**: Preact Router (official router for Preact)
- **Styling**: Tailwind CSS dengan tema hijau dan abu-abu
- **Build Tool**: Vite (fast build tool for modern web apps)
- **TypeScript**: Untuk type safety dan developer experience

#### 1.2.2 WebRTC & Real-time
- **WebRTC Client**: Simple-peer (simplified WebRTC implementation)
- **WebSocket**: Native WebSocket API atau socket.io-client
- **Media Handling**: Custom media utilities untuk camera, microphone, dan screen sharing

#### 1.2.3 UI/UX
- **Component Library**: Custom components dengan Tailwind CSS
- **Icons**: Heroicons atau custom SVG icons
- **Animations**: Framer Motion atau CSS animations
- **Responsive Design**: Mobile-first approach dengan Tailwind CSS

#### 1.2.4 Development Tools
- **Linting**: ESLint dengan Preact plugin
- **Formatting**: Prettier untuk code formatting
- **Testing**: Vitest untuk unit testing dan Testing Library untuk component testing
- **Bundle Analysis**: Rollup Plugin Visualizer

## 2. Struktur Folder Frontend

```
frontend/
├── public/                        # Static assets
│   ├── favicon.ico                # Favicon
│   ├── manifest.json              # PWA manifest
│   └── assets/                    # Images, fonts, etc.
│       ├── icons/                 # SVG icons
│       ├── images/                # Images
│       └── fonts/                 # Custom fonts
├── src/                          # Source code
│   ├── components/                # Reusable components
│   │   ├── common/               # Common components used across the app
│   │   │   ├── Button/           # Button component
│   │   │   │   ├── index.tsx     # Component export
│   │   │   │   ├── Button.tsx    # Button implementation
│   │   │   │   └── styles.css    # Button styles
│   │   │   ├── Input/            # Input component
│   │   │   ├── Modal/            # Modal component
│   │   │   ├── Loading/          # Loading component
│   │   │   └── Avatar/           # Avatar component
│   │   ├── layout/               # Layout components
│   │   │   ├── Header/           # Header component
│   │   │   ├── Sidebar/          # Sidebar component
│   │   │   ├── Footer/           # Footer component
│   │   │   └── MainLayout/       # Main layout wrapper
│   │   ├── forms/                # Form components
│   │   │   ├── LoginForm/        # Login form
│   │   │   ├── RegisterForm/     # Register form
│   │   │   ├── CreateRoomForm/   # Create room form
│   │   │   └── ProfileForm/      # Profile form
│   │   └── meeting/              # Meeting-specific components
│   │       ├── VideoPlayer/      # Video player component
│   │       ├── AudioControls/    # Audio controls component
│   │       ├── ParticipantList/  # Participant list component
│   │       ├── ChatBox/          # Chat box component
│   │       ├── SettingsPanel/    # Settings panel component
│   │       ├── ScreenShare/      # Screen share component
│   │       └── Recording/        # Recording component
│   ├── pages/                    # Page components
│   │   ├── Home/                 # Home page
│   │   │   ├── index.tsx         # Home page component
│   │   │   └── styles.css        # Home page styles
│   │   ├── Login/                # Login page
│   │   ├── Register/             # Register page
│   │   ├── Dashboard/            # Dashboard page
│   │   ├── Meeting/              # Meeting page
│   │   ├── Profile/              # Profile page
│   │   └── NotFound/             # 404 page
│   ├── stores/                   # State management (Zustand)
│   │   ├── authStore.ts          # Authentication store
│   │   ├── userStore.ts          # User store
│   │   ├── roomStore.ts          # Room store
│   │   ├── uiStore.ts            # UI store
│   │   └── webrtcStore.ts        # WebRTC store
│   ├── services/                 # API and service layer
│   │   ├── api/                  # API services
│   │   │   ├── authApi.ts        # Authentication API
│   │   │   ├── userApi.ts        # User API
│   │   │   └── roomApi.ts        # Room API
│   │   ├── websocket/            # WebSocket service
│   │   │   ├── index.ts          # WebSocket service export
│   │   │   ├── WebSocketService.ts # WebSocket implementation
│   │   │   └── types.ts          # WebSocket types
│   │   ├── webrtc/               # WebRTC service
│   │   │   ├── index.ts          # WebRTC service export
│   │   │   ├── PeerManager.ts    # WebRTC peer manager
│   │   │   ├── MediaHandler.ts   # Media handler
│   │   │   └── types.ts          # WebRTC types
│   │   └── storage/              # Storage service
│   │       ├── index.ts          # Storage service export
│   │       ├── LocalStorage.ts   # Local storage implementation
│   │       └── SessionStorage.ts  # Session storage implementation
│   ├── hooks/                    # Custom hooks
│   │   ├── useAuth.ts            # Authentication hook
│   │   ├── useUser.ts            # User data hook
│   │   ├── useRoom.ts            # Room data hook
│   │   ├── useWebRTC.ts          # WebRTC hook
│   │   ├── useWebSocket.ts       # WebSocket hook
│   │   ├── useMedia.ts           # Media devices hook
│   │   ├── useLocalStorage.ts    # Local storage hook
│   │   └── useDebounce.ts        # Debounce hook
│   ├── utils/                    # Utility functions
│   │   ├── dateUtils.ts          # Date utilities
│   │   ├── validationUtils.ts    # Validation utilities
│   │   ├── formatUtils.ts        # Format utilities
│   │   ├── constants.ts          # App constants
│   │   └── config.ts             # App configuration
│   ├── contexts/                 # React contexts
│   │   ├── ThemeContext.tsx      # Theme context
│   │   └── AuthContext.tsx       # Auth context
│   ├── styles/                   # Global styles
│   │   ├── global.css            # Global CSS
│   │   ├── tailwind.css          # Tailwind CSS
│   │   └── theme.css             # Theme-specific styles
│   ├── types/                    # TypeScript type definitions
│   │   ├── auth.ts               # Auth types
│   │   ├── user.ts               # User types
│   │   ├── room.ts               # Room types
│   │   ├── webrtc.ts             # WebRTC types
│   │   ├── api.ts                # API types
│   │   └── index.ts              # Type exports
│   ├── App.tsx                   # Root App component
│   ├── main.tsx                  # App entry point
│   └── index.css                 # App entry CSS
├── tests/                        # Test files
│   ├── components/               # Component tests
│   ├── pages/                    # Page tests
│   ├── services/                 # Service tests
│   ├── hooks/                    # Hook tests
│   └── utils/                    # Utility tests
├── docs/                         # Documentation
│   ├── components.md             # Component documentation
│   ├── api.md                    # API documentation
│   └── deployment.md             # Deployment documentation
├── .env.example                  # Environment variables example
├── .eslintrc.js                  # ESLint configuration
├── .prettierrc.js                # Prettier configuration
├── tailwind.config.js            # Tailwind CSS configuration
├── tsconfig.json                 # TypeScript configuration
├── vite.config.ts                # Vite configuration
├── package.json                  # Dependencies and scripts
└── README.md                     # Project README
```

## 3. Tema Hijau dan Abu-abu

### 3.1 Desain Tema

Tema hijau dan abu-abu untuk aplikasi WebRTC meeting dirancang untuk memberikan kesan profesional, modern, dan menenangkan. Berikut adalah palet warna yang akan digunakan:

#### 3.1.1 Palet Warna

```css
/* Primary Colors (Hijau) */
--green-50: #f0fdf4;
--green-100: #dcfce7;
--green-200: #bbf7d0;
--green-300: #86efac;
--green-400: #4ade80;
--green-500: #22c55e; /* Primary Green */
--green-600: #16a34a;
--green-700: #15803d;
--green-800: #166534;
--green-900: #14532d;

/* Secondary Colors (Abu-abu) */
--gray-50: #f9fafb;
--gray-100: #f3f4f6;
--gray-200: #e5e7eb;
--gray-300: #d1d5db;
--gray-400: #9ca3af;
--gray-500: #6b7280;
--gray-600: #4b5563;
--gray-700: #374151;
--gray-800: #1f2937;
--gray-900: #111827;

/* Accent Colors */
--accent-500: #10b981; /* Emerald untuk aksen */
--accent-600: #059669;

/* Status Colors */
--success-500: #22c55e;
--warning-500: #f59e0b;
--error-500: #ef4444;
--info-500: #3b82f6;
```

#### 3.1.2 Konfigurasi Tailwind CSS

```javascript
// tailwind.config.js
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Hijau sebagai warna utama
        primary: {
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0',
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e', // Primary Green
          600: '#16a34a',
          700: '#15803d',
          800: '#166534',
          900: '#14532d',
        },
        // Abu-abu sebagai warna sekunder
        secondary: {
          50: '#f9fafb',
          100: '#f3f4f6',
          200: '#e5e7eb',
          300: '#d1d5db',
          400: '#9ca3af',
          500: '#6b7280',
          600: '#4b5563',
          700: '#374151',
          800: '#1f2937',
          900: '#111827',
        },
        // Aksen untuk强调元素
        accent: {
          500: '#10b981',
          600: '#059669',
        },
        // Status colors
        success: '#22c55e',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#3b82f6',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        'custom': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
        'custom-lg': '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
```

#### 3.1.3 Global Styles

```css
/* src/styles/global.css */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

:root {
  /* Primary Colors (Hijau) */
  --color-primary-50: #f0fdf4;
  --color-primary-100: #dcfce7;
  --color-primary-200: #bbf7d0;
  --color-primary-300: #86efac;
  --color-primary-400: #4ade80;
  --color-primary-500: #22c55e;
  --color-primary-600: #16a34a;
  --color-primary-700: #15803d;
  --color-primary-800: #166534;
  --color-primary-900: #14532d;

  /* Secondary Colors (Abu-abu) */
  --color-secondary-50: #f9fafb;
  --color-secondary-100: #f3f4f6;
  --color-secondary-200: #e5e7eb;
  --color-secondary-300: #d1d5db;
  --color-secondary-400: #9ca3af;
  --color-secondary-500: #6b7280;
  --color-secondary-600: #4b5563;
  --color-secondary-700: #374151;
  --color-secondary-800: #1f2937;
  --color-secondary-900: #111827;

  /* Accent Colors */
  --color-accent-500: #10b981;
  --color-accent-600: #059669;

  /* Status Colors */
  --color-success: #22c55e;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
  --color-info: #3b82f6;

  /* Border Radius */
  --radius-sm: 0.25rem;
  --radius-md: 0.375rem;
  --radius-lg: 0.5rem;
  --radius-xl: 0.75rem;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);

  /* Transitions */
  --transition-fast: 150ms ease-in-out;
  --transition-normal: 300ms ease-in-out;
  --transition-slow: 500ms ease-in-out;
}

* {
  box-sizing: border-box;
}

html {
  font-size: 16px;
  scroll-behavior: smooth;
}

body {
  margin: 0;
  font-family: 'Inter', system-ui, sans-serif;
  background-color: var(--color-secondary-50);
  color: var(--color-secondary-900);
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* Custom Scrollbar */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: var(--color-secondary-100);
  border-radius: var(--radius-md);
}

::-webkit-scrollbar-thumb {
  background: var(--color-secondary-300);
  border-radius: var(--radius-md);
}

::-webkit-scrollbar-thumb:hover {
  background: var(--color-secondary-400);
}

/* Focus Styles */
*:focus {
  outline: 2px solid var(--color-primary-500);
  outline-offset: 2px;
}

/* Button Styles */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--transition-fast);
  border: none;
  cursor: pointer;
  white-space: nowrap;
}

.btn-primary {
  background-color: var(--color-primary-500);
  color: white;
}

.btn-primary:hover {
  background-color: var(--color-primary-600);
}

.btn-secondary {
  background-color: var(--color-secondary-200);
  color: var(--color-secondary-800);
}

.btn-secondary:hover {
  background-color: var(--color-secondary-300);
}

.btn-outline {
  background-color: transparent;
  color: var(--color-primary-500);
  border: 1px solid var(--color-primary-500);
}

.btn-outline:hover {
  background-color: var(--color-primary-500);
  color: white;
}

/* Input Styles */
.input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--color-secondary-300);
  border-radius: var(--radius-md);
  background-color: white;
  transition: border-color var(--transition-fast);
}

.input:focus {
  border-color: var(--color-primary-500);
}

/* Card Styles */
.card {
  background-color: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: 1.5rem;
  transition: box-shadow var(--transition-normal);
}

.card:hover {
  box-shadow: var(--shadow-lg);
}

/* Animation Classes */
.animate-fade-in {
  animation: fadeIn 0.5s ease-in-out;
}

.animate-slide-up {
  animation: slideUp 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    transform: translateY(10px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

/* Utility Classes */
.text-primary {
  color: var(--color-primary-500);
}

.text-secondary {
  color: var(--color-secondary-600);
}

.bg-primary {
  background-color: var(--color-primary-500);
}

.bg-secondary {
  background-color: var(--color-secondary-100);
}

.border-primary {
  border-color: var(--color-primary-500);
}

.border-secondary {
  border-color: var(--color-secondary-300);
}
```

## 4. Komponen Utama Frontend

### 4.1 Layout Components

#### 4.1.1 Header Component

```tsx
// src/components/layout/Header/Header.tsx
import { h } from 'preact';
import { Link } from 'preact-router';
import { useAuthStore } from '../../../stores/authStore';
import { Button } from '../../common/Button/Button';
import { Avatar } from '../../common/Avatar/Avatar';
import './Header.css';

export function Header() {
  const { user, logout } = useAuthStore();

  return (
    <header className="header">
      <div className="header-container">
        <div className="header-logo">
          <Link href="/">
            <h1 className="logo-text">MeetGreen</h1>
          </Link>
        </div>
        
        <nav className="header-nav">
          <Link href="/dashboard" className="nav-link">Dashboard</Link>
          <Link href="/rooms" className="nav-link">Rooms</Link>
          <Link href="/contacts" className="nav-link">Contacts</Link>
        </nav>
        
        <div className="header-user">
          {user ? (
            <div className="user-menu">
              <Avatar src={user.avatar} alt={user.firstName} size="md" />
              <div className="user-dropdown">
                <Link href="/profile" className="dropdown-item">Profile</Link>
                <Link href="/settings" className="dropdown-item">Settings</Link>
                <Button variant="secondary" size="sm" onClick={logout}>
                  Logout
                </Button>
              </div>
            </div>
          ) : (
            <div className="auth-buttons">
              <Link href="/login">
                <Button variant="outline" size="sm">Login</Button>
              </Link>
              <Link href="/register">
                <Button size="sm">Register</Button>
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
```

```css
/* src/components/layout/Header/Header.css */
.header {
  background-color: white;
  box-shadow: var(--shadow-md);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-logo {
  display: flex;
  align-items: center;
}

.logo-text {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-primary-600);
  margin: 0;
}

.header-nav {
  display: flex;
  gap: 1.5rem;
}

.nav-link {
  color: var(--color-secondary-700);
  text-decoration: none;
  font-weight: 500;
  transition: color var(--transition-fast);
}

.nav-link:hover {
  color: var(--color-primary-600);
}

.header-user {
  display: flex;
  align-items: center;
}

.user-menu {
  position: relative;
}

.user-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.5rem;
  background-color: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  min-width: 200px;
  padding: 0.5rem;
  display: none;
  z-index: 10;
}

.user-menu:hover .user-dropdown {
  display: block;
}

.dropdown-item {
  display: block;
  padding: 0.5rem 0.75rem;
  color: var(--color-secondary-700);
  text-decoration: none;
  border-radius: var(--radius-sm);
  transition: background-color var(--transition-fast);
}

.dropdown-item:hover {
  background-color: var(--color-secondary-100);
}

.auth-buttons {
  display: flex;
  gap: 0.75rem;
}

@media (max-width: 768px) {
  .header-nav {
    display: none;
  }
  
  .header-container {
    padding: 1rem;
  }
}
```

#### 4.1.2 Sidebar Component

```tsx
// src/components/layout/Sidebar/Sidebar.tsx
import { h } from 'preact';
import { Link } from 'preact-router';
import { Avatar } from '../../common/Avatar/Avatar';
import { Button } from '../../common/Button/Button';
import './Sidebar.css';

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

interface MenuItem {
  icon: string;
  label: string;
  href: string;
  count?: number;
}

export function Sidebar({ isOpen, onClose }: SidebarProps) {
  const menuItems: MenuItem[] = [
    { icon: 'dashboard', label: 'Dashboard', href: '/dashboard' },
    { icon: 'video', label: 'My Rooms', href: '/rooms', count: 5 },
    { icon: 'users', label: 'Contacts', href: '/contacts', count: 12 },
    { icon: 'calendar', label: 'Schedule', href: '/schedule' },
    { icon: 'settings', label: 'Settings', href: '/settings' },
  ];

  return (
    <aside className={`sidebar ${isOpen ? 'sidebar-open' : ''}`}>
      <div className="sidebar-header">
        <h2 className="sidebar-title">Menu</h2>
        <Button
          variant="ghost"
          size="sm"
          icon="close"
          onClick={onClose}
          className="sidebar-close"
        />
      </div>
      
      <nav className="sidebar-nav">
        <ul className="nav-list">
          {menuItems.map((item) => (
            <li key={item.href} className="nav-item">
              <Link href={item.href} className="nav-link" onClick={onClose}>
                <span className="nav-icon">{item.icon}</span>
                <span className="nav-label">{item.label}</span>
                {item.count && (
                  <span className="nav-count">{item.count}</span>
                )}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
      
      <div className="sidebar-footer">
        <div className="user-profile">
          <Avatar src="/assets/avatar.jpg" alt="User" size="md" />
          <div className="user-info">
            <div className="user-name">John Doe</div>
            <div className="user-status">Online</div>
          </div>
        </div>
        <Button variant="outline" size="sm" icon="logout">
          Logout
        </Button>
      </div>
    </aside>
  );
}
```

```css
/* src/components/layout/Sidebar/Sidebar.css */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 260px;
  height: 100vh;
  background-color: white;
  box-shadow: var(--shadow-lg);
  z-index: 1000;
  transform: translateX(-100%);
  transition: transform var(--transition-normal);
  display: flex;
  flex-direction: column;
}

.sidebar-open {
  transform: translateX(0);
}

.sidebar-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--color-secondary-200);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-secondary-800);
  margin: 0;
}

.sidebar-close {
  color: var(--color-secondary-600);
}

.sidebar-nav {
  flex: 1;
  padding: 1rem 0;
  overflow-y: auto;
}

.nav-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-item {
  margin-bottom: 0.25rem;
}

.nav-link {
  display: flex;
  align-items: center;
  padding: 0.75rem 1.5rem;
  color: var(--color-secondary-700);
  text-decoration: none;
  transition: background-color var(--transition-fast);
  position: relative;
}

.nav-link:hover {
  background-color: var(--color-secondary-100);
}

.nav-link.active {
  background-color: var(--color-primary-50);
  color: var(--color-primary-600);
}

.nav-link.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background-color: var(--color-primary-500);
}

.nav-icon {
  margin-right: 0.75rem;
  font-size: 1.25rem;
}

.nav-label {
  flex: 1;
}

.nav-count {
  background-color: var(--color-primary-500);
  color: white;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
}

.sidebar-footer {
  padding: 1.5rem;
  border-top: 1px solid var(--color-secondary-200);
}

.user-profile {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.user-info {
  margin-left: 0.75rem;
}

.user-name {
  font-weight: 600;
  color: var(--color-secondary-800);
}

.user-status {
  font-size: 0.875rem;
  color: var(--color-secondary-600);
}

@media (max-width: 768px) {
  .sidebar {
    width: 100%;
    max-width: 300px;
  }
}
```

### 4.2 Meeting Components

#### 4.2.1 VideoPlayer Component

```tsx
// src/components/meeting/VideoPlayer/VideoPlayer.tsx
import { h } from 'preact';
import { useEffect, useRef, useState } from 'preact/hooks';
import { Avatar } from '../../common/Avatar/Avatar';
import { Button } from '../../common/Button/Button';
import './VideoPlayer.css';

interface VideoPlayerProps {
  stream?: MediaStream;
  isLocal?: boolean;
  username?: string;
  isAudioEnabled?: boolean;
  isVideoEnabled?: boolean;
  isScreenSharing?: boolean;
  isSpeaking?: boolean;
  onToggleAudio?: () => void;
  onToggleVideo?: () => void;
  className?: string;
}

export function VideoPlayer({
  stream,
  isLocal = false,
  username = 'User',
  isAudioEnabled = true,
  isVideoEnabled = true,
  isScreenSharing = false,
  isSpeaking = false,
  onToggleAudio,
  onToggleVideo,
  className = '',
}: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isVideoLoaded, setIsVideoLoaded] = useState(false);

  useEffect(() => {
    if (videoRef.current && stream) {
      videoRef.current.srcObject = stream;
      setIsVideoLoaded(true);
    } else {
      setIsVideoLoaded(false);
    }
  }, [stream]);

  const handleVideoLoaded = () => {
    setIsVideoLoaded(true);
  };

  return (
    <div className={`video-player ${className} ${isSpeaking ? 'speaking' : ''}`}>
      {isVideoEnabled && stream ? (
        <video
          ref={videoRef}
          autoPlay
          playsInline
          muted={isLocal}
          onLoadedData={handleVideoLoaded}
          className={`video-element ${isVideoLoaded ? 'loaded' : ''} ${isScreenSharing ? 'screen-share' : ''}`}
        />
      ) : (
        <div className="video-placeholder">
          <Avatar
            src=""
            alt={username}
            size="lg"
            className="avatar-placeholder"
          />
          <div className="username">{username}</div>
          {isLocal && (
            <div className="local-badge">You</div>
          )}
        </div>
      )}
      
      <div className="video-overlay">
        <div className="video-info">
          <div className="username">{username}</div>
          <div className="status-indicators">
            {!isAudioEnabled && (
              <div className="status-indicator muted">
                <span className="icon">mic_off</span>
              </div>
            )}
            {!isVideoEnabled && (
              <div className="status-indicator video-off">
                <span className="icon">videocam_off</span>
              </div>
            )}
            {isScreenSharing && (
              <div className="status-indicator screen-share">
                <span className="icon">screen_share</span>
              </div>
            )}
          </div>
        </div>
        
        {isLocal && (
          <div className="video-controls">
            <Button
              variant="icon"
              size="sm"
              icon={isAudioEnabled ? 'mic' : 'mic_off'}
              onClick={onToggleAudio}
              className={isAudioEnabled ? '' : 'active'}
            />
            <Button
              variant="icon"
              size="sm"
              icon={isVideoEnabled ? 'videocam' : 'videocam_off'}
              onClick={onToggleVideo}
              className={isVideoEnabled ? '' : 'active'}
            />
          </div>
        )}
      </div>
    </div>
  );
}
```

```css
/* src/components/meeting/VideoPlayer/VideoPlayer.css */
.video-player {
  position: relative;
  background-color: var(--color-secondary-800);
  border-radius: var(--radius-lg);
  overflow: hidden;
  aspect-ratio: 16 / 9;
  box-shadow: var(--shadow-md);
  transition: box-shadow var(--transition-normal);
}

.video-player:hover {
  box-shadow: var(--shadow-lg);
}

.video-player.speaking {
  box-shadow: 0 0 0 3px var(--color-primary-500);
}

.video-element {
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0;
  transition: opacity var(--transition-normal);
}

.video-element.loaded {
  opacity: 1;
}

.video-element.screen-share {
  object-fit: contain;
}

.video-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: var(--color-secondary-800);
  color: white;
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  margin-bottom: 1rem;
}

.username {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
}

.local-badge {
  position: absolute;
  top: 1rem;
  left: 1rem;
  background-color: var(--color-primary-500);
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: 600;
}

.video-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, rgba(0, 0, 0, 0.3) 0%, transparent 30%, transparent 70%, rgba(0, 0, 0, 0.3) 100%);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 1rem;
  opacity: 0;
  transition: opacity var(--transition-fast);
  pointer-events: none;
}

.video-player:hover .video-overlay {
  opacity: 1;
}

.video-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-indicators {
  display: flex;
  gap: 0.25rem;
}

.status-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background-color: rgba(0, 0, 0, 0.6);
  border-radius: 50%;
  color: white;
}

.status-indicator.active {
  background-color: var(--color-error);
}

.video-controls {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  pointer-events: all;
}

@media (max-width: 768px) {
  .video-player {
    aspect-ratio: 4 / 3;
  }
  
  .video-overlay {
    opacity: 1;
    background: linear-gradient(to bottom, rgba(0, 0, 0, 0.5) 0%, transparent 20%, transparent 80%, rgba(0, 0, 0, 0.5) 100%);
  }
}
```

#### 4.2.2 Meeting Controls Component

```tsx
// src/components/meeting/MeetingControls/MeetingControls.tsx
import { h } from 'preact';
import { Button } from '../../common/Button/Button';
import './MeetingControls.css';

interface MeetingControlsProps {
  isAudioEnabled: boolean;
  isVideoEnabled: boolean;
  isScreenSharing: boolean;
  isRecording: boolean;
  participantCount: number;
  onToggleAudio: () => void;
  onToggleVideo: () => void;
  onToggleScreenShare: () => void;
  onToggleRecording: () => void;
  onShowParticipants: () => void;
  onShowChat: () => void;
  onLeaveMeeting: () => void;
}

export function MeetingControls({
  isAudioEnabled,
  isVideoEnabled,
  isScreenSharing,
  isRecording,
  participantCount,
  onToggleAudio,
  onToggleVideo,
  onToggleScreenShare,
  onToggleRecording,
  onShowParticipants,
  onShowChat,
  onLeaveMeeting,
}: MeetingControlsProps) {
  return (
    <div className="meeting-controls">
      <div className="controls-left">
        <Button
          variant="icon"
          size="lg"
          icon={isAudioEnabled ? 'mic' : 'mic_off'}
          onClick={onToggleAudio}
          className={isAudioEnabled ? '' : 'active'}
          title={isAudioEnabled ? 'Mute audio' : 'Unmute audio'}
        />
        <Button
          variant="icon"
          size="lg"
          icon={isVideoEnabled ? 'videocam' : 'videocam_off'}
          onClick={onToggleVideo}
          className={isVideoEnabled ? '' : 'active'}
          title={isVideoEnabled ? 'Turn off video' : 'Turn on video'}
        />
        <Button
          variant="icon"
          size="lg"
          icon={isScreenSharing ? 'stop_screen_share' : 'screen_share'}
          onClick={onToggleScreenShare}
          className={isScreenSharing ? 'active' : ''}
          title={isScreenSharing ? 'Stop screen sharing' : 'Share screen'}
        />
        <Button
          variant="icon"
          size="lg"
          icon={isRecording ? 'stop_circle' : 'fiber_manual_record'}
          onClick={onToggleRecording}
          className={isRecording ? 'active recording' : ''}
          title={isRecording ? 'Stop recording' : 'Start recording'}
        />
      </div>
      
      <div className="controls-center">
        <Button
          variant="icon"
          size="lg"
          icon="people"
          onClick={onShowParticipants}
          title="Show participants"
          badge={participantCount > 1 ? participantCount : undefined}
        />
        <Button
          variant="icon"
          size="lg"
          icon="chat"
          onClick={onShowChat}
          title="Show chat"
        />
      </div>
      
      <div className="controls-right">
        <Button
          variant="danger"
          size="lg"
          icon="call_end"
          onClick={onLeaveMeeting}
          title="Leave meeting"
        >
          Leave
        </Button>
      </div>
    </div>
  );
}
```

```css
/* src/components/meeting/MeetingControls/MeetingControls.css */
.meeting-controls {
  position: fixed;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  padding: 0.75rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 1.5rem;
  z-index: 100;
}

.controls-left,
.controls-center,
.controls-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.controls-center {
  border-left: 1px solid var(--color-secondary-200);
  border-right: 1px solid var(--color-secondary-200);
  padding-left: 1.5rem;
  padding-right: 1.5rem;
}

.controls-right {
  margin-left: 0.5rem;
}

.meeting-controls .btn {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-secondary-100);
  color: var(--color-secondary-700);
  transition: all var(--transition-fast);
}

.meeting-controls .btn:hover {
  background-color: var(--color-secondary-200);
  transform: scale(1.05);
}

.meeting-controls .btn.active {
  background-color: var(--color-error);
  color: white;
}

.meeting-controls .btn.recording {
  animation: pulse 2s infinite;
}

.meeting-controls .btn-danger {
  background-color: var(--color-error);
  color: white;
  width: auto;
  padding: 0 1rem;
  border-radius: var(--radius-md);
}

.meeting-controls .btn-danger:hover {
  background-color: #dc2626;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.7);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(239, 68, 68, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0);
  }
}

@media (max-width: 768px) {
  .meeting-controls {
    bottom: 1rem;
    padding: 0.5rem 1rem;
    gap: 1rem;
  }
  
  .controls-left,
  .controls-center,
  .controls-right {
    gap: 0.5rem;
  }
  
  .controls-center {
    padding-left: 1rem;
    padding-right: 1rem;
  }
  
  .meeting-controls .btn {
    width: 40px;
    height: 40px;
  }
  
  .meeting-controls .btn-danger {
    padding: 0 0.75rem;
    font-size: 0.875rem;
  }
}
```

### 4.3 State Management dengan Zustand

#### 4.3.1 Auth Store

```ts
// src/stores/authStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authService } from '../services/api/authApi';
import { User } from '../types/auth';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
  clearError: () => void;
  setUser: (user: User) => void;
  setToken: (token: string) => void;
}

interface RegisterData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      login: async (email: string, password: string) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authService.login(email, password);
          set({
            user: response.user,
            token: response.token,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            isLoading: false,
            error: error instanceof Error ? error.message : 'Login failed',
          });
        }
      },

      register: async (data: RegisterData) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authService.register(data);
          set({
            user: response.user,
            token: response.token,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            isLoading: false,
            error: error instanceof Error ? error.message : 'Registration failed',
          });
        }
      },

      logout: () => {
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          isLoading: false,
          error: null,
        });
      },

      clearError: () => {
        set({ error: null });
      },

      setUser: (user: User) => {
        set({ user });
      },

      setToken: (token: string) => {
        set({ token, isAuthenticated: true });
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
```

#### 4.3.2 WebRTC Store

```ts
// src/stores/webrtcStore.ts
import { create } from 'zustand';
import { MediaStream, RTCPeerConnection } from '../types/webrtc';
import { PeerManager } from '../services/webrtc/PeerManager';

interface WebRTCState {
  localStream: MediaStream | null;
  remoteStreams: Map<string, MediaStream>;
  peerConnections: Map<string, RTCPeerConnection>;
  isAudioEnabled: boolean;
  isVideoEnabled: boolean;
  isScreenSharing: boolean;
  isRecording: boolean;
  currentRoomId: string | null;
  peerManager: PeerManager | null;
  
  // Actions
  initializeMedia: () => Promise<void>;
  toggleAudio: () => void;
  toggleVideo: () => void;
  toggleScreenShare: () => Promise<void>;
  toggleRecording: () => void;
  joinRoom: (roomId: string) => void;
  leaveRoom: () => void;
  addRemoteStream: (peerId: string, stream: MediaStream) => void;
  removeRemoteStream: (peerId: string) => void;
  cleanup: () => void;
}

export const useWebRTCStore = create<WebRTCState>((set, get) => ({
  localStream: null,
  remoteStreams: new Map(),
  peerConnections: new Map(),
  isAudioEnabled: true,
  isVideoEnabled: true,
  isScreenSharing: false,
  isRecording: false,
  currentRoomId: null,
  peerManager: null,

  initializeMedia: async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: true,
        video: true,
      });
      
      set({ localStream: stream });
    } catch (error) {
      console.error('Error accessing media devices:', error);
      throw error;
    }
  },

  toggleAudio: () => {
    const { localStream, isAudioEnabled } = get();
    if (localStream) {
      const audioTracks = localStream.getAudioTracks();
      audioTracks.forEach(track => {
        track.enabled = !isAudioEnabled;
      });
      set({ isAudioEnabled: !isAudioEnabled });
    }
  },

  toggleVideo: () => {
    const { localStream, isVideoEnabled } = get();
    if (localStream) {
      const videoTracks = localStream.getVideoTracks();
      videoTracks.forEach(track => {
        track.enabled = !isVideoEnabled;
      });
      set({ isVideoEnabled: !isVideoEnabled });
    }
  },

  toggleScreenShare: async () => {
    const { localStream, isScreenSharing } = get();
    
    if (isScreenSharing) {
      // Stop screen sharing and switch back to camera
      if (localStream) {
        const screenTrack = localStream.getVideoTracks().find(track => 
          track.label.includes('screen')
        );
        if (screenTrack) {
          screenTrack.stop();
          localStream.removeTrack(screenTrack);
        }
        
        // Get camera stream again
        try {
          const cameraStream = await navigator.mediaDevices.getUserMedia({
            video: true,
          });
          const cameraTrack = cameraStream.getVideoTracks()[0];
          localStream.addTrack(cameraTrack);
        } catch (error) {
          console.error('Error accessing camera:', error);
        }
      }
      set({ isScreenSharing: false });
    } else {
      // Start screen sharing
      try {
        const screenStream = await navigator.mediaDevices.getDisplayMedia({
          video: true,
        });
        
        if (localStream) {
          // Remove camera track
          const cameraTrack = localStream.getVideoTracks().find(track => 
            !track.label.includes('screen')
          );
          if (cameraTrack) {
            localStream.removeTrack(cameraTrack);
          }
          
          // Add screen track
          const screenTrack = screenStream.getVideoTracks()[0];
          localStream.addTrack(screenTrack);
          
          // Handle screen share end
          screenTrack.onended = () => {
            get().toggleScreenShare();
          };
        }
        
        set({ isScreenSharing: true });
      } catch (error) {
        console.error('Error starting screen share:', error);
      }
    }
  },

  toggleRecording: () => {
    set({ isRecording: !get().isRecording });
  },

  joinRoom: (roomId: string) => {
    set({ currentRoomId: roomId });
    // Initialize peer manager and WebRTC connections
    const peerManager = new PeerManager();
    set({ peerManager });
  },

  leaveRoom: () => {
    const { peerManager, localStream } = get();
    
    if (peerManager) {
      peerManager.destroy();
    }
    
    if (localStream) {
      localStream.getTracks().forEach(track => track.stop());
    }
    
    set({
      currentRoomId: null,
      peerManager: null,
      remoteStreams: new Map(),
      peerConnections: new Map(),
    });
  },

  addRemoteStream: (peerId: string, stream: MediaStream) => {
    const { remoteStreams } = get();
    remoteStreams.set(peerId, stream);
    set({ remoteStreams: new Map(remoteStreams) });
  },

  removeRemoteStream: (peerId: string) => {
    const { remoteStreams } = get();
    remoteStreams.delete(peerId);
    set({ remoteStreams: new Map(remoteStreams) });
  },

  cleanup: () => {
    const { localStream, peerManager } = get();
    
    if (localStream) {
      localStream.getTracks().forEach(track => track.stop());
    }
    
    if (peerManager) {
      peerManager.destroy();
    }
    
    set({
      localStream: null,
      remoteStreams: new Map(),
      peerConnections: new Map(),
      isAudioEnabled: true,
      isVideoEnabled: true,
      isScreenSharing: false,
      isRecording: false,
      currentRoomId: null,
      peerManager: null,
    });
  },
}));
```

## 5. Service Layer

### 5.1 API Service

```ts
// src/services/api/authApi.ts
import { User, LoginResponse, RegisterData } from '../../types/auth';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

class AuthService {
  private getAuthHeaders() {
    const token = localStorage.getItem('auth-token');
    return {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  async login(email: string, password: string): Promise<LoginResponse> {
    const response = await fetch(`${API_BASE_URL}/v1/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Login failed');
    }

    const data = await response.json();
    
    // Store token in localStorage
    localStorage.setItem('auth-token', data.data.token);
    
    return data.data;
  }

  async register(data: RegisterData): Promise<LoginResponse> {
    const response = await fetch(`${API_BASE_URL}/v1/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Registration failed');
    }

    const responseData = await response.json();
    
    // Store token in localStorage
    localStorage.setItem('auth-token', responseData.data.token);
    
    return responseData.data;
  }

  async logout(): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/v1/auth/logout`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Logout failed');
    }

    // Remove token from localStorage
    localStorage.removeItem('auth-token');
  }

  async refreshToken(): Promise<{ token: string }> {
    const response = await fetch(`${API_BASE_URL}/v1/auth/refresh`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Token refresh failed');
    }

    const data = await response.json();
    
    // Update token in localStorage
    localStorage.setItem('auth-token', data.data.token);
    
    return data.data;
  }

  async getProfile(): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/v1/users/profile`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to get profile');
    }

    const data = await response.json();
    return data.data;
  }

  async updateProfile(data: Partial<User>): Promise<User> {
    const response = await fetch(`${API_BASE_URL}/v1/users/profile`, {
      method: 'PUT',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to update profile');
    }

    const responseData = await response.json();
    return responseData.data;
  }
}

export const authService = new AuthService();
```

### 5.2 WebSocket Service

```ts
// src/services/websocket/WebSocketService.ts
import { EventEmitter } from 'events';

interface WebSocketMessage {
  type: string;
  roomId?: string;
  userId?: string;
  targetUserId?: string;
  payload?: any;
}

interface WebSocketOptions {
  url: string;
  reconnectAttempts?: number;
  reconnectInterval?: number;
}

export class WebSocketService extends EventEmitter {
  private socket: WebSocket | null = null;
  private url: string;
  private reconnectAttempts: number;
  private reconnectInterval: number;
  private currentReconnectAttempt: number = 0;
  private isConnected: boolean = false;
  private reconnectTimer: NodeJS.Timeout | null = null;

  constructor(options: WebSocketOptions) {
    super();
    this.url = options.url;
    this.reconnectAttempts = options.reconnectAttempts || 5;
    this.reconnectInterval = options.reconnectInterval || 3000;
  }

  connect(): void {
    try {
      this.socket = new WebSocket(this.url);
      this.setupEventListeners();
    } catch (error) {
      console.error('WebSocket connection error:', error);
      this.handleReconnect();
    }
  }

  private setupEventListeners(): void {
    if (!this.socket) return;

    this.socket.onopen = () => {
      console.log('WebSocket connected');
      this.isConnected = true;
      this.currentReconnectAttempt = 0;
      this.emit('connected');
    };

    this.socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        this.emit('message', message);
        
        // Emit specific events based on message type
        if (message.type) {
          this.emit(message.type, message);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    this.socket.onclose = (event) => {
      console.log('WebSocket disconnected:', event.code, event.reason);
      this.isConnected = false;
      this.emit('disconnected', event);
      this.handleReconnect();
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.emit('error', error);
    };
  }

  private handleReconnect(): void {
    if (this.currentReconnectAttempt < this.reconnectAttempts) {
      this.currentReconnectAttempt++;
      console.log(`Attempting to reconnect... (${this.currentReconnectAttempt}/${this.reconnectAttempts})`);
      
      this.reconnectTimer = setTimeout(() => {
        this.connect();
      }, this.reconnectInterval);
    } else {
      console.error('Max reconnection attempts reached');
      this.emit('maxReconnectReached');
    }
  }

  send(message: WebSocketMessage): void {
    if (this.socket && this.isConnected) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  joinRoom(roomId: string, userId: string, username: string): void {
    this.send({
      type: 'join_room',
      roomId,
      userId,
      payload: { username },
    });
  }

  leaveRoom(roomId: string, userId: string, username: string): void {
    this.send({
      type: 'user_left',
      roomId,
      userId,
      payload: { username },
    });
  }

  sendChatMessage(roomId: string, userId: string, username: string, message: string): void {
    this.send({
      type: 'chat',
      roomId,
      userId,
      payload: { message },
    });
  }

  sendWebRTCOffer(roomId: string, userId: string, targetUserId: string, sdp: string): void {
    this.send({
      type: 'offer',
      roomId,
      userId,
      targetUserId,
      payload: { sdp },
    });
  }

  sendWebRTCAnswer(roomId: string, userId: string, targetUserId: string, sdp: string): void {
    this.send({
      type: 'answer',
      roomId,
      userId,
      targetUserId,
      payload: { sdp },
    });
  }

  sendWebRTCIceCandidate(roomId: string, userId: string, targetUserId: string, candidate: any): void {
    this.send({
      type: 'ice-candidate',
      roomId,
      userId,
      targetUserId,
      payload: { candidate },
    });
  }

  sendMuteStatus(roomId: string, userId: string, muted: boolean): void {
    this.send({
      type: 'mute',
      roomId,
      userId,
      payload: { muted },
    });
  }

  sendVideoStatus(roomId: string, userId: string, videoEnabled: boolean): void {
    this.send({
      type: 'video',
      roomId,
      userId,
      payload: { videoEnabled },
    });
  }

  sendScreenShareStatus(roomId: string, userId: string, enabled: boolean): void {
    this.send({
      type: 'screen_share',
      roomId,
      userId,
      payload: { enabled },
    });
  }

  disconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }

    this.isConnected = false;
  }

  get connected(): boolean {
    return this.isConnected;
  }
}
```

### 5.3 WebRTC Service

```ts
// src/services/webrtc/PeerManager.ts
import SimplePeer from 'simple-peer';
import { EventEmitter } from 'events';

interface PeerManagerOptions {
  initiator: boolean;
  trickle: boolean;
  config?: RTCConfiguration;
}

export class PeerManager extends EventEmitter {
  private peer: SimplePeer.Instance | null = null;
  private stream: MediaStream | null = null;
  private options: PeerManagerOptions;

  constructor(options: PeerManagerOptions) {
    super();
    this.options = {
      trickle: true,
      config: {
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' },
        ],
      },
      ...options,
    };
  }

  init(stream: MediaStream): void {
    this.stream = stream;
    
    this.peer = new SimplePeer({
      initiator: this.options.initiator,
      trickle: this.options.trickle,
      config: this.options.config,
      streams: [stream],
    });

    this.setupEventListeners();
  }

  private setupEventListeners(): void {
    if (!this.peer) return;

    this.peer.on('signal', (data) => {
      this.emit('signal', data);
    });

    this.peer.on('connect', () => {
      console.log('Peer connected');
      this.emit('connect');
    });

    this.peer.on('data', (data) => {
      this.emit('data', data);
    });

    this.peer.on('stream', (stream) => {
      console.log('Remote stream received');
      this.emit('stream', stream);
    });

    this.peer.on('track', (track, stream) => {
      console.log('Remote track received');
      this.emit('track', track, stream);
    });

    this.peer.on('close', () => {
      console.log('Peer connection closed');
      this.emit('close');
    });

    this.peer.on('error', (error) => {
      console.error('Peer error:', error);
      this.emit('error', error);
    });
  }

  signal(data: any): void {
    if (this.peer) {
      this.peer.signal(data);
    }
  }

  send(data: string | Uint8Array): void {
    if (this.peer) {
      this.peer.send(data);
    }
  }

  addTrack(track: MediaStreamTrack, stream: MediaStream): void {
    if (this.peer) {
      this.peer.addTrack(track, stream);
    }
  }

  replaceTrack(oldTrack: MediaStreamTrack, newTrack: MediaStreamTrack): void {
    if (this.peer) {
      this.peer.replaceTrack(oldTrack, newTrack, this.stream!);
    }
  }

  removeTrack(track: MediaStreamTrack, stream: MediaStream): void {
    if (this.peer) {
      this.peer.removeTrack(track, stream);
    }
  }

  destroy(): void {
    if (this.peer) {
      this.peer.destroy();
      this.peer = null;
    }
    this.stream = null;
    this.removeAllListeners();
  }

  get connected(): boolean {
    return this.peer ? this.peer.connected : false;
  }
}
```

## 6. Page Components

### 6.1 Meeting Page

```tsx
// src/pages/Meeting/Meeting.tsx
import { h } from 'preact';
import { useEffect, useState, useRef } from 'preact/hooks';
import { useParams, useLocation } from 'preact-router';
import { VideoPlayer } from '../../components/meeting/VideoPlayer/VideoPlayer';
import { MeetingControls } from '../../components/meeting/MeetingControls/MeetingControls';
import { ParticipantList } from '../../components/meeting/ParticipantList/ParticipantList';
import { ChatBox } from '../../components/meeting/ChatBox/ChatBox';
import { Button } from '../../components/common/Button/Button';
import { useAuthStore } from '../../stores/authStore';
import { useWebRTCStore } from '../../stores/webrtcStore';
import { useWebSocketStore } from '../../stores/webSocketStore';
import { useRoomStore } from '../../stores/roomStore';
import './Meeting.css';

export function Meeting() {
  const { roomId } = useParams<{ roomId: string }>();
  const location = useLocation();
  const { user } = useAuthStore();
  const {
    localStream,
    remoteStreams,
    isAudioEnabled,
    isVideoEnabled,
    isScreenSharing,
    isRecording,
    initializeMedia,
    toggleAudio,
    toggleVideo,
    toggleScreenShare,
    toggleRecording,
    joinRoom,
    leaveRoom,
    addRemoteStream,
    removeRemoteStream,
  } = useWebRTCStore();
  
  const { connect, disconnect, sendMessage } = useWebSocketStore();
  const { currentRoom, joinRoom: joinRoomApi, leaveRoom: leaveRoomApi } = useRoomStore();
  
  const [showParticipants, setShowParticipants] = useState(false);
  const [showChat, setShowChat] = useState(false);
  const [isInitializing, setIsInitializing] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  const remoteVideosRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const initMeeting = async () => {
      try {
        setIsInitializing(true);
        setError(null);
        
        // Initialize media devices
        await initializeMedia();
        
        // Join room via API
        if (roomId && user) {
          await joinRoomApi(roomId, '');
          
          // Connect to WebSocket
          connect(roomId, user.id, user.firstName);
          
          // Join WebRTC room
          joinRoom(roomId);
        }
      } catch (err) {
        console.error('Error initializing meeting:', err);
        setError(err instanceof Error ? err.message : 'Failed to initialize meeting');
      } finally {
        setIsInitializing(false);
      }
    };

    initMeeting();

    return () => {
      // Cleanup on unmount
      leaveRoom();
      disconnect();
    };
  }, [roomId, user]);

  const handleToggleAudio = () => {
    toggleAudio();
    // Send mute status to other participants
    if (roomId && user) {
      sendMessage({
        type: 'mute',
        roomId,
        userId: user.id,
        payload: { muted: !isAudioEnabled },
      });
    }
  };

  const handleToggleVideo = () => {
    toggleVideo();
    // Send video status to other participants
    if (roomId && user) {
      sendMessage({
        type: 'video',
        roomId,
        userId: user.id,
        payload: { videoEnabled: !isVideoEnabled },
      });
    }
  };

  const handleToggleScreenShare = async () => {
    await toggleScreenShare();
    // Send screen share status to other participants
    if (roomId && user) {
      sendMessage({
        type: 'screen_share',
        roomId,
        userId: user.id,
        payload: { enabled: !isScreenSharing },
      });
    }
  };

  const handleLeaveMeeting = async () => {
    try {
      if (roomId && user) {
        await leaveRoomApi(roomId);
      }
      
      leaveRoom();
      disconnect();
      
      // Navigate to dashboard
      // navigate('/dashboard');
    } catch (err) {
      console.error('Error leaving meeting:', err);
      setError(err instanceof Error ? err.message : 'Failed to leave meeting');
    }
  };

  const remoteStreamsArray = Array.from(remoteStreams.values());
  const participantCount = remoteStreamsArray.length + 1; // +1 for local user

  if (isInitializing) {
    return (
      <div className="meeting-loading">
        <div className="loading-spinner"></div>
        <p>Initializing meeting...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="meeting-error">
        <h2>Error</h2>
        <p>{error}</p>
        <Button onClick={handleLeaveMeeting}>Go to Dashboard</Button>
      </div>
    );
  }

  return (
    <div className="meeting-container">
      <div className="meeting-header">
        <h1>{currentRoom?.name || 'Meeting'}</h1>
        <div className="meeting-info">
          <span className="room-id">Room ID: {roomId}</span>
          <span className="participant-count">{participantCount} participants</span>
        </div>
      </div>
      
      <div className="meeting-content">
        <div className="videos-container" ref={remoteVideosRef}>
          {/* Local video */}
          {localStream && (
            <div className="video-container local-video">
              <VideoPlayer
                stream={localStream}
                isLocal={true}
                username={user?.firstName || 'You'}
                isAudioEnabled={isAudioEnabled}
                isVideoEnabled={isVideoEnabled}
                isScreenSharing={isScreenSharing}
                onToggleAudio={handleToggleAudio}
                onToggleVideo={handleToggleVideo}
              />
            </div>
          )}
          
          {/* Remote videos */}
          {remoteStreamsArray.map((stream, index) => (
            <div key={index} className="video-container remote-video">
              <VideoPlayer
                stream={stream}
                isLocal={false}
                username={`Participant ${index + 1}`}
              />
            </div>
          ))}
        </div>
        
        {/* Participant List */}
        {showParticipants && (
          <div className="participants-panel">
            <ParticipantList
              onClose={() => setShowParticipants(false)}
              participantCount={participantCount}
            />
          </div>
        )}
        
        {/* Chat Box */}
        {showChat && (
          <div className="chat-panel">
            <ChatBox
              onClose={() => setShowChat(false)}
              onSendMessage={(message) => {
                if (roomId && user) {
                  sendMessage({
                    type: 'chat',
                    roomId,
                    userId: user.id,
                    payload: { message },
                  });
                }
              }}
            />
          </div>
        )}
      </div>
      
      {/* Meeting Controls */}
      <MeetingControls
        isAudioEnabled={isAudioEnabled}
        isVideoEnabled={isVideoEnabled}
        isScreenSharing={isScreenSharing}
        isRecording={isRecording}
        participantCount={participantCount}
        onToggleAudio={handleToggleAudio}
        onToggleVideo={handleToggleVideo}
        onToggleScreenShare={handleToggleScreenShare}
        onToggleRecording={toggleRecording}
        onShowParticipants={() => {
          setShowParticipants(true);
          setShowChat(false);
        }}
        onShowChat={() => {
          setShowChat(true);
          setShowParticipants(false);
        }}
        onLeaveMeeting={handleLeaveMeeting}
      />
    </div>
  );
}
```

```css
/* src/pages/Meeting/Meeting.css */
.meeting-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--color-secondary-900);
  color: white;
  overflow: hidden;
}

.meeting-header {
  padding: 1rem 1.5rem;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 10;
}

.meeting-header h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.meeting-info {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.room-id, .participant-count {
  background-color: rgba(255, 255, 255, 0.1);
  padding: 0.375rem 0.75rem;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
}

.meeting-content {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.videos-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
  padding: 1rem;
  height: 100%;
  overflow-y: auto;
}

.video-container {
  position: relative;
  border-radius: var(--radius-lg);
  overflow: hidden;
  background-color: var(--color-secondary-800);
}

.local-video {
  position: absolute;
  bottom: 1rem;
  right: 1rem;
  width: 300px;
  height: 200px;
  z-index: 5;
  box-shadow: var(--shadow-lg);
}

.participants-panel, .chat-panel {
  position: absolute;
  top: 0;
  right: 0;
  width: 320px;
  height: 100%;
  background-color: white;
  box-shadow: var(--shadow-xl);
  z-index: 20;
  transform: translateX(100%);
  transition: transform var(--transition-normal);
}

.participants-panel.open, .chat-panel.open {
  transform: translateX(0);
}

.meeting-loading, .meeting-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: var(--color-secondary-900);
  color: white;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top: 4px solid var(--color-primary-500);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.meeting-error {
  background-color: var(--color-error);
}

.meeting-error h2 {
  margin-bottom: 1rem;
}

@media (max-width: 768px) {
  .meeting-header {
    padding: 1rem;
  }
  
  .meeting-header h1 {
    font-size: 1.25rem;
  }
  
  .meeting-info {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .videos-container {
    grid-template-columns: 1fr;
    padding: 0.5rem;
  }
  
  .local-video {
    width: 120px;
    height: 90px;
    bottom: 0.5rem;
    right: 0.5rem;
  }
  
  .participants-panel, .chat-panel {
    width: 100%;
  }
}
```

## 7. Kesimpulan

Desain arsitektur frontend untuk aplikasi WebRTC meeting dengan Preact dan tema hijau-abu-abu ini mencakup:

1. **Arsitektur Komponen** dengan pemisahan yang jelas antara layout, page, dan reusable components
2. **State Management** menggunakan Zustand untuk performa yang optimal dan developer experience yang baik
3. **Tema Hijau dan Abu-abu** yang memberikan kesan profesional dan modern dengan palet warna yang konsisten
4. **Service Layer** untuk mengelola komunikasi dengan backend, WebSocket, dan WebRTC
5. **Responsive Design** yang berfungsi baik di desktop maupun mobile

Dengan arsitektur ini, aplikasi frontend akan memiliki:
- Performa yang optimal dengan Preact yang lightweight
- Developer experience yang baik dengan TypeScript dan tooling modern
- Maintainability yang tinggi dengan struktur yang terorganisir
- User experience yang baik dengan tema yang konsisten dan UI yang intuitif
- Scalability yang baik dengan state management dan service layer yang terpisah