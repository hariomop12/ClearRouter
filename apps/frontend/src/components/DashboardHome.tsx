import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import api from '../services/api';

interface UsageStats {
  total_requests: number;
  total_tokens: number;
  total_cost: number;
  top_models: any[];
  top_providers: any[];
  daily_breakdown: any[];
}

interface Credits {
  total_credits: number;
  used_credits: number;
  remaining_credits: number;
}

const DashboardHome: React.FC = () => {
  const { state } = useAuth();
  const { user } = state;
  const [stats, setStats] = useState<UsageStats | null>(null);
  const [credits, setCredits] = useState<Credits | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Fetch analytics data with individual error handling
      let analyticsData = null;
      let creditsData = null;
      
      try {
        const analyticsResponse = await api.get('/analytics/usage?days=30');
        analyticsData = analyticsResponse.data;
        console.log('Analytics data received:', analyticsData);
      } catch (analyticsErr) {
        console.warn('Analytics API failed:', analyticsErr);
        // Set default analytics data
        analyticsData = {
          total_requests: 0,
          total_tokens: 0,
          total_cost: 0,
          top_models: [],
          top_providers: [],
          daily_breakdown: []
        };
      }
      
      try {
        const creditsResponse = await api.get('/credits');
        creditsData = creditsResponse.data;
        console.log('Credits data received:', creditsData);
      } catch (creditsErr) {
        console.warn('Credits API failed:', creditsErr);
        // Set default credits data
        creditsData = {
          total_credits: 0,
          used_credits: 0,
          remaining_credits: 0
        };
      }
      
      setStats(analyticsData);
      setCredits(creditsData);
    } catch (err) {
      console.error('Error fetching dashboard data:', err);
      setError('Failed to load dashboard data');
      // Set fallback data to prevent crashes
      setStats({
        total_requests: 0,
        total_tokens: 0,
        total_cost: 0,
        top_models: [],
        top_providers: [],
        daily_breakdown: []
      });
      setCredits({
        total_credits: 0,
        used_credits: 0,
        remaining_credits: 0
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDashboardData();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500"></div>
      </div>
    );
  }

  return (
    <>
      <div className="mb-8">
        <h2 className="text-3xl font-bold text-white mb-2">
          Welcome back, {user?.name}!
        </h2>
        <p className="text-gray-300">
          Here's your ClearRouter usage overview for the last 30 days
        </p>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-500/20 border border-red-500/50 rounded-lg">
          <p className="text-red-200">{error}</p>
        </div>
      )}

      {/* Analytics Cards */}
      <div className="grid md:grid-cols-3 gap-6 mb-8">
        {/* Organization Credits */}
        <div className="bg-gray-900/50 backdrop-blur-sm rounded-2xl border border-gray-800 p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-gray-400 text-sm font-medium">User Credits</h3>
            <div className="p-2 bg-green-500/20 rounded-lg">
              <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
              </svg>
            </div>
          </div>
          <div className="text-3xl font-bold text-white mb-2">
            {credits && credits.remaining_credits ? credits.remaining_credits.toLocaleString() : '0'}
          </div>
          <div className="text-sm text-gray-400">
            Available Balance
          </div>
        </div>

        {/* Total Requests */}
        <div className="bg-gray-900/50 backdrop-blur-sm rounded-2xl border border-gray-800 p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-gray-400 text-sm font-medium">Total Requests</h3>
            <div className="p-2 bg-blue-500/20 rounded-lg">
              <svg className="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
          </div>
          <div className="text-3xl font-bold text-white mb-2">
            {stats && stats.total_requests ? stats.total_requests.toLocaleString() : '0'}
          </div>
          <div className="text-sm text-gray-400">
            Last 30 days
          </div>
        </div>

        {/* Tokens Used */}
        <div className="bg-gray-900/50 backdrop-blur-sm rounded-2xl border border-gray-800 p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-gray-400 text-sm font-medium">Tokens Used</h3>
            <div className="p-2 bg-purple-500/20 rounded-lg">
              <svg className="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
          </div>
          <div className="text-3xl font-bold text-white mb-2">
            {stats && stats.total_tokens ? stats.total_tokens.toLocaleString() : '0'}
          </div>
          <div className="text-sm text-gray-400">
            Last 30 days
          </div>
        </div>
      </div>

      {/* Top Models and Providers */}
      {stats && (
        <div className="grid md:grid-cols-2 gap-8 mb-8">
          {/* Top Models */}
          <div className="bg-gray-900/50 backdrop-blur-sm rounded-2xl border border-gray-800 p-6">
            <h3 className="text-xl font-bold text-white mb-4">Top Models</h3>
            {stats.top_models && stats.top_models.length > 0 ? (
              <div className="space-y-3">
                {stats.top_models.slice(0, 5).map((model: any, index: number) => (
                  <div key={index} className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                    <div>
                      <p className="text-white font-medium">{model.model || 'Unknown Model'}</p>
                      <p className="text-sm text-gray-400">{model.requests || 0} requests</p>
                    </div>
                    <div className="text-right">
                      <p className="text-white font-medium">{model.tokens ? model.tokens.toLocaleString() : '0'}</p>
                      <p className="text-xs text-gray-400">tokens</p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-400">No model usage data available</p>
            )}
          </div>

          {/* Top Providers */}
          <div className="bg-gray-900/50 backdrop-blur-sm rounded-2xl border border-gray-800 p-6">
            <h3 className="text-xl font-bold text-white mb-4">Top Providers</h3>
            {stats.top_providers && stats.top_providers.length > 0 ? (
              <div className="space-y-3">
                {stats.top_providers.slice(0, 5).map((provider: any, index: number) => (
                  <div key={index} className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                    <div>
                      <p className="text-white font-medium capitalize">{provider.provider || 'Unknown Provider'}</p>
                      <p className="text-sm text-gray-400">{provider.requests || 0} requests</p>
                    </div>
                    <div className="text-right">
                      <p className="text-white font-medium">${provider.cost ? provider.cost.toFixed(4) : '0.0000'}</p>
                      <p className="text-xs text-gray-400">cost</p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-400">No provider usage data available</p>
            )}
          </div>
        </div>
      )}

      {/* Quick Actions */}
      <div className="bg-gray-900/50 backdrop-blur-sm rounded-2xl border border-gray-800 p-6">
        <h3 className="text-xl font-bold text-white mb-6">Quick Actions</h3>
        <div className="grid md:grid-cols-4 gap-4">
          <Link
            to="/dashboard/keys"
            className="flex items-center space-x-3 p-4 bg-gray-800/50 rounded-lg hover:bg-gray-700/50 transition-colors group"
          >
            <div className="w-10 h-10 bg-gradient-to-r from-blue-600 to-blue-700 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v-2H7v-2H4a1 1 0 01-1-1v-1m0 0a6 6 0 0113.255-3.257M0 9h2.25" />
              </svg>
            </div>
            <div>
              <p className="text-white font-medium">Create API Key</p>
              <p className="text-xs text-gray-400">Generate new key</p>
            </div>
          </Link>

          <Link
            to="/dashboard/chat"
            className="flex items-center space-x-3 p-4 bg-gray-800/50 rounded-lg hover:bg-gray-700/50 transition-colors group"
          >
            <div className="w-10 h-10 bg-gradient-to-r from-green-600 to-green-700 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <div>
              <p className="text-white font-medium">Start Chat</p>
              <p className="text-xs text-gray-400">AI conversation</p>
            </div>
          </Link>

          <Link
            to="/dashboard/models"
            className="flex items-center space-x-3 p-4 bg-gray-800/50 rounded-lg hover:bg-gray-700/50 transition-colors group"
          >
            <div className="w-10 h-10 bg-gradient-to-r from-purple-600 to-purple-700 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
            </div>
            <div>
              <p className="text-white font-medium">Browse Models</p>
              <p className="text-xs text-gray-400">Available AI models</p>
            </div>
          </Link>

          <Link
            to="/dashboard/credits"
            className="flex items-center space-x-3 p-4 bg-gray-800/50 rounded-lg hover:bg-gray-700/50 transition-colors group"
          >
            <div className="w-10 h-10 bg-gradient-to-r from-yellow-600 to-yellow-700 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
              </svg>
            </div>
            <div>
              <p className="text-white font-medium">Add Credits</p>
              <p className="text-xs text-gray-400">Purchase more</p>
            </div>
          </Link>
        </div>
      </div>
    </>
  );
};

export default DashboardHome;