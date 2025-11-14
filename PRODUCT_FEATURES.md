# Sentinel: Website Intelligence & Security Platform - Feature Overview

## 1. Introduction: Beyond Simple Analytics

Sentinel is an open-source, self-hostable platform that redefines website analytics. In a world where businesses are forced to choose between invasive, complex tools and overly simplistic ones, Sentinel provides a third option: a powerful, privacy-first platform that offers deep insights while simultaneously protecting your website from automated threats.

Our vision is to provide **Website Intelligence**, not just analytics. This means empowering developers, marketers, and founders with clear, actionable data, while giving them the tools to ensure the quality and security of their traffic.

## 2. Core Features: The Foundation of Insight

At its core, Sentinel provides a comprehensive and easy-to-understand analytics dashboard.

*   **Privacy-First Tracking:** Our tracking script is lightweight (<2KB), cookieless, and fully compliant with GDPR/CCPA, ensuring you can gather insights without compromising your users' privacy.
*   **Essential Analytics Dashboard:** Get an at-a-glance overview of your website's performance with key metrics:
    *   **Total Page Views:** The total number of pages viewed.
    *   **Unique Visitors:** The number of individual visitors to your site.
    *   **Bounce Rate:** The percentage of visitors who leave after viewing only one page.
    *   **Average Visit Duration:** The average amount of time users spend on your site.
*   **Detailed Reports:** Dive deeper into your data with reports on:
    *   **Top Pages:** See which pages are most popular.
    *   **Top Referrers:** Understand where your traffic is coming from.
    *   **Audience Insights:** Break down your audience by Country, Browser, and Operating System.
*   **Simple Site Management:** Easily manage multiple websites from a single, intuitive interface. Each site gets its own unique tracking ID and dashboard.

## 3. Differentiator Features: The Sentinel Advantage

This is where Sentinel truly shines. We go beyond traditional analytics to provide features that give you a competitive edge.

### 3.1. Advanced Security & Threat Detection

A significant portion of web traffic is automated, consisting of bots, scrapers, and ad-fraud clicks. This traffic pollutes your data and can cost you money. Sentinel gives you the tools to fight back.

*   **Traffic Quality Score:** Our signature feature, this single KPI (0-100) gives you an immediate understanding of your traffic's health. It's a composite score based on bot detection, bounce rate anomalies, and other patterns, allowing you to see at a glance if you're under attack.
*   **Advanced Bot Detection:** We use a multi-layered approach to identify non-human traffic, including:
    *   **User-Agent Heuristics:** We check against a list of known bot signatures.
    *   **Data Center IP Blocking:** We identify and penalize traffic coming from known data centers (ASNs), a common source of bot traffic.
*   **The Sentinel Firewall:** Take control of your traffic with our easy-to-use firewall. You can create rules to block traffic by:
    *   **IP Address or CIDR Range**
    *   **Country**
    *   **ASN (Data Center)**

### 3.2. Deep User & Performance Insights

Understand not just *what* your users are doing, but *how* and *why*.

*   **Session Replay:** Go beyond the numbers and watch recordings of real user sessions. See every mouse movement, click, and scroll, allowing you to identify user friction points, debug issues, and optimize your user experience.
*   **Web Vitals Monitoring:** Website performance is critical for user experience and SEO. Sentinel automatically tracks Core Web Vitals:
    *   **LCP (Largest Contentful Paint):** Measures loading performance.
    *   **CLS (Cumulative Layout Shift):** Measures visual stability.
    *   **FID (First Input Delay):** Measures interactivity.
    These metrics are displayed on your dashboard, giving you a clear view of your site's performance.

### 3.3. Conversion Tracking

*   **Funnels & Goals:** Understand how users navigate through your site and where they drop off. With our Funnels feature, you can define multi-step conversion funnels (e.g., `Homepage -> Pricing -> Checkout`) and visualize the conversion rate at each step. This is invaluable for optimizing your marketing campaigns and user flows.

## 4. Deployment: Own Your Data

Sentinel is designed to be easily self-hostable. We provide a production-ready `docker-compose.yml` file that allows you to get Sentinel up and running on your own infrastructure in under five minutes. By self-hosting, you maintain complete ownership and control of your data.

## 5. Getting Started

1.  Clone the Sentinel repository.
2.  Run `docker-compose up -d`.
3.  Navigate to `http://localhost:5173` to create your account.
4.  Add your first site and copy the tracking script to your website's HTML.
5.  Start gaining insights!

---

Thank you for choosing Sentinel. We're excited to help you unlock the true potential of your website's data.
