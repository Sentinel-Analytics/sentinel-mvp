import React, { useEffect, useRef, useState } from 'react';
import { api } from '../api';

const SessionReplay = () => {
    const playerRef = useRef(null);
    const [sites, setSites] = useState([]);
    const [selectedSite, setSelectedSite] = useState(null);
    const [sessions, setSessions] = useState([]);
    const [selectedSession, setSelectedSession] = useState(null);
    const [loadingEvents, setLoadingEvents] = useState(false);
    const [playerLoaded, setPlayerLoaded] = useState(false);

    useEffect(() => {
        const fetchSites = async () => {
            try {
                const sitesData = await api.getSites();
                setSites(sitesData || []);
                if (sitesData && sitesData.length > 0) {
                    setSelectedSite(sitesData[0]);
                }
            } catch (error) {
                console.error("Failed to fetch sites:", error);
            }
        };
        fetchSites();
    }, []);

    useEffect(() => {
        const fetchSessions = async () => {
            if (selectedSite) {
                try {
                    const sessionIDs = await api.listSessions(selectedSite.id);
                    setSessions(sessionIDs || []);
                    setSelectedSession(sessionIDs && sessionIDs.length > 0 ? sessionIDs[0] : null);
                } catch (error) {
                    console.error("Failed to fetch sessions:", error);
                    setSessions([]);
                    setSelectedSession(null);
                }
            }
        };
        fetchSessions();
    }, [selectedSite]);

    useEffect(() => {
        const loadRrwebPlayer = () => {
            // Inject CSS
            const cssLink = document.createElement('link');
            cssLink.rel = 'stylesheet';
            cssLink.href = 'https://cdn.jsdelivr.net/npm/rrweb-player@latest/dist/style.css';
            document.head.appendChild(cssLink);

            // Inject Script
            const script = document.createElement('script');
            script.src = 'https://cdn.jsdelivr.net/npm/rrweb-player@latest/dist/index.js';
            script.async = true;
            document.body.appendChild(script);

            script.onload = () => {
                setPlayerLoaded(true); // Signal that the player script is ready
            };

            return () => {
                document.head.removeChild(cssLink);
                document.body.removeChild(script);
            };
        };

        loadRrwebPlayer();
    }, []);

    useEffect(() => {
        const fetchAndPlaySession = async () => {
            if (selectedSite && selectedSession && playerLoaded && window.rrwebPlayer) {
                setLoadingEvents(true);
                try {
                    const fetchedEvents = await api.getSessionEvents(selectedSite.id, selectedSession);
                    const events = Array.isArray(fetchedEvents) ? fetchedEvents : [];
                    
                    if (playerRef.current) {
                        playerRef.current.innerHTML = ''; // Clear previous player
                        
                        if (events.length > 0) {
                            new window.rrwebPlayer({
                                target: playerRef.current,
                                props: {
                                    events,
                                    width: playerRef.current.clientWidth,
                                    height: 600, // Adjust as needed
                                    autoPlay: true,
                                },
                            });
                        } else {
                            // Optionally, display a message if there are no events for the session
                            playerRef.current.innerHTML = '<div class="flex items-center justify-center h-64 text-slate-400">No events recorded for this session.</div>';
                        }
                    }
                } catch (error) {
                    console.error("Failed to fetch session events:", error);
                    if (playerRef.current) {
                        playerRef.current.innerHTML = '<div class="flex items-center justify-center h-64 text-red-400">Failed to load session events.</div>';
                    }
                } finally {
                    setLoadingEvents(false);
                }
            }
        };
        fetchAndPlaySession();
    }, [selectedSite, selectedSession, playerLoaded]);

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4 text-slate-200">Session Replay</h1>

            <div className="mb-4 flex space-x-4">
                <div className="flex-1">
                    <label htmlFor="site-select" className="block text-sm font-medium text-slate-400">Select Site</label>
                    <select
                        id="site-select"
                        className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md bg-slate-700 text-slate-200"
                        value={selectedSite ? selectedSite.id : ''}
                        onChange={(e) => setSelectedSite(sites.find(site => site.id === e.target.value))}
                    >
                        {sites.map(site => (
                            <option key={site.id} value={site.id}>{site.name}</option>
                        ))}
                    </select>
                </div>

                <div className="flex-1">
                    <label htmlFor="session-select" className="block text-sm font-medium text-slate-400">Select Session</label>
                    <select
                        id="session-select"
                        className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md bg-slate-700 text-slate-200"
                        value={selectedSession || ''}
                        onChange={(e) => setSelectedSession(e.target.value)}
                    >
                        {sessions.map(sessionID => (
                            <option key={sessionID} value={sessionID}>{sessionID}</option>
                        ))}
                    </select>
                </div>
            </div>

            <div className="bg-slate-800 rounded-lg shadow-lg p-4">
                {loadingEvents ? (
                    <div className="flex items-center justify-center h-64 text-slate-400">
                        Loading session events...
                    </div>
                ) : selectedSession ? (
                    <div ref={playerRef} className="rrweb-player-container"></div>
                ) : (
                    <div className="flex items-center justify-center h-64 text-slate-400">
                        No session selected or no sessions available for this site.
                    </div>
                )}
            </div>
        </div>
    );
};

export default SessionReplay;
