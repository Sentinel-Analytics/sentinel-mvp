(function() {
    const scriptTag = document.querySelector('script[data-site-id]');
    if (!scriptTag) {
        console.error('Sentinel: data-site-id attribute not found on script tag.');
        return;
    }
    const siteId = scriptTag.getAttribute('data-site-id');
    const apiEndpoint = 'https://api-sentinel.getmusterup.com/track';
    const sessionEndpoint = 'https://api-sentinel.getmusterup.com/session';

    // Inject rrweb
    const rrwebScript = document.createElement('script');
    rrwebScript.src = 'https://cdn.jsdelivr.net/npm/rrweb@latest/dist/rrweb.min.js';
    document.head.appendChild(rrwebScript);

    // Inject web-vitals
    const webVitalsScript = document.createElement('script');
    webVitalsScript.type = 'module';
    webVitalsScript.textContent = `
        import { onCLS, onFID, onLCP } from 'https://unpkg.com/web-vitals@4?module';
        
        const sendVitals = ({ name, value }) => {
            // This function will be defined later in the main script
            window.trackWebVitals({ [name]: value });
        };

        onCLS(sendVitals, { reportAllChanges: true });
        onFID(sendVitals, { reportAllChanges: true });
        onLCP(sendVitals, { reportAllChanges: true });
    `;
    document.head.appendChild(webVitalsScript);


    rrwebScript.onload = function() {
        let events = [];
        let lastUrl = location.href;
        let vitals = {};

        function track(vitalsToSend) {
            const data = {
                siteId: siteId,
                url: window.location.href,
                referrer: document.referrer || '',
                screenWidth: window.screen.width,
                ...vitalsToSend
            };

            fetch(apiEndpoint, {
                method: 'POST',
                body: JSON.stringify(data),
                headers: {
                    'Content-Type': 'application/json'
                },
                keepalive: true,
            }).catch(error => console.error('Sentinel tracking error:', error));
        }

        // Expose a global function for web-vitals to call
        window.trackWebVitals = (vital) => {
            vitals = { ...vitals, ...vital };
        };


        // --- SPA Tracking ---
        // Track initial page view
        setTimeout(() => track(vitals), 100); // Send initial vitals after a short delay

        const originalPushState = history.pushState;
        history.push_state = function(...args) {
            originalPushState.apply(this, args);
            if (location.href !== lastUrl) {
                lastUrl = location.href;
                vitals = {}; // Reset vitals for new page view
                setTimeout(() => track(vitals), 100); // Send vitals for new page
            }
        };

        window.addEventListener('popstate', () => {
            if (location.href !== lastUrl) {
                lastUrl = location.href;
                vitals = {}; // Reset vitals for new page view
                setTimeout(() => track(vitals), 100); // Send vitals for new page
            }
        });

        // --- RRWeb Session Recording ---
        rrweb.record({
            emit(event) {
                events.push(event);
            },
        });

        let sessionId = null;

        // Save events every 10 seconds
        setInterval(() => {
            if (events.length > 0) {
                const body = JSON.stringify({ siteId: siteId, events: events, sessionId: sessionId });
                events = []; // Clear buffer
                fetch(sessionEndpoint, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: body,
                    keepalive: true
                })
                .then(response => response.json())
                .then(data => {
                    if (data.sessionId) {
                        sessionId = data.sessionId;
                    }
                })
                .catch(err => console.error('Sentinel session recording error:', err));
            }
        }, 10 * 1000);
    }
})();