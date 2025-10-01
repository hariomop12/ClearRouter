import React from 'react';
import { Link, useLocation } from 'react-router-dom';

const Sidebar: React.FC = () => {
  const location = useLocation();
  const navItems = [
    { name: 'Dashboard', path: '/dashboard', icon: '🏠' },
    { name: 'API Keys', path: '/dashboard/keys', icon: '🔑' },
    { name: 'Chat', path: '/dashboard/chat', icon: '💬' },
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
    </aside>
  );
};

export default Sidebar;
