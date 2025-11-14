# Proposal: Building Sentinel's Differentiator Features

This document outlines the strategic plan, technology choices, and implementation process for the features in Phase 2. My approach is guided by three core principles:

1.  **Leverage the Existing Stack:** We will build upon the solid foundation of Go (backend) and React (frontend) to ensure consistency and maintainability.
2.  **Prioritize Performance:** Any new tracking capabilities must have a negligible impact on the end-user's website performance. All new features will be designed with efficiency in mind.
3.  **Champion Privacy:** For features like session replays, privacy is paramount. We will build in privacy-by-default mechanisms, such as anonymizing user input.

---

## The Implementation Plan

Here is a feature-by-feature breakdown of the technologies and processes I will use:

### Part 1: Advanced Bot & Threat Detection

This is a foundational, backend-only feature that enriches our existing data.

*   **Process:**
    1.  **ASN Database:** I will download the `GeoLite2-ASN.mmdb` database from MaxMind, which maps IP addresses to Autonomous System Numbers (i.e., data centers like AWS, Google Cloud).
    2.  **Backend Logic:** In `analytics.go`, I will extend the `TrackHandler`. For each incoming event, I will:
        *   Perform an ASN lookup on the visitor's IP using the `oschwald/geoip2-golang` library we already use for country lookups.
        *   Check the User-Agent against a curated list of known bot signatures.
        *   Calculate a `TrustScore` (e.g., starting at 100, deducting points for bot-like signals).
    3.  **Database:** I will add a `TrustScore` (integer) column to the `events` table in our ClickHouse schema.

*   **Technologies:**
    *   **Go:** For all backend logic.
    *   **MaxMind GeoLite2-ASN:** For data center detection.
    *   **ClickHouse:** To store the new `TrustScore` data point.

---

### Part 2: The Sentinel Firewall

This feature will give users direct control over their traffic.

*   **Process:**
    1.  **Backend API:** I will create a new `firewall.go` file to handle all firewall logic. This will include a new RESTful API endpoint (`/api/firewall`) supporting `GET`, `POST`, and `DELETE` to manage rules.
    2.  **Database (SQLite):** A new `firewall_rules` table will be created in our existing SQLite database to store rules for each site.
    3.  **Blocking Logic:** The `TrackHandler` will be updated to query the firewall rules *before* processing an event. If a request matches a rule, the server will respond with a `403 Forbidden` status, stopping the request dead in its tracks.
    4.  **Frontend UI:** I will create a new `FirewallPage.jsx` in the frontend. This page will feature a simple form to add new rules (by IP, Country, or ASN) and a table to display and delete existing rules.

*   **Technologies:**
    *   **Go:** For the API and blocking logic.
    *   **SQLite:** To store firewall rules.
    *   **React:** To build the new user interface for managing the firewall.

---

### Part 3: Session Replays, Heatmaps & User Issue Detection

This is the most significant and impactful group of features. I'll tackle them in order, starting with the foundation: event capturing.

*   **Process & Technologies:**

    1.  **Event Recording (The Foundation):**
        *   **Technology:** I will integrate the **`rrweb`** library into our `tracker-v3.js` script. `rrweb` is the industry-standard, open-source library for recording and replaying web sessions. It's lightweight, efficient, and captures everything needed for replays and heatmaps.
        *   **Process:** The tracker will be updated to not just send a single pageview, but to capture DOM changes, mouse movements, and clicks. It will batch these events and send them to a new backend endpoint. It will also automatically detect and flag "rage clicks" and "dead clicks."

    2.  **Backend Ingestion & Storage:**
        *   **Process:** I will create a new backend endpoint (e.g., `/api/session/record`) to receive the batched events from `rrweb`. To handle the potentially large volume of data, I will create a new `session_events` table in ClickHouse, optimized for storing this time-series data.
        *   **Privacy:** The `rrweb` recorder will be configured to automatically anonymize all user input fields, ensuring sensitive data never leaves the user's browser.

    3.  **Frontend Playback & Visualization:**
        *   **Session Replays (F-209):** I will create a new `SessionReplaysPage.jsx`. This page will list recent user sessions. Clicking on a session will open a player view that uses the `rrweb-player` component to reconstruct and play back the user's exact journey.
        *   **Heatmaps & Issue Detection (F-210, F-211):** Once we are capturing click data, I will create a new `InsightsPage.jsx`.
            *   Initially, this page will display simple reports, such as a list of the top pages where "rage clicks" are occurring.
            *   For visual heatmaps, I will use a library like **`heatmap.js`** to overlay the aggregated click data onto a representation of the user's webpage, providing a clear visual guide to user engagement.

---
