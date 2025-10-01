import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const UserInfoSection: React.FC = () => {
  const { state } = useAuth();
  const { user } = state;

  return (
    <div className="px-4 py-4 border-t border-gray-800">
      <div className="flex items-center space-x-3 p-3 bg-gray-800/50 rounded-lg">
        <div className="w-10 h-10 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full flex items-center justify-center flex-shrink-0">
          <span className="text-sm font-bold text-white">
            {user?.name?.charAt(0).toUpperCase()}
          </span>
        </div>
        <div className="min-w-0 flex-1">
          <p className="text-sm font-medium text-white truncate">
            {user?.name}
          </p>
          <p className="text-xs text-gray-400 truncate">
            {user?.email}
          </p>
        </div>
      </div>
    </div>
  );
};

const Sidebar: React.FC = () => {
  const location = useLocation();
  const navItems = [
    { name: 'Dashboard', path: '/dashboard', icon: '🏠' },
    { name: 'API Keys', path: '/dashboard/keys', icon: '🔑' },
    { name: 'Chat', path: '/dashboard/chat', icon: '💬' },
    { name: 'Models', path: '/dashboard/models', icon: '🤖' },
    { name: 'Add Credits', path: '/dashboard/credits', icon: '💰' },
  ];

  return (
    <aside className="fixed top-0 left-0 h-full w-64 bg-gray-900 border-r border-gray-800 flex flex-col z-20">
      <div className="px-6 py-6 border-b border-gray-800">
        <h2 className="text-2xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">ClearRouter</h2>
      </div>
      <nav className="flex-1 px-4 py-6 space-y-2">
        {navItems.map(item => (
          <Link
            key={item.path}
            to={item.path}
            className={`flex items-center px-4 py-2 rounded-lg text-lg font-medium transition-colors ${location.pathname === item.path ? 'bg-gradient-to-r from-purple-600 to-pink-600 text-white' : 'text-gray-300 hover:bg-gray-800 hover:text-white'}`}
          >
            <span className="mr-3 text-xl">{item.icon}</span>
            {item.name}
          </Link>
        ))}
      </nav>
      
      {/* User Info Section */}
      <UserInfoSection />
    </aside>
  );
};

export default Sidebar;
