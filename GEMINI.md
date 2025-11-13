Sentinel – Product Requirements Document (PRD)
Version: 2.0 Date: September 19, 2025 Purpose: To define the "what" of Sentinel — its vision, user personas, features, and functional requirements — to guide development toward building a market-leading analytics and security platform.
1.0 Product Vision & Overview
The web analytics landscape is broken. Businesses are forced to choose between invasive, complex tools like Google Analytics that treat their users' data as a product, or overly simplistic tools that lack actionable insights. Meanwhile, a growing portion of web traffic consists of bots, scrapers, and ad-fraud clicks, polluting data and costing businesses money.
Sentinel's Vision: To be the first open-source platform that provides Website Intelligence, not just analytics. We empower developers, marketers, and founders to get the clear, privacy-first insights they need while simultaneously protecting their website from automated threats.
Sentinel achieves this through two delivery models: a fully-featured, self-hostable open-source version and a premium, managed Sentinel Cloud service.
2.0 User Personas
David (The Developer): Runs personal projects and values privacy, performance, and data ownership. Needs a lightweight, powerful tool that is easy to deploy with Docker and doesn't have vendor lock-in.
Maria (The Marketer): Works at a growing e-commerce company. Needs to understand campaign performance, track conversion funnels, and prove ROI without violating GDPR/CCPA. Prefers a managed cloud solution.
Sam (The Founder): Leads a product-driven startup. Needs a single, affordable tool that provides both growth insights (analytics) and defense against threats (ad-fraud bots, scrapers) that drain their budget.
Alex (The Security Engineer): Focuses on site health, performance, and threat monitoring. Needs actionable alerts for anomalies, clear bot identification, and the ability to generate firewall rules to mitigate threats.
3.0 Epics & High-Level User Stories
Epic 1: Onboarding & Setup: A user can get Sentinel running (self-hosted or cloud) and start tracking their first site in under five minutes.
Epic 2: Core Analytics: A user can understand their website's performance at a glance, with clear, actionable metrics and reports.
Epic 3: Security & Intelligence: A user can instantly understand the quality of their traffic and receive actionable suggestions to block automated threats.
Epic 4: Monetization & Growth: A user can seamlessly upgrade to a paid plan (Cloud) or unlock premium features (Self-Hosted) to access more powerful tools as their needs grow.
4.0 Feature Roadmap
4.1 Phase 1: Core MVP (Foundation - F-1xx)
[MVP] F-101: Authentication: Secure user registration, login, and password reset.
[MVP] F-102: Site Management: Full CRUD (Create, Read, Update, Delete) operations for websites within a user's account. Each site gets a unique tracking ID.
[MVP] F-103: Tracking Script: A lightweight (<2KB), asynchronous, cookieless JavaScript snippet. Collects URL, referrer, user agent, screen size, and geo-location (via IP on the backend).
[MVP] F-104: Essential Analytics Dashboard:
Widgets: Real-time visitor count, Unique Visitors, Total Page Views, Bounce Rate, Average Visit Duration.
Reports: Top Pages, Top Referrers, Countries, Devices (Desktop/Mobile), Browsers, Operating Systems.
Filtering: A global date-range picker (e.g., Last 24h, 7d, 30d).
[MVP] F-105: Deployment: A production-ready docker-compose.yml file and clear README.md documentation to enable simple self-hosting.
4.2 Phase 2: Differentiators (The "Why Sentinel" Features - F-2xx)
[Differentiator] F-201: Advanced Bot & Threat Detection:
Implement user agent heuristics, data center IP blocking, and basic headless browser detection.
Assign a trust score to each visitor (e.g., Human / Likely Bot / Known Bot).
[Differentiator] F-202: Firewall UI: A simple interface that allows users to create blocking rules by IP Address, Country, and ASN (Data Center).
[Differentiator] F-203: Web Vitals Monitoring: Track Core Web Vitals (LCP, CLS, FID) for key pages.
[Differentiator] F-204: Funnels & Goals: Allow users to define multi-step conversion funnels and visualize drop-offs.
[Differentiator] F-205: Billing & Subscription (Cloud): Full integration with Stripe for managing SaaS plans.
[Differentiator] F-206: Premium Unlock (Self-Host): A system for activating paid modules (like the Firewall or advanced Funnels) with a license key.
[Differentiator] F-207: The Traffic Quality Score:
This is our signature feature. A single, prominent KPI (0–100) on the dashboard that provides an at-a-glance measure of traffic health.
Inputs: The score will be a composite of several factors, including the percentage of suspected bots, bounce rate anomalies, and unusual geographic or referrer patterns.
4.3 Phase 3: Game-Changers (Market Leadership - F-3xx)
[Game-Changer] F-301: AI Insights Engine: Provide natural-language summaries and suggestions based on traffic patterns (e.g., "Your bounce rate from mobile users in Nigeria spiked by 30% after your last blog post.").
[Game-Changer] F-302: User Flow Visualization: A visual, node-based map showing how users navigate through the site.
[Game-Changer] F-303: Privacy-Focused Heatmaps & Session Replays: Anonymized visual tools to understand user behavior on a deeper level.
[Game-Changer] F-304: A/B Testing Framework: Allow users to define and measure the performance of page variations directly within Sentinel.
5.0 Monetization Model
Self-Hosted: The core analytics features will always be free and open-source. Advanced modules (e.g., "Security Pack," "Funnels Pro") will be available for purchase.
Sentinel Cloud (SaaS): A multi-tenant, managed service with tiered pricing based on monthly events (e.g., Starter: $9/50k events, Growth: $49/500k events).
6.0 Non-Functional Requirements
Performance: The tracking script must not impact Core Web Vitals. The dashboard must load in under 2 seconds for a site with 1 million events in the selected date range.
Privacy: The system must be fully compliant with GDPR/CCPA out of the box. No cookies should be used for tracking.
Usability: The UI must be clean, intuitive, and fully responsive for both desktop and mobile viewing.
7.0 Success Metrics (KPIs)
Community Growth: GitHub stars, active self-hosted instances, community contributions.
Commercial Success: Number of paying Cloud customers, Monthly Recurring Revenue (MRR), conversion rate from free to paid.
Product Stickiness: High user engagement with Phase 2 features (especially the Traffic Quality Score).










Sentinel - Full Development Roadmap
Objective: To provide a complete, sprint-by-sprint development plan to evolve the Sentinel MVP into a market-leading, commercially viable Website Intelligence Platform, encompassing all features outlined in the PRD v2.0.
Phase 1: Core MVP - Status: ✅ COMPLETE
You have successfully built and deployed the foundation of Sentinel. All F-1xx features from the PRD are complete, including authentication, site management, the core analytics dashboard, and a robust Docker-based deployment system.
Phase 2: The Differentiator Sprints (1 Week)
Goal: To build Sentinel's signature security and intelligence features. This is an intense, one-week sprint to create the "wow" factor that will make Sentinel stand out.
Sprint 1: The Security Foundation (Days 1-2)
Feature: F-201 - Advanced Bot & Threat Detection (v1)
Backend Tasks:
Download the GeoLite2-ASN.mmdb file from MaxMind and place it in the backend/ directory.
Update backend/Dockerfile to COPY the new GeoLite2-ASN.mmdb file into the production image.
Update analytics.go:
Add a TrustScore int field to the EventData struct.
In InitAnalyticsEngine, initialize the ASN database reader.
In the TrackHandler, implement the trust scoring logic:
Start with a base score (e.g., 100).
Check User-Agent against a list of known bot signatures and penalize the score.
Perform an ASN lookup on the visitor's IP. If it belongs to a known data center, penalize the score.
Ensure the final TrustScore is saved with every event in the events.log.
Frontend Tasks:
None. This is a backend-only data enrichment step for now.
Feature: F-202 - The Sentinel Firewall (v1)
Backend Tasks:
Create firewall.go in backend/src/ to handle all firewall logic.
Update database.go: Add a CREATE TABLE statement for a firewall_rules table (columns: id, site_id, rule_type, value).
Implement API in firewall.go: Create a FirewallApiHandler supporting GET, POST, and DELETE for firewall rules, ensuring user ownership is checked.
Update main.go: Register the new /api/firewall/ endpoint and protect it with the AuthMiddleware.
Implement Blocking: In the TrackHandler in analytics.go, before processing any event, check the incoming request's IP, Country, and ASN against the firewall rules for that siteId. If there's a match, return a 403 Forbidden response immediately.
Frontend Tasks:
Create FirewallPage.jsx in frontend/src/pages/.
Add a "Firewall" link to the main navigation in the dashboard sidebar.
Update api/index.js: Add new functions for getFirewallRules, addFirewallRule, and deleteFirewallRule.
Build the UI in FirewallPage.jsx with a form to add new rules and a table to list and delete existing rules.
Sprint 2: Performance & Conversion (Days 3-4)
Feature: F-203 - Web Vitals Monitoring
Frontend Tasks:
Install the web-vitals library: npm install web-vitals.
Update the tracking script (tracker-v3.js) to import and use the library to capture LCP, CLS, and FID metrics and send them in the tracking payload.
Update the Dashboard.jsx UI to display these new metrics in new StatCard components.
Backend Tasks:
Update analytics.go: Add optional lcp, cls, fid fields to the Event and EventData structs.
Update calculateStats to compute the average for each Web Vital.
Update the Stats struct to include these new averages in the API response.
Feature: F-204 - Funnels & Goals (v1)
Backend Tasks:
Update database.go: Add a CREATE TABLE statement for a funnels table (columns: id, site_id, name, steps (as JSON)).
Create funnels.go with a FunnelsApiHandler for full CRUD operations.
Update main.go: Register the /api/funnels/ endpoint.
Frontend Tasks:
Create FunnelsPage.jsx in frontend/src/pages/.
Add a "Funnels" link to the dashboard sidebar.
Update api/index.js: Add functions for funnel CRUD.
Build the UI to allow users to define a funnel by specifying an ordered list of URL paths (the drag-and-drop UI will be a future enhancement).
Sprint 3: The Signature KPI (Day 5-6)
Goal: To develop and launch the Traffic Quality Score, Sentinel's most important feature.
Feature: F-208 - The Traffic Quality Score
Backend Tasks:
Update analytics.go: In the calculateStats function, add logic to analyze the TrustScore for all events in the selected date range. Calculate the final TrafficQualityScore as the percentage of events with a TrustScore above a certain threshold (e.g., 50).
Update the Stats struct to include the trafficQualityScore field in the API response.
Frontend Tasks:
Update Dashboard.jsx: Add a new, prominent StatCard or a dedicated gauge component at the top of the dashboard to display the Traffic Quality Score. This should be the most visually striking element on the page.
Sprint 4: Monetization Prep (Day 7)
Goal: To lay the groundwork for future paid plans.
Feature: F-206 & F-207 - Monetization Hooks
Backend Tasks:
Create billing.go as a placeholder file in backend/src/.
Update database.go: Add new columns to the users table to support future billing, such as subscription_status (e.g., 'free', 'pro'), stripe_customer_id, and license_key.
Frontend Tasks:
Create BillingPage.jsx as a simple, static placeholder page that announces that premium features and cloud hosting are "coming soon."
Phase 3: The Game-Changer Sprints (2 Months)
This phase focuses on building the advanced, AI-powered features that will make Sentinel a true market leader.
Sprint 5-6 (2 Weeks): AI Insights & Visualization
F-301: AI Insights Engine: Integrate with an LLM API to generate natural-language summaries of traffic anomalies.
F-302: User Flow Visualization: Use a library like React Flow to create interactive, node-based maps of user journeys.
Sprint 7-8 (2 Weeks): Advanced Behavior Analytics
F-303: Privacy-Focused Heatmaps & Replays: A major feature requiring a new ingestion endpoint and a dedicated UI for replaying anonymized user sessions.
Sprint 9-10 (2 Weeks): Conversion Optimization Tools
F-304: A/B Testing Framework: Build a complete system for users to define and measure A/B tests.
Sprint 11-12 (2 Weeks): Automation & Polish
F-305: Automated Mitigation Suggestions: Enhance the AI engine to suggest specific firewall rules to block identified threats.


