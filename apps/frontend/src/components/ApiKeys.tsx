import React, { useState, useEffect } from 'react';
import api from '../services/api';

interface ApiKey {
  id: string;
  user_id: string;
  api_key: string;
  active: boolean;
  created_at: string;
}

const ApiKeys: React.FC = () => {
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchApiKeys = async () => {
    try {
      setLoading(true);
      const response = await api.get('/keys');
      setApiKeys(response.data);
    } catch (err) {
      setError('Failed to fetch API keys');
      console.error('Error fetching API keys:', err);
    } finally {
      setLoading(false);
    }
  };

  const createApiKey = async () => {
    try {
      setCreating(true);
      const response = await api.post('/keys/create');
      setApiKeys([response.data, ...apiKeys]);
    } catch (err) {
      setError('Failed to create API key');
      console.error('Error creating API key:', err);
    } finally {
      setCreating(false);
    }
  };

  const deleteApiKey = async (keyId: string) => {
    try {
      await api.delete(`/keys/${keyId}`);
      setApiKeys(apiKeys.filter(key => key.id !== keyId));
    } catch (err) {
      setError('Failed to delete API key');
      console.error('Error deleting API key:', err);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    // TODO: Add toast notification
    alert('API key copied to clipboard!');
  };

  useEffect(() => {
    fetchApiKeys();
  }, []);

  return (
    <div className="max-w-6xl mx-auto">
      <div className="mb-8">
        <h2 className="text-3xl font-bold text-white mb-4">API Keys</h2>
        <p className="text-gray-300">Manage your API keys for accessing ClearRouter services</p>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-500/20 border border-red-500/50 rounded-lg">
          <p className="text-red-200">{error}</p>
          <button 
            onClick={() => setError(null)}
            className="mt-2 text-red-300 hover:text-red-100 text-sm underline"
          >
            Dismiss
          </button>
        </div>
      )}

      {/* Create API Key Button */}
      <div className="mb-8">
        <button
          onClick={createApiKey}
          disabled={creating}
          className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-medium rounded-lg transition-all transform hover:scale-105 disabled:opacity-50 disabled:transform-none"
        >
          {creating ? (
            <>
              <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Creating...
            </>
          ) : (
            <>
              <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Create New API Key
            </>
          )}
        </button>
      </div>

      {/* API Keys List */}
      <div className="bg-white/5 backdrop-blur-sm rounded-2xl border border-gray-800 overflow-hidden">
        <div className="px-6 py-4 bg-gray-800/50 border-b border-gray-700">
          <h3 className="text-xl font-semibold text-white">Your API Keys</h3>
        </div>

        {loading ? (
          <div className="p-8 text-center">
            <div className="inline-flex items-center text-gray-300">
              <svg className="animate-spin -ml-1 mr-3 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Loading API keys...
            </div>
          </div>
        ) : apiKeys.length === 0 ? (
          <div className="p-8 text-center">
            <div className="text-gray-400 mb-4">
              <svg className="w-12 h-12 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v-2H7v-2H4a1 1 0 01-1-1v-1m0 0a6 6 0 0113.255-3.257M0 9h2.25" />
              </svg>
              <p className="text-lg">No API keys found</p>
              <p className="text-sm">Create your first API key to get started</p>
            </div>
          </div>
        ) : (
          <div className="divide-y divide-gray-700">
            {apiKeys.map((apiKey) => (
              <div key={apiKey.id} className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <div className="flex items-center space-x-3 mb-2">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        apiKey.active 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {apiKey.active ? 'Active' : 'Inactive'}
                      </span>
                      <span className="text-sm text-gray-400">
                        Created: {new Date(apiKey.created_at).toLocaleDateString()}
                      </span>
                    </div>
                    
                    <div className="flex items-center space-x-3">
                      <code className="bg-gray-800 text-green-300 px-3 py-2 rounded-lg font-mono text-sm flex-1 truncate">
                        {apiKey.api_key}
                      </code>
                      <button
                        onClick={() => copyToClipboard(apiKey.api_key)}
                        className="inline-flex items-center px-3 py-2 border border-gray-600 rounded-lg text-sm text-gray-300 hover:text-white hover:border-gray-500 transition-colors"
                      >
                        <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                        </svg>
                        Copy
                      </button>
                    </div>
                  </div>
                  
                  <div className="ml-4">
                    <button
                      onClick={() => deleteApiKey(apiKey.id)}
                      className="inline-flex items-center px-3 py-2 border border-red-600 rounded-lg text-sm text-red-300 hover:text-white hover:bg-red-600 transition-all"
                    >
                      <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Usage Instructions */}
      <div className="mt-8 bg-blue-500/10 border border-blue-500/30 rounded-lg p-6">
        <h4 className="text-lg font-semibold text-blue-300 mb-2">How to use your API key</h4>
        <p className="text-blue-200 mb-4">Include your API key in the Authorization header when making requests to the ClearRouter API:</p>
        <pre className="bg-gray-900 text-green-300 p-4 rounded-lg overflow-x-auto text-sm">
{`curl -X POST 'http://localhost:8080/v1/chat/completions' \\
  -H 'Authorization: Bearer YOUR_API_KEY' \\
  -H 'Content-Type: application/json' \\
  -d '{
    "model": "gemini-2.5-flash-lite",
    "messages": [
      {"role": "user", "content": "Hello, world!"}
    ]
  }'`}
        </pre>
      </div>
    </div>
  );
};

export default ApiKeys;