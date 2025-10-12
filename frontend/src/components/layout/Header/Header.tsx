import { h } from 'preact'
import { Link } from 'preact-router/match'
import { Button } from '../../common/Button/Button'
import './Header.css'

export function Header() {
  return (
    <header className="header">
      <div className="container">
        <div className="header-content">
          <Link href="/" className="logo">
            <div className="logo-icon">
              <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
                <circle cx="16" cy="16" r="14" stroke="currentColor" strokeWidth="2" />
                <path d="M10 16L14 20L22 12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </div>
            <span className="logo-text">WebRTC Meeting</span>
          </Link>
          
          <nav className="nav">
            <Link href="/" className="nav-link">Beranda</Link>
            <Link href="/features" className="nav-link">Fitur</Link>
            <Link href="/pricing" className="nav-link">Harga</Link>
            <Link href="/about" className="nav-link">Tentang</Link>
          </nav>
          
          <div className="header-actions">
            <Button variant="ghost" size="sm">
              Masuk
            </Button>
            <Button variant="primary" size="sm">
              Daftar
            </Button>
          </div>
        </div>
      </div>
    </header>
  )
}