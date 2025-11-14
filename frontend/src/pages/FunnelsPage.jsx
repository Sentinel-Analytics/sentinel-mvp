import React, { useState, useEffect } from 'react';
import { api } from '../api';

const FunnelsPage = () => {
    const [sites, setSites] = useState([]);
    const [selectedSite, setSelectedSite] = useState(null);
    const [funnels, setFunnels] = useState([]);
    const [newFunnelName, setNewFunnelName] = useState('');
    const [newFunnelSteps, setNewFunnelSteps] = useState(['']);

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
        const fetchFunnels = async () => {
            if (selectedSite) {
                try {
                    const funnelsData = await api.listFunnels(selectedSite.id);
                    setFunnels(funnelsData || []);
                } catch (error) {
                    console.error("Failed to fetch funnels:", error);
                    setFunnels([]);
                }
            }
        };
        fetchFunnels();
    }, [selectedSite]);

    const handleAddStep = () => {
        setNewFunnelSteps([...newFunnelSteps, '']);
    };

    const handleStepChange = (index, value) => {
        const updatedSteps = [...newFunnelSteps];
        updatedSteps[index] = value;
        setNewFunnelSteps(updatedSteps);
    };

    const handleRemoveStep = (index) => {
        const updatedSteps = [...newFunnelSteps];
        updatedSteps.splice(index, 1);
        setNewFunnelSteps(updatedSteps);
    };

    const handleCreateFunnel = async (e) => {
        e.preventDefault();
        if (!selectedSite || !newFunnelName.trim() || newFunnelSteps.some(step => !step.trim())) {
            alert("Please fill in all fields.");
            return;
        }
        try {
            await api.createFunnel({
                siteId: selectedSite.id,
                name: newFunnelName,
                steps: newFunnelSteps,
            });
            setNewFunnelName('');
            setNewFunnelSteps(['']);
            // Refetch funnels
            const funnelsData = await api.listFunnels(selectedSite.id);
            setFunnels(funnelsData || []);
        } catch (error) {
            console.error("Failed to create funnel:", error);
        }
    };

    const handleDeleteFunnel = async (funnelId) => {
        if (window.confirm("Are you sure you want to delete this funnel?")) {
            try {
                await api.deleteFunnel(funnelId);
                // Refetch funnels
                const funnelsData = await api.listFunnels(selectedSite.id);
                setFunnels(funnelsData || []);
            } catch (error) {
                console.error("Failed to delete funnel:", error);
            }
        }
    };

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4 text-slate-200">Funnels & Goals</h1>

            <div className="mb-4">
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

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div>
                    <h2 className="text-xl font-semibold mb-4 text-slate-200">Create New Funnel</h2>
                    <form onSubmit={handleCreateFunnel} className="space-y-4 bg-slate-800 p-6 rounded-lg">
                        <div>
                            <label htmlFor="funnel-name" className="block text-sm font-medium text-slate-400">Funnel Name</label>
                            <input
                                type="text"
                                id="funnel-name"
                                value={newFunnelName}
                                onChange={(e) => setNewFunnelName(e.target.value)}
                                className="mt-1 block w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-md text-slate-200 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-slate-400">Funnel Steps (URL Paths)</label>
                            {newFunnelSteps.map((step, index) => (
                                <div key={index} className="flex items-center space-x-2 mt-2">
                                    <input
                                        type="text"
                                        value={step}
                                        onChange={(e) => handleStepChange(index, e.target.value)}
                                        placeholder="/example/path"
                                        className="block w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-md text-slate-200 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                    />
                                    <button type="button" onClick={() => handleRemoveStep(index)} className="text-red-500 hover:text-red-700">Remove</button>
                                </div>
                            ))}
                            <button type="button" onClick={handleAddStep} className="mt-2 text-indigo-400 hover:text-indigo-600">Add Step</button>
                        </div>
                        <button type="submit" className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded-md">Create Funnel</button>
                    </form>
                </div>
                <div>
                    <h2 className="text-xl font-semibold mb-4 text-slate-200">Existing Funnels</h2>
                    <div className="space-y-4">
                        {funnels.map(funnel => (
                            <div key={funnel.id} className="bg-slate-800 p-4 rounded-lg">
                                <div className="flex justify-between items-center">
                                    <h3 className="font-semibold text-slate-200">{funnel.name}</h3>
                                    <button onClick={() => handleDeleteFunnel(funnel.id)} className="text-red-500 hover:text-red-700">Delete</button>
                                </div>
                                <ol className="list-decimal list-inside mt-2 text-slate-400">
                                    {funnel.steps.map((step, index) => (
                                        <li key={index}>{step}</li>
                                    ))}
                                </ol>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default FunnelsPage;
