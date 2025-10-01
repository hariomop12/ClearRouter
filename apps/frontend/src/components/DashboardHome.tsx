import React from 'react';
import { useAuth } from '../contexts/AuthContext';

const DashboardHome: React.FC = () => {
  const { state } = useAuth();
  const { user } = state;

  return (
    <>
      <div className="text-center mb-12">
        <h2 className="text-4xl font-bold text-white mb-4">
          🎉 You are logged in!
        </h2>
        <p className="text-xl text-gray-300">
          Welcome to your ClearRouter dashboard
        </p>
      </div>

      {/* User Info Card */}
      <div className="max-w-2xl mx-auto">
        <div className="bg-white/5 backdrop-blur-sm rounded-2xl border border-gray-800 p-8">
          <div className="text-center mb-8">
            <div className="w-20 h-20 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <span className="text-2xl font-bold text-white">
                {user?.name?.charAt(0).toUpperCase()}
              </span>
            </div>
            <h3 className="text-2xl font-bold text-white mb-2">User Profile</h3>
            <p className="text-gray-400">Your account information</p>
          </div>

          <div className="space-y-6">
            <div className="flex items-center justify-between p-4 bg-gray-800/50 rounded-lg">
              <div className="flex items-center space-x-3">
                <svg className="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                <span className="text-gray-300 font-medium">Name</span>
              </div>
              <span className="text-white">{user?.name}</span>
            </div>

            <div className="flex items-center justify-between p-4 bg-gray-800/50 rounded-lg">
              <div className="flex items-center space-x-3">
                <svg className="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                <span className="text-gray-300 font-medium">Email</span>
              </div>
              <span className="text-white">{user?.email}</span>
            </div>

            <div className="flex items-center justify-between p-4 bg-gray-800/50 rounded-lg">
              <div className="flex items-center space-x-3">
                <svg className="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v-2H7v-2H4a1 1 0 01-1-1v-1m0 0a6 6 0 0113.255-3.257M0 9h2.25" />
                </svg>
                <span className="text-gray-300 font-medium">User ID</span>
              </div>
              <span className="text-white font-mono text-sm">{user?.id}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Features Overview */}
      <div className="max-w-4xl mx-auto mt-12">
        <h3 className="text-2xl font-bold text-white text-center mb-8">Features</h3>
        <div className="grid md:grid-cols-3 gap-6">
          <div className="bg-white/5 backdrop-blur-sm rounded-xl border border-gray-800 p-6 text-center">
            <div className="w-12 h-12 bg-gradient-to-r from-blue-600 to-blue-700 rounded-lg flex items-center justify-center mx-auto mb-4">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <h4 className="text-lg font-semibold text-white mb-2">AI Chat</h4>
            <p className="text-gray-400 text-sm">Start conversations with multiple AI providers</p>
          </div>

          <div className="bg-white/5 backdrop-blur-sm rounded-xl border border-gray-800 p-6 text-center">
            <div className="w-12 h-12 bg-gradient-to-r from-green-600 to-green-700 rounded-lg flex items-center justify-center mx-auto mb-4">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v-2H7v-2H4a1 1 0 01-1-1v-1m0 0a6 6 0 0113.255-3.257M0 9h2.25" />
              </svg>
            </div>
            <h4 className="text-lg font-semibold text-white mb-2">API Keys</h4>
            <p className="text-gray-400 text-sm">Manage your API keys and integrations</p>
          </div>

          <div className="bg-white/5 backdrop-blur-sm rounded-xl border border-gray-800 p-6 text-center">
            <div className="w-12 h-12 bg-gradient-to-r from-yellow-600 to-yellow-700 rounded-lg flex items-center justify-center mx-auto mb-4">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
              </svg>
            </div>
            <h4 className="text-lg font-semibold text-white mb-2">Credits</h4>
            <p className="text-gray-400 text-sm">Monitor your usage and purchase credits</p>
          </div>
        </div>
      </div>
    </>
  );
};

export default DashboardHome;