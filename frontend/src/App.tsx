import { h } from 'preact'
import { Router } from 'preact-router'
import { Route } from 'preact-router'

// Pages
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Meeting from './pages/Meeting'
import Profile from './pages/Profile'
import Contacts from './pages/Contacts'
import Rooms from './pages/Rooms'
import Schedule from './pages/Schedule'
import Settings from './pages/Settings'
import Notifications from './pages/Notifications'
import Recordings from './pages/Recordings'
import NotFound from './pages/NotFound'

function App() {
  return (
    <div className="min-h-screen bg-secondary-50">
      <Router>
        <Route path="/" component={Home} />
        <Route path="/login" component={Login} />
        <Route path="/register" component={Register} />
        <Route path="/dashboard" component={Dashboard} />
        <Route path="/meeting/:roomId" component={Meeting} />
        <Route path="/profile" component={Profile} />
        <Route path="/contacts" component={Contacts} />
        <Route path="/rooms" component={Rooms} />
        <Route path="/schedule" component={Schedule} />
        <Route path="/settings" component={Settings} />
        <Route path="/notifications" component={Notifications} />
        <Route path="/recordings" component={Recordings} />
        <Route default component={NotFound} />
      </Router>
    </div>
  )
}

export default App
