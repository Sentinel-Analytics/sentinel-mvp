import React, { useState, useEffect } from 'react';
import { api } from '../api'; // Assuming api functions exist

const FirewallPage = () => {
  const [rules, setRules] = useState([]);
  const [ruleType, setRuleType] = useState('ip');
  const [value, setValue] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const siteId = localStorage.getItem('siteId'); // Or get it from context/URL

  useEffect(() => {
    if (siteId) {
      fetchRules();
    }
  }, [siteId]);

  const fetchRules = async () => {
    try {
      setLoading(true);
      const fetchedRules = await api.getFirewallRules(siteId);
      setRules(fetchedRules);
      setError(null);
    } catch (err) {
      setError('Failed to fetch firewall rules.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddRule = async (e) => {
    e.preventDefault();
    if (!value.trim()) {
      setError('Rule value cannot be empty.');
      return;
    }
    try {
      await api.addFirewallRule(siteId, { rule_type: ruleType, value });
      setValue('');
      fetchRules(); // Refresh the list
    } catch (err) {
      setError('Failed to add rule.');
      console.error(err);
    }
  };

  const handleDeleteRule = async (ruleId) => {
    try {
      await api.deleteFirewallRule(siteId, ruleId);
      fetchRules(); // Refresh the list
    } catch (err) {
      setError('Failed to delete rule.');
      console.error(err);
    }
  };

  if (!siteId) {
    return <div className="text-center p-8">Please select a site to manage its firewall.</div>;
  }

  if (loading) {
    return <div className="text-center p-8">Loading firewall rules...</div>;
  }

  return (
    <div className="p-4 sm:p-6 lg:p-8">
      <h1 className="text-2xl font-bold mb-6">Sentinel Firewall</h1>

      <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6 mb-8">
        <h2 className="text-xl font-semibold mb-4">Add New Rule</h2>
        <form onSubmit={handleAddRule} className="flex flex-col sm:flex-row gap-4">
          <select
            value={ruleType}
            onChange={(e) => setRuleType(e.target.value)}
            className="p-2 border rounded-md bg-white dark:bg-gray-700"
          >
            <option value="ip">IP Address</option>
            <option value="country">Country (2-letter code)</option>
            <option value="asn">ASN</option>
          </select>
          <input
            type="text"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder="Enter value to block"
            className="flex-grow p-2 border rounded-md bg-white dark:bg-gray-700"
          />
          <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600">
            Add Rule
          </button>
        </form>
        {error && <p className="text-red-500 mt-4">{error}</p>}
      </div>

      <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-4">Existing Rules</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full">
            <thead>
              <tr className="border-b">
                <th className="text-left p-2">Type</th>
                <th className="text-left p-2">Value</th>
                <th className="text-right p-2">Actions</th>
              </tr>
            </thead>
            <tbody>
              {rules.length > 0 ? (
                rules.map((rule) => (
                  <tr key={rule.id} className="border-b">
                    <td className="p-2">{rule.rule_type}</td>
                    <td className="p-2">{rule.value}</td>
                    <td className="text-right p-2">
                      <button
                        onClick={() => api.deleteFirewallRule(siteId, rule.id)}
                        className="bg-red-500 text-white px-3 py-1 rounded-md hover:bg-red-600"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan="3" className="text-center p-4">
                    No firewall rules defined.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default FirewallPage;
