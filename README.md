# Sentinel

![GitHub Repo stars](https://img.shields.io/github/stars/your-username/sentinel?style=social)
![GitHub issues](https://img.shields.io/github/issues/your-username/sentinel)
![GitHub license](https://img.shields.io/github/license/your-username/sentinel)
![Docker](https://img.shields.io/badge/docker-ready-blue)

**A privacy-first, open-source web analytics platform.**  
Get the insights you need without compromising user privacy.

Sentinel provides simple, beautiful, and actionable website analytics. It's lightweight, easy to self-host, and gives you **100% data ownership**. Built for developers, marketers, and businesses who value privacy.

---

## ✨ Key Features (MVP v1)

- 👤 **Multi-User & Multi-Site**: Manage multiple websites from a single, secure account.
- 📊 **Essential Analytics**: Unique Visitors, Page Views, Bounce Rate, Average Visit Time, Top Pages, and Top Referrers.
- 🌍 **Visitor Insights**: Reports on Browser, Operating System, and Country.
- 📅 **Date Filtering**: Analyze traffic over the last 24 hours, 7 days, or 30 days.
- 🚀 **Incredibly Lightweight**: Tracking script is tiny (<2KB) and won’t slow down your site.
- 🔒 **100% Self-Hosted**: Run on your own infrastructure with Docker. You own your data, always.

---

## 🚀 Quick Start (One-Click Deployment)

Get your own **Sentinel** instance running in under 5 minutes.

### Prerequisites
- A server (VPS or home server) with **docker** and **docker-compose** installed.  
- The **GeoLite2-Country.mmdb** file downloaded from MaxMind.

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/sentinel.git
cd sentinel
```

### 2. Place GeoIP Database
Place your downloaded `GeoLite2-Country.mmdb` file in the **root** of this project directory.

### 3. Configure Environment
Create a `.env` file in the project root with your database connection details:

```env
DATABASE_URL=postgres://sentinel:your_super_secret_password@db:5432/sentinel?sslmode=disable
```

Update your `docker-compose.yml` to match the password you chose for the `POSTGRES_PASSWORD` variable.

### 4. Run the Application
```bash
docker compose up --build -d
```

### 5. You're Live! 🎉
- **Create Your Admin Account:**  
  Go to [http://your-server-ip:8000/signup](http://your-server-ip:8000/signup)

- **Log In:**  
  Go to [http://your-server-ip:8000/login](http://your-server-ip:8000/login)

---

## 🔮 What's Next? The Road to v2

This MVP is just the beginning. Our vision for **Sentinel** is to be a complete website intelligence platform. Coming in **v2**:

- 🛡️ **Traffic Quality Score & Bot Detection** – Go beyond simple numbers and detect real human traffic.
- 🔥 **Real-Time Firewall** – Block malicious traffic from specific countries, data centers (ASNs), and known bots.
- ⚡ **Core Web Vitals** – Monitor site performance and its impact on user experience directly in your dashboard.

---

## 🤝 Community & Feedback

Sentinel is built **for the community, by the community**. Your feedback is invaluable in shaping the future of this project.

- 🐛 **Report a Bug**: Open an issue and let us know.
- 💡 **Request a Feature**: Start a discussion in the feature requests section.
- 💬 **Join the Conversation**: (Coming Soon) Discord server for chat, support, and announcements.

---

## ❤️ Support the Project

If you find Sentinel useful and want to support its continued development:

- ⭐ Star the project on GitHub — the easiest and most impactful way to show support!  


---

## 📜 License

This project is licensed under the **MIT License**.
