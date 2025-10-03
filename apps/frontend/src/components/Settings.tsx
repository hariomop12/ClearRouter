import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { userAPI } from '../services/api';

const Settings: React.FC = () => {
  const { state, logout, updateUser } = useAuth();
  const { user } = state;
  const [name, setName] = useState(user?.name || '');
  const [isUpdatingName, setIsUpdatingName] = useState(false);
  const [isDeletingAccount, setIsDeletingAccount] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleUpdateUsername = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      setError('Name is required');
      return;
    }

    setIsUpdatingName(true);
    setError('');
    setMessage('');

    try {
      await userAPI.updateUsername(name.trim());
      setMessage('Username updated successfully!');
      // Update user in auth context and localStorage
      if (state.user) {
        const updatedUser = { ...state.user, name: name.trim() };
        updateUser(updatedUser);
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to update username';
      setError(errorMessage);
    } finally {
      setIsUpdatingName(false);
    }
  };

  const handleDeleteAccount = async () => {
    setIsDeletingAccount(true);
    setError('');
    setMessage('');

    try {
      await userAPI.deleteAccount();
      setMessage('Account deleted successfully. Redirecting...');
      // Logout and redirect
      setTimeout(() => {
        logout();
      }, 2000);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to delete account';
      setError(errorMessage);
    } finally {
      setIsDeletingAccount(false);
      setShowDeleteConfirm(false);
    }
  };

  return (
    <div className="p-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white mb-2">Account Settings</h1>
        <p className="text-gray-400">Manage your account preferences and settings</p>
      </div>

      {/* Messages */}
      {message && (
        <div className="mb-6 p-4 bg-green-900/50 border border-green-700 rounded-lg">
          <p className="text-green-300">{message}</p>
        </div>
      )}

      {error && (
        <div className="mb-6 p-4 bg-red-900/50 border border-red-700 rounded-lg">
          <p className="text-red-300">{error}</p>
        </div>
      )}

      {/* Profile Information */}
      <div className="bg-gray-800/50 rounded-xl p-6 mb-8">
        <h2 className="text-xl font-semibold text-white mb-4">Profile Information</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Email Address
            </label>
            <input
              type="email"
              value={user?.email || ''}
              disabled
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-gray-400 cursor-not-allowed"
            />
            <p className="text-xs text-gray-500 mt-1">Email cannot be changed</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Account ID
            </label>
            <input
              type="text"
              value={user?.id || ''}
              disabled
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-gray-400 cursor-not-allowed font-mono text-sm"
            />
          </div>
        </div>
      </div>

      {/* Update Username */}
      <div className="bg-gray-800/50 rounded-xl p-6 mb-8">
        <h2 className="text-xl font-semibold text-white mb-4">Update Username</h2>
        
        <form onSubmit={handleUpdateUsername} className="max-w-md">
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Display Name
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="Enter your display name"
              minLength={2}
              maxLength={255}
              required
            />
          </div>
          
          <button
            type="submit"
            disabled={isUpdatingName || !name.trim() || name === user?.name}
            className="bg-gradient-to-r from-purple-600 to-pink-600 text-white px-6 py-2 rounded-lg font-medium hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed transition-opacity"
          >
            {isUpdatingName ? 'Updating...' : 'Update Username'}
          </button>
        </form>
      </div>

      {/* Danger Zone */}
      <div className="bg-red-900/20 border border-red-800 rounded-xl p-6">
        <h2 className="text-xl font-semibold text-red-300 mb-4">Danger Zone</h2>
        <p className="text-gray-300 mb-4">
          Once you delete your account, there is no going back. Please be certain.
        </p>
        
        {!showDeleteConfirm ? (
          <button
            onClick={() => setShowDeleteConfirm(true)}
            className="bg-red-600 text-white px-6 py-2 rounded-lg font-medium hover:bg-red-700 transition-colors"
          >
            Delete Account
          </button>
        ) : (
          <div className="space-y-4">
            <div className="p-4 bg-red-900/30 border border-red-700 rounded-lg">
              <p className="text-red-200 font-medium mb-2">⚠️ This action cannot be undone!</p>
              <p className="text-red-300 text-sm">
                This will permanently delete your account, API keys, chat history, credits, and all associated data.
              </p>
            </div>
            
            <div className="flex space-x-3">
              <button
                onClick={handleDeleteAccount}
                disabled={isDeletingAccount}
                className="bg-red-600 text-white px-6 py-2 rounded-lg font-medium hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {isDeletingAccount ? 'Deleting...' : 'Yes, Delete My Account'}
              </button>
              
              <button
                onClick={() => setShowDeleteConfirm(false)}
                disabled={isDeletingAccount}
                className="bg-gray-600 text-white px-6 py-2 rounded-lg font-medium hover:bg-gray-700 disabled:opacity-50 transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Settings;