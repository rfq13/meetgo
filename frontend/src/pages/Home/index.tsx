import { h } from 'preact'
import { Link } from 'preact-router'
import { Button } from '../../components/common/Button/Button'
import './Home.css'

export default function Home() {
  return (
    <div className="home">
      {/* Hero Section */}
      <section className="hero">
        <div className="container">
          <div className="hero-content">
            <h1 className="hero-title">
              MeetGreen
              <span className="hero-subtitle">Video Conferencing with Nature</span>
            </h1>
            <p className="hero-description">
              Platform video conferencing modern dengan tema hijau yang menenangkan. 
              Bergabunglah dalam rapat online dengan pengalaman yang menyenangkan dan profesional.
            </p>
            <div className="hero-actions">
              <Link href="/register">
                <Button variant="primary" size="lg" className="mr-4">
                  Mulai Gratis
                </Button>
              </Link>
              <Link href="/login">
                <Button variant="outline" size="lg">
                  Masuk
                </Button>
              </Link>
            </div>
          </div>
          <div className="hero-image">
            <div className="hero-placeholder">
              <div className="video-grid">
                <div className="video-tile"></div>
                <div className="video-tile"></div>
                <div className="video-tile"></div>
                <div className="video-tile"></div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="features">
        <div className="container">
          <h2 className="section-title">Fitur Unggulan</h2>
          <div className="features-grid">
            <div className="feature-card">
              <div className="feature-icon">
                <div className="icon-video"></div>
              </div>
              <h3>Video HD Quality</h3>
              <p>Nikmati kualitas video tinggi dengan teknologi WebRTC terkini</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">
                <div className="icon-screen"></div>
              </div>
              <h3>Screen Sharing</h3>
              <p>Bagikan layar Anda dengan mudah untuk presentasi yang efektif</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">
                <div className="icon-chat"></div>
              </div>
              <h3>Real-time Chat</h3>
              <p>Komunikasi dengan peserta melalui chat real-time yang responsif</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">
                <div className="icon-record"></div>
              </div>
              <h3>Recording</h3>
              <p>Rekam rapat penting Anda untuk dokumentasi dan referensi</p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="cta">
        <div className="container">
          <div className="cta-content">
            <h2>Siap Memulai Rapat Online?</h2>
            <p>Bergabung dengan ribuan pengguna yang telah mempercayai MeetGreen</p>
            <Link href="/register">
              <Button variant="primary" size="lg">
                Daftar Sekarang
              </Button>
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}