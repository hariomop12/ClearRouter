import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../services/api';

interface Provider {
  provider_id: string;
  model_name: string;
  input_price: number;
  output_price: number;
  context_size: number;
  max_output: number;
  streaming: boolean;
  vision?: boolean;
  tools?: boolean;
}

interface Model {
  id: string;
  name: string;
  family: string;
  providers: Provider[];
  status: string;
  json_output: boolean;
}

interface ModelsResponse {
  data: Model[];
}

const Models: React.FC = () => {
  const [models, setModels] = useState<Model[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedFamily, setSelectedFamily] = useState<string>('all');
  const [searchTerm, setSearchTerm] = useState<string>('');

  useEffect(() => {
    const fetchModels = async () => {
      try {
        setLoading(true);
        const response = await api.get<ModelsResponse>('/models');
        setModels(response.data.data);
      } catch (err: any) {
        setError('Failed to fetch models');
        console.error('Error fetching models:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchModels();
  }, []);

  const families = ['all', ...Array.from(new Set(models.map(model => model.family)))];

  const filteredModels = models.filter(model => {
    const matchesFamily = selectedFamily === 'all' || model.family === selectedFamily;
    const matchesSearch = model.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
                         model.id.toLowerCase().includes(searchTerm.toLowerCase());
    return matchesFamily && matchesSearch;
  });

  const formatPrice = (price: number): string => {
    if (price === 0) return 'Free';
    if (price < 0.000001) return `$${(price * 1000000).toFixed(2)}/1M tokens`;
    if (price < 0.001) return `$${(price * 1000).toFixed(3)}/1K tokens`;
    return `$${price.toFixed(4)}/token`;
  };

  const formatContextSize = (size: number): string => {
    if (size >= 1000000) return `${(size / 1000000).toFixed(1)}M`;
    if (size >= 1000) return `${(size / 1000).toFixed(0)}K`;
    return size.toString();
  };

  const getProviderLogo = (family: string): string => {
    switch (family) {
      case 'openai':
        return '🤖';
      case 'google':
        return '🔥';
      default:
        return '⚡';
    }
  };

  const getProviderColor = (family: string): string => {
    switch (family) {
      case 'openai':
        return 'from-green-500 to-emerald-600';
      case 'google':
        return 'from-blue-500 to-indigo-600';
      default:
        return 'from-purple-500 to-pink-600';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-violet-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="flex items-center justify-center min-h-[60vh]">
            <div className="text-center">
              <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-purple-400 mx-auto mb-4"></div>
              <p className="text-gray-300">Loading models...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-violet-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="flex items-center justify-center min-h-[60vh]">
            <div className="text-center">
              <div className="text-red-400 text-6xl mb-4">⚠️</div>
              <p className="text-red-300 text-lg">{error}</p>
              <button 
                onClick={() => window.location.reload()} 
                className="mt-4 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-violet-900">
      {/* Header */}
      <header className="bg-black/20 backdrop-blur-sm border-b border-gray-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link to="/" className="flex items-center">
              <h1 className="text-2xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
                ClearRouter
              </h1>
            </Link>
            <nav className="flex items-center space-x-6">
              <Link to="/" className="text-gray-300 hover:text-white transition-colors">
                Home
              </Link>
              <Link to="/models" className="text-purple-400 font-medium">
                Models
              </Link>
              <Link to="/login" className="text-gray-300 hover:text-white transition-colors">
                Login
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Page Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl sm:text-5xl font-bold text-white mb-4">
            AI Models
          </h1>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto">
            Access the most powerful AI models from leading providers through our unified API
          </p>
        </div>

        {/* Filters */}
        <div className="mb-8 flex flex-col sm:flex-row gap-4 items-center justify-between">
          {/* Search */}
          <div className="relative flex-1 max-w-md">
            <svg className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              type="text"
              placeholder="Search models..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
            />
          </div>

          {/* Family Filter */}
          <div className="flex flex-wrap gap-2">
            {families.map((family) => (
              <button
                key={family}
                onClick={() => setSelectedFamily(family)}
                className={`px-4 py-2 rounded-lg font-medium transition-all capitalize ${
                  selectedFamily === family
                    ? 'bg-purple-600 text-white'
                    : 'bg-gray-800/50 text-gray-300 hover:bg-gray-700/50 hover:text-white'
                }`}
              >
                {family === 'all' ? 'All Providers' : family}
              </button>
            ))}
          </div>
        </div>

        {/* Models Count */}
        <div className="mb-6">
          <p className="text-gray-400">
            Showing {filteredModels.length} of {models.length} models
          </p>
        </div>

        {/* Models Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredModels.map((model) => {
            const provider = model.providers[0]; // Use first provider for display
            return (
              <div
                key={model.id}
                className="bg-white/5 backdrop-blur-sm border border-gray-800 rounded-2xl p-6 hover:bg-white/10 transition-all duration-300 hover:scale-105 hover:border-gray-700"
              >
                {/* Header */}
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center space-x-3">
                    <div className={`w-12 h-12 rounded-xl bg-gradient-to-r ${getProviderColor(model.family)} flex items-center justify-center text-2xl`}>
                      {getProviderLogo(model.family)}
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold text-white">{model.name}</h3>
                      <p className="text-sm text-gray-400 capitalize">{model.family}</p>
                    </div>
                  </div>
                  
                  {/* Status Badge */}
                  <div className="flex items-center space-x-2">
                    {provider.streaming && (
                      <span className="px-2 py-1 text-xs bg-green-900/50 text-green-300 rounded-full">
                        Streaming
                      </span>
                    )}
                    {provider.vision && (
                      <span className="px-2 py-1 text-xs bg-blue-900/50 text-blue-300 rounded-full">
                        Vision
                      </span>
                    )}
                    {provider.tools && (
                      <span className="px-2 py-1 text-xs bg-purple-900/50 text-purple-300 rounded-full">
                        Tools
                      </span>
                    )}
                  </div>
                </div>

                {/* Pricing */}
                <div className="grid grid-cols-2 gap-4 mb-4">
                  <div className="bg-gray-800/30 rounded-lg p-3">
                    <p className="text-xs text-gray-400 mb-1">Input</p>
                    <p className="text-sm font-medium text-white">{formatPrice(provider.input_price)}</p>
                  </div>
                  <div className="bg-gray-800/30 rounded-lg p-3">
                    <p className="text-xs text-gray-400 mb-1">Output</p>
                    <p className="text-sm font-medium text-white">{formatPrice(provider.output_price)}</p>
                  </div>
                </div>

                {/* Specs */}
                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-400">Context Size</span>
                    <span className="text-sm text-white font-medium">{formatContextSize(provider.context_size)}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-400">Max Output</span>
                    <span className="text-sm text-white font-medium">{formatContextSize(provider.max_output)}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-400">JSON Output</span>
                    <span className={`text-sm font-medium ${model.json_output ? 'text-green-400' : 'text-red-400'}`}>
                      {model.json_output ? 'Yes' : 'No'}
                    </span>
                  </div>
                </div>

                {/* Try Button */}
                <button className="w-full mt-4 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white py-2 px-4 rounded-lg font-medium transition-all duration-200 hover:scale-105">
                  Try {model.name}
                </button>
              </div>
            );
          })}
        </div>

        {/* No Results */}
        {filteredModels.length === 0 && (
          <div className="text-center py-12">
            <div className="text-gray-400 text-6xl mb-4">🔍</div>
            <p className="text-gray-300 text-lg">No models found matching your criteria</p>
            <button 
              onClick={() => {
                setSearchTerm('');
                setSelectedFamily('all');
              }}
              className="mt-4 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
            >
              Clear Filters
            </button>
          </div>
        )}
      </main>
    </div>
  );
};

export default Models;