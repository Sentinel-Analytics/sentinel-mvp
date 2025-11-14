# Phase 2: Differentiator Features

This document outlines the key features that will be implemented during Phase 2. The goal of this phase is to build Sentinel's signature intelligence features that provide deep, actionable insights beyond traditional analytics.

### Feature Checklist

- [ ] **F-201: Advanced Bot & Threat Detection**: Assign a `TrustScore` to each visitor based on heuristics like user agent, data center IP, and headless browser detection. This is the foundation for our security features.

- [ ] **F-202: The Sentinel Firewall**: Build the UI and API to allow users to create blocking rules by IP Address, Country, and ASN (Data Center).

- [ ] **F-203: Web Vitals Monitoring**: Enhance the tracking script to capture and report on Core Web Vitals (LCP, CLS, FID) for key pages.

- [ ] **F-204: Funnels & Goals**: Create a UI for users to define multi-step conversion funnels and visualize drop-offs.

- [ ] **F-208: The Traffic Quality Score**: Develop our signature KPI. This will be a single, prominent score (0-100) on the dashboard that provides an at-a-glance measure of traffic health based on bot scores, bounce rates, and other signals.

---

### User-Suggested Intelligence Features

- [ ] **F-209: Privacy-Focused Session Replays (v1)**: Record and replay user sessions to understand their journey, identify pain points, and see where they encounter issues. All sensitive input will be anonymized by default.

- [ ] **F-210: Engagement Heatmaps**: Generate visual heatmaps to show where users click and scroll the most on key pages, providing insight into what content is most engaging.

- [ ] **F-211: User Issue Detection**: Automatically identify user frustration signals like "rage clicks" (repeated, rapid clicking) and "dead clicks" (clicks on non-interactive elements) to pinpoint UI/UX problems.
