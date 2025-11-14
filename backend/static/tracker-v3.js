(function() {
    const scriptTag = document.querySelector('script[data-site-id]');
    if (!scriptTag) {
        console.error('Sentinel: data-site-id attribute not found on script tag.');
        return;
    }
    const siteId = scriptTag.getAttribute('data-site-id');
    const apiEndpoint = 'https://api-sentinel.getmusterup.com/track';

    let lastUrl = location.href;

    function track() {
        const data = {
            siteId: siteId,
            url: window.location.href,
            referrer: document.referrer || '',
            screenWidth: window.screen.width,
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

    // --- SPA Tracking ---
    // 1. Track initial page view
    track();

    // 2. Wrap history.pushState to track URL changes
    const originalPushState = history.pushState;
    history.pushState = function(...args) {
        originalPushState.apply(this, args);
        if (location.href !== lastUrl) {
            lastUrl = location.href;
            track();
        }
    };

    // 3. Listen for popstate event (browser back/forward buttons)
    window.addEventListener('popstate', () => {
        if (location.href !== lastUrl) {
            lastUrl = location.href;
            track();
        }
    });
})();
