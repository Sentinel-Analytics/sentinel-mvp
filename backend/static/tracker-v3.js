(function() {
    // The script tag will have the data-site-id attribute, which we will add later.
    const scriptTag = document.querySelector('script[data-site-id]');
    if (!scriptTag) {
        console.error('Sentinel: data-site-id attribute not found on script tag.');
        return;
    }
    const siteId = scriptTag.getAttribute('data-site-id');
    const aPIEndpoint = 'https://api-sentinel.getmusterup.com/track' // Relative URL is fine

    function track() {
        const data = {
            siteId: siteId, // Include the Site ID
            url: window.location.href,
            referrer: document.referrer || '',
            screenWidth: window.screen.width,
        };

        fetch(aPIEndpoint, {
            method: 'POST',
            body: JSON.stringify(data),
            headers: { 
                'Content-Type': 'application/json' 
            },
            keepalive: true,
        }).catch(error => console.error('Sentinel tracking error:', error));
    }
    
    // Track the initial page view
    track();
})();
