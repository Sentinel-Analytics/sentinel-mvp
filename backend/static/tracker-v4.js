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

    // Dynamically import and initialize web-vitals
    const initWebVitals = () => {
        import('https://unpkg.com/web-vitals@4?module')
            .then(({ onCLS, onFID, onLCP }) => {
                const sendVitals = ({ name, value }) => {
                    window.trackWebVitals({ [name]: value });
                };
                onCLS(sendVitals);
                onFID(sendVitals);
                onLCP(sendVitals);
            })
            .catch(err => {
                console.error('Sentinel: Failed to load web-vitals library.', err);
            });
    };


    rrwebScript.onload = function() {
        // Start Web Vitals monitoring
        initWebVitals();

        let events = [];
        let lastUrl = location.href;

        function track(payload = {}) {
            const data = {
                siteId: siteId,
                url: window.location.href,
                referrer: document.referrer || '',
                screenWidth: window.screen.width,
                ...payload
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
            track(vital);
        };


        // --- SPA Tracking ---
        // Track initial page view
        track();

        const originalPushState = history.pushState;
        history.pushState = function(...args) {
            originalPushState.apply(this, args);
            if (location.href !== lastUrl) {
                lastUrl = location.href;
                track(); // Track navigation
            }
        };

        window.addEventListener('popstate', () => {
            if (location.href !== lastUrl) {
                lastUrl = location.href;
                track(); // Track navigation
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