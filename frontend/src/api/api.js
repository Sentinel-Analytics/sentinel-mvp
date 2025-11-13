const API_URL = "https://api-sentinel.getmusterup.com"

const request = async (endpoint, options = {}) => {
  const defaultOptions = {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  }

  try {
    const response = await fetch(`${API_URL}${endpoint}`, { ...defaultOptions, ...options })
    const data = await response.json().catch(() => ({}))

    if (!response.ok) {
      throw new Error(data.error || `HTTP error! status: ${response.status}`)
    }
    return data
  } catch (error) {
    console.error(`API request to ${endpoint} failed:`, error)
    throw error
  }
}

export const api = {
  login: (email, password) => {
    const formData = new URLSearchParams()
    formData.append("email", email)
    formData.append("password", password)
    return request("/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: formData,
    })
  },

  signup: (email, password) => {
    const formData = new URLSearchParams()
    formData.append("email", email)
    formData.append("password", password)
    return request("/auth/signup", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: formData,
    })
  },

  logout: () => request("/logout"),
  getSites: () => request("/api/sites/"),
  addSite: (name) =>
    request("/api/sites/", {
      method: "POST",
      body: JSON.stringify({ name }),
    }),
  getDashboardStats: (siteId, days) => request(`/api/dashboard?siteId=${siteId}&days=${days}`),
}
