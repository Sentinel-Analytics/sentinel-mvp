import axios from "axios";

const API_URL = "https://api-sentinel.getmusterup.com/";

const axiosInstance = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

const request = async (config) => {
  try {
    const response = await axiosInstance(config);
    return response.data;
  } catch (error) {
    console.error(`API request to ${config.url} failed:`, error);
    throw error.response ? error.response.data : error;
  }
};

export const api = {
  login: (email, password) => {
    return request({
      url: "/auth/login",
      method: "POST",
      data: { email, password },
    });
  },

  signup: (email, password) => {
    return request({
      url: "/auth/signup",
      method: "POST",
      data: { email, password },
    });
  },

  logout: () => request({ url: "/logout" }),
  getSites: () => request({ url: "/api/sites/" }),
  addSite: (name) =>
    request({
      url: "/api/sites/",
      method: "POST",
      data: { name },
    }),
  deleteSite: (id) =>
    request({
      url: "/api/sites/",
      method: "DELETE",
      data: { id },
    }),
  getDashboardStats: (siteId, days) =>
    request({
      url: `/api/dashboard?siteId=${siteId}&days=${days}`,
    }),
  getSessionEvents: (siteId, sessionId) =>
    request({
      url: `/api/session/events?siteId=${siteId}&sessionId=${sessionId}`,
    }),
  listSessions: (siteId) =>
    request({
      url: `/api/sessions?siteId=${siteId}`,
    }),
  listFunnels: (siteId) =>
    request({
      url: `/api/funnels/?siteId=${siteId}`,
    }),
  createFunnel: (funnel) =>
    request({
      url: "/api/funnels/",
      method: "POST",
      data: funnel,
    }),
  updateFunnel: (funnel) =>
    request({
      url: "/api/funnels/",
      method: "PUT",
      data: funnel,
    }),
  deleteFunnel: (id) =>
    request({
      url: `/api/funnels/?id=${id}`,
      method: "DELETE",
    }),
};
