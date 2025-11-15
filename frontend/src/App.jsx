import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom"
import Dashboard from "./pages/Dashboard"
import Login from "./pages/Login"
import Signup from "./pages/SignUp"
import SessionReplay from "./pages/SessionReplay"
import FunnelsPage from "./pages/FunnelsPage"
import FirewallPage from "./pages/FirewallPage"

const ProtectedRoute = ({ children }) => {
  // The Dashboard component now handles its own authentication check
  return children
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />
        <Route path="/session-replay" element={<SessionReplay />} />
        <Route path="/funnels" element={<FunnelsPage />} />
        <Route path="/firewall" element={<FirewallPage />} />
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Router>
  )
}

export default App
