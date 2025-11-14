"use client"

import { useState, useEffect } from "react"
import { Eye, Users, TrendingDown, Clock, Plus, LogOut, Copy, Check, Trash2 } from "lucide-react"
import { api } from "../api"
import Logo from "../components/Logo"
import StatCard from "../components/StatCard"
import BarChart from "../components/BarChart"
import LineChart from "../components/LineChart"
import DoughnutChart from "../components/DoughnutChart"

const Dashboard = () => {
  const [sites, setSites] = useState([])
  const [selectedSite, setSelectedSite] = useState(null)
  const [dashboardData, setDashboardData] = useState(null)
  const [loading, setLoading] = useState(false)
  const [newSiteName, setNewSiteName] = useState("")
  const [copied, setCopied] = useState(false)
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [authLoading, setAuthLoading] = useState(true)

  useEffect(() => {
    const checkAuth = async () => {
      try {
        // Try to fetch sites to verify authentication
        await api.getSites()
        setIsAuthenticated(true)
        fetchSites()
      } catch (error) {
        console.error("Authentication check failed:", error)
        setIsAuthenticated(false)
        // Redirect to login if not authenticated
        window.location.href = "/login"
      } finally {
        setAuthLoading(false)
      }
    }
    checkAuth()
  }, [])

  const fetchSites = async () => {
    try {
      const sitesData = await api.getSites()
      setSites(sitesData || [])
      if (sitesData && sitesData.length > 0) {
        setSelectedSite(sitesData[0])
      }
    } catch (error) {
      console.error("Failed to fetch sites:", error)
      setSites([])
    }
  }

  const fetchDashboardData = async (siteId) => {
    setLoading(true)
    try {
      const data = await api.getDashboardStats(siteId, 30)
      setDashboardData(data)
    } catch (error) {
      console.error("Failed to fetch dashboard data:", error)
      setDashboardData(null)
    } finally {
      setLoading(false)
    }
  }

  const addSite = async (e) => {
    e.preventDefault()
    if (!newSiteName.trim()) return

    try {
      const newSite = await api.addSite(newSiteName)
      setNewSiteName("")
      fetchSites()
      if (newSite) {
        setSelectedSite(newSite)
      }
    } catch (error) {
      console.error("Failed to add site:", error)
    }
  }

  const handleDeleteSite = async (siteId) => {
    if (!window.confirm("Are you sure you want to delete this site? This action cannot be undone.")) {
      return
    }

    try {
      await api.deleteSite(siteId)
      // Refetch sites after deletion
      const updatedSites = await api.getSites()
      setSites(updatedSites || [])

      // If the deleted site was the selected one, clear the selection
      if (selectedSite && selectedSite.id === siteId) {
        setSelectedSite(updatedSites && updatedSites.length > 0 ? updatedSites[0] : null)
      }
    } catch (error) {
      console.error("Failed to delete site:", error)
    }
  }

  const handleLogout = async () => {
    try {
      await api.logout()
      setIsAuthenticated(false)
      window.location.href = "/login"
    } catch (error) {
      console.error("Logout failed:", error)
      // Force logout even if API call fails
      setIsAuthenticated(false)
      window.location.href = "/login"
    }
  }

  const copyTrackingScript = async () => {
    if (!selectedSite) return

    const trackingScript = `<script src="https://api-sentinel.getmusterup.com/static/tracker-v3.js" data-site-id="${selectedSite.id}"></script>`

    try {
      await navigator.clipboard.writeText(trackingScript)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (error) {
      console.error("Failed to copy:", error)
    }
  }

  useEffect(() => {
    if (selectedSite) {
      fetchDashboardData(selectedSite.id)
    }
  }, [selectedSite])

  if (authLoading) {
    return (
      <div className="min-h-screen bg-slate-900 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return null // Will redirect to login
  }

  return (
    <div
      className="min-h-screen bg-slate-900"
      style={{
        backgroundImage: `radial-gradient(circle at 1px 1px, rgba(148, 163, 184, 0.15) 1px, transparent 0)`,
        backgroundSize: "20px 20px",
      }}
    >
      <div className="flex">
        {/* Sidebar */}
        <div className="w-80 bg-slate-800 border-r border-slate-700 min-h-screen p-6">
          <Logo className="mb-8" />

          {/* Sites List */}
          <div className="mb-8">
            <h3 className="text-slate-200 font-semibold mb-4">Your Sites</h3>
            <div className="space-y-2">
              {sites.length > 0 ? (
                sites.map((site) => (
                  <div key={site.id} className="flex items-center justify-between">
                    <button
                      onClick={() => setSelectedSite(site)}
                      className={`w-full text-left p-3 rounded-lg transition-colors ${
                        selectedSite?.id === site.id ? "bg-indigo-600 text-white" : "text-slate-300 hover:bg-slate-700"
                      }`}
                    >
                      {site.name}
                    </button>
                    <button
                      onClick={() => handleDeleteSite(site.id)}
                      className="p-2 text-slate-400 hover:text-red-500 transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                ))
              ) : (
                <div className="text-slate-400 text-sm p-3">No sites added yet. Add your first site below.</div>
              )}
            </div>
          </div>

          {/* Add Site Form */}
          <div className="mb-8">
            <h3 className="text-slate-200 font-semibold mb-4">Add New Site</h3>
            <form onSubmit={addSite} className="space-y-3">
              <input
                type="text"
                value={newSiteName}
                onChange={(e) => setNewSiteName(e.target.value)}
                placeholder="Site name"
                className="w-full p-3 bg-slate-700 border border-slate-600 rounded-lg text-slate-200 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
              <button
                type="submit"
                className="w-full flex items-center justify-center space-x-2 p-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
              >
                <Plus className="w-4 h-4" />
                <span>Add Site</span>
              </button>
            </form>
          </div>

          {/* Logout Button */}
          <button
            onClick={handleLogout}
            className="w-full flex items-center justify-center space-x-2 p-3 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg transition-colors"
          >
            <LogOut className="w-4 h-4" />
            <span>Logout</span>
          </button>
        </div>

        {/* Main Content */}
        <div className="flex-1 p-8">
          {selectedSite ? (
            <>
              <div className="mb-8">
                <h1 className="text-3xl font-bold text-slate-200 mb-2">Dashboard for {selectedSite.name}</h1>
                <p className="text-slate-400">Analytics and insights for your website</p>
              </div>

              <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 mb-8">
                <h3 className="text-slate-200 font-semibold mb-3">Tracking Script</h3>
                <p className="text-slate-400 text-sm mb-4">
                  Add this code to your website's HTML to start tracking analytics
                </p>
                <div className="relative">
                  <pre className="bg-slate-900 border border-slate-600 rounded-lg p-4 text-sm text-slate-300 overflow-x-auto">
                    <code>{`<script src="https://api-sentinel.getmusterup.com/static/tracker-v3.js" data-site-id="${selectedSite.id}"></script>`}</code>
                  </pre>
                  <button
                    onClick={copyTrackingScript}
                    className="absolute top-3 right-3 flex items-center space-x-2 px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white text-sm rounded-md transition-colors"
                  >
                    {copied ? (
                      <>
                        <Check className="w-4 h-4" />
                        <span>Copied!</span>
                      </>
                    ) : (
                      <>
                        <Copy className="w-4 h-4" />
                        <span>Copy</span>
                      </>
                    )}
                  </button>
                </div>
              </div>

              {dashboardData ? (
                <>
                  {/* Stats Cards */}
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                    <StatCard
                      title="Total Views"
                      value={dashboardData.totalViews?.toLocaleString() || "0"}
                      icon={Eye}
                      change={12.5}
                      changeType="positive"
                    />
                    <StatCard
                      title="Unique Visitors"
                      value={dashboardData.uniqueVisitors?.toLocaleString() || "0"}
                      icon={Users}
                      change={8.2}
                      changeType="positive"
                    />
                    <StatCard
                      title="Bounce Rate"
                      value={`${dashboardData.bounceRate?.toFixed(1) || 0}%`}
                      icon={TrendingDown}
                      change={-2.1}
                      changeType="positive"
                    />
                    <StatCard
                      title="Avg. Visit Time"
                      value={dashboardData.avgVisitTime || "0s"}
                      icon={Clock}
                      change={15.3}
                      changeType="positive"
                    />
                  </div>

                  {/* Charts Grid */}
                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                    <LineChart
                      title="Daily Visitors"
                      data={[]}
                      labels={[]}
                    />
                    <BarChart
                      title="Top Pages"
                      data={dashboardData.topPages?.map(p => p.count) || []}
                      labels={dashboardData.topPages?.map(p => p.value) || []}
                    />
                  </div>

                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <DoughnutChart
                      title="Visitors by Country"
                      data={dashboardData.topCountries?.map(c => c.count) || []}
                      labels={dashboardData.topCountries?.map(c => c.value) || []}
                    />
                    <DoughnutChart
                      title="Visitors by OS"
                      data={dashboardData.topOS?.map(os => os.count) || []}
                      labels={dashboardData.topOS?.map(os => os.value) || []}
                    />
                  </div>
                </>
              ) : loading ? (
                <div className="flex items-center justify-center h-64">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
                </div>
              ) : (
                <div className="flex items-center justify-center h-64">
                  <div className="text-center">
                    <div className="text-slate-400 text-lg mb-2">No analytics data available</div>
                    <div className="text-slate-500 text-sm">
                      Install the tracking script above and wait for visitors to see analytics data here.
                    </div>
                  </div>
                </div>
              )}
            </>
          ) : (
            <div className="flex items-center justify-center h-64">
              <div className="text-center">
                <div className="text-slate-400 text-lg mb-2">No site selected</div>
                <div className="text-slate-500 text-sm">Add a site from the sidebar to view analytics data.</div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Dashboard
