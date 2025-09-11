Sentinel
A privacy-first, open-source web analytics platform. Get the insights you need without compromising user privacy.

Sentinel provides simple, beautiful, and actionable website analytics. It's lightweight, easy to self-host, and gives you 100% data ownership. Built for developers, marketers, and businesses who value privacy.

✨ Key Features (MVP v1)
👤 Multi-User & Multi-Site: Manage multiple websites from a single, secure account.

📊 Essential Analytics: Get all the key metrics: Unique Visitors, Page Views, Bounce Rate, Average Visit Time, Top Pages, and Top Referrers.

🌍 Visitor Insights: Understand your audience with reports on their Browser, Operating System, and Country.

📅 Date Filtering: Analyze your traffic over the last 24 hours, 7 days, or 30 days.

🚀 Incredibly Lightweight: The tracking script is tiny (<2KB) and won't slow down your site.

🔒 100% Self-Hosted: Run it on your own infrastructure with Docker. You own your data, always.

🚀 Quick Start (One-Click Deployment)
Get your own Sentinel instance running in under 5 minutes.

Prerequisites:

A server (VPS or home server) with docker and docker-compose installed.

The GeoLite2-Country.mmdb file downloaded from MaxMind.

1. Clone the Repository

git clone [https://github.com/your-username/sentinel.git](https://github.com/your-username/sentinel.git)
cd sentinel

2. Place GeoIP Database
Place your downloaded GeoLite2-Country.mmdb file in the root of this project directory.

3. Configure Environment
Create a file named .env in the project root. This will hold your database connection details.

# .env
DATABASE_URL=postgres://sentinel:your_super_secret_password@db:5432/sentinel?sslmode=disable

Then, update your docker-compose.yml to match the password you chose for the POSTGRES_PASSWORD variable.

4. Run the Application

docker compose up --build -d

5. You're Live!
Your Sentinel instance is now running!

Create Your Admin Account: Go to http://your-server-ip:8000/signup

Log In: Go to http://your-server-ip:8000/login

🔮 What's Next? The Road to v2
This MVP is just the beginning. Our vision for Sentinel is to be a complete website intelligence platform. Here's a sneak peek at what we're working on for the next major release:

🛡️ Traffic Quality Score & Bot Detection: Go beyond simple numbers and understand how much of your traffic is from real humans.

🔥 Real-Time Firewall: Act on the data by creating simple rules to block malicious traffic from specific countries, data centers (ASNs), and known bots.

⚡ Core Web Vitals: Monitor your site's performance and its impact on user experience directly within your Sentinel dashboard.

🤝 Community & Feedback
Sentinel is built for the community, by the community. Your feedback is invaluable in shaping the future of this project.

🐛 Report a Bug: Found an issue? Please open an issue and let us know.

💡 Request a Feature: Have a great idea? Start a discussion in the feature requests section.

💬 Join the Conversation: (Coming Soon) We'll be launching a Discord server for real-time chat, support, and announcements.

❤️ Support the Project
If you find Sentinel useful and want to support its continued development, please consider:

⭐ Starring the project on GitHub. This is the easiest and most impactful way to show your support!

☕ Buy me a coffee. A small donation helps fuel development and is greatly appreciated.

📜 License
This project is licensed under the MIT License.