import React from 'react';

const Info: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-900">
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent mb-4">
            ClearRouter Information
          </h1>
          <p className="text-gray-300 text-lg">
            Everything you need to know about using our AI chat service
          </p>
        </div>

        {/* What are Tokens Section */}
        <div className="bg-gray-800/50 rounded-xl p-8 mb-8 border border-gray-700">
          <h2 className="text-2xl font-semibold text-white mb-6 flex items-center">
            <span className="bg-blue-500 rounded-full w-8 h-8 flex items-center justify-center text-sm font-bold mr-3">1</span>
            What are Tokens?
          </h2>
          <div className="text-gray-300 space-y-4">
            <p>
              Think of <strong className="text-blue-400">tokens</strong> as small pieces of text that AI models use to understand and generate responses.
            </p>
            <div className="bg-gray-700/50 rounded-lg p-4">
              <p><strong className="text-yellow-400">Example:</strong></p>
              <p>The sentence <span className="text-green-400">"Hello, how are you?"</span> might be broken into:</p>
              <ul className="list-disc list-inside mt-2 ml-4">
                <li><span className="text-green-400">"Hello"</span> = 1 token</li>
                <li><span className="text-green-400">","</span> = 1 token</li>
                <li><span className="text-green-400">"how"</span> = 1 token</li>
                <li><span className="text-green-400">"are"</span> = 1 token</li>
                <li><span className="text-green-400">"you"</span> = 1 token</li>
                <li><span className="text-green-400">"?"</span> = 1 token</li>
              </ul>
              <p className="mt-2"><strong>Total:</strong> 6 tokens</p>
            </div>
            <p>
              <strong className="text-purple-400">Simple rule:</strong> 1 token ≈ 4 characters or 0.75 words on average.
            </p>
          </div>
        </div>

        {/* How We Calculate Price Section */}
        <div className="bg-gray-800/50 rounded-xl p-8 mb-8 border border-gray-700">
          <h2 className="text-2xl font-semibold text-white mb-6 flex items-center">
            <span className="bg-green-500 rounded-full w-8 h-8 flex items-center justify-center text-sm font-bold mr-3">2</span>
            How We Calculate Prices
          </h2>
          <div className="text-gray-300 space-y-6">
            <p>Our pricing is simple and transparent:</p>
            
            <div className="bg-gray-700/50 rounded-lg p-6">
              <h3 className="text-xl font-semibold text-yellow-400 mb-4">💰 Price Formula</h3>
              <div className="bg-gray-900 rounded-lg p-4 font-mono text-green-400 text-lg">
                Total Cost = (Input Tokens × Input Price) + (Output Tokens × Output Price)
              </div>
            </div>

            <div className="grid md:grid-cols-2 gap-6">
              <div className="bg-blue-900/30 rounded-lg p-4 border border-blue-500/30">
                <h4 className="text-blue-400 font-semibold mb-2">📥 Input Tokens</h4>
                <p>Your question or message to the AI</p>
                <p className="text-sm text-gray-400 mt-1">What you send to the AI model</p>
              </div>
              <div className="bg-purple-900/30 rounded-lg p-4 border border-purple-500/30">
                <h4 className="text-purple-400 font-semibold mb-2">📤 Output Tokens</h4>
                <p>The AI's response back to you</p>
                <p className="text-sm text-gray-400 mt-1">What the AI generates for you</p>
              </div>
            </div>

            <div className="bg-gray-700/50 rounded-lg p-6">
              <h3 className="text-xl font-semibold text-yellow-400 mb-4">🧮 Example Calculation</h3>
              <div className="space-y-2">
                <p><strong>Your question:</strong> "What is artificial intelligence?" (5 tokens)</p>
                <p><strong>AI response:</strong> "AI is computer technology that can think and learn..." (50 tokens)</p>
                <p><strong>Model:</strong> GPT-4o Mini</p>
                <div className="bg-gray-900 rounded p-3 mt-3">
                  <p>• Input: 5 tokens × $0.15/1M tokens = $0.00000075</p>
                  <p>• Output: 50 tokens × $0.60/1M tokens = $0.00003000</p>
                  <p className="text-green-400 font-semibold">• Total Cost: $0.00003075</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Available AI Models Section */}
        <div className="bg-gray-800/50 rounded-xl p-8 mb-8 border border-gray-700">
          <h2 className="text-2xl font-semibold text-white mb-6 flex items-center">
            <span className="bg-purple-500 rounded-full w-8 h-8 flex items-center justify-center text-sm font-bold mr-3">3</span>
            Available AI Models
          </h2>
          <div className="text-gray-300 space-y-4">
            <p>We support <strong className="text-green-400">66+ AI models</strong> from 5 major providers:</p>
            
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div className="bg-gray-700/50 rounded-lg p-4">
                <h4 className="text-blue-400 font-semibold mb-2">🤖 OpenAI</h4>
                <ul className="text-sm space-y-1">
                  <li>• GPT-4o, GPT-4o Mini</li>
                  <li>• GPT-4, GPT-3.5 Turbo</li>
                  <li>• o1, o3 (Reasoning)</li>
                </ul>
              </div>
              <div className="bg-gray-700/50 rounded-lg p-4">
                <h4 className="text-purple-400 font-semibold mb-2">🧠 Anthropic</h4>
                <ul className="text-sm space-y-1">
                  <li>• Claude 3.5 Sonnet</li>
                  <li>• Claude 3 Opus</li>
                  <li>• Claude 3 Haiku</li>
                </ul>
              </div>
              <div className="bg-gray-700/50 rounded-lg p-4">
                <h4 className="text-green-400 font-semibold mb-2">🌟 Google</h4>
                <ul className="text-sm space-y-1">
                  <li>• Gemini 2.5 Pro</li>
                  <li>• Gemini 1.5 Flash</li>
                  <li>• Gemma 3 Series</li>
                </ul>
              </div>
              <div className="bg-gray-700/50 rounded-lg p-4">
                <h4 className="text-orange-400 font-semibold mb-2">🚀 DeepSeek</h4>
                <ul className="text-sm space-y-1">
                  <li>• DeepSeek V3</li>
                  <li>• DeepSeek Coder</li>
                  <li>• DeepSeek R1</li>
                </ul>
              </div>
              <div className="bg-gray-700/50 rounded-lg p-4">
                <h4 className="text-red-400 font-semibold mb-2">⚡ Mistral</h4>
                <ul className="text-sm space-y-1">
                  <li>• Mistral Large</li>
                  <li>• Mixtral 8x7B</li>
                  <li>• Pixtral (Vision)</li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        {/* How to Use Section */}
        <div className="bg-gray-800/50 rounded-xl p-8 mb-8 border border-gray-700">
          <h2 className="text-2xl font-semibold text-white mb-6 flex items-center">
            <span className="bg-orange-500 rounded-full w-8 h-8 flex items-center justify-center text-sm font-bold mr-3">4</span>
            How to Use ClearRouter
          </h2>
          <div className="text-gray-300 space-y-6">
            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <h3 className="text-xl font-semibold text-blue-400 mb-4">🌐 Web Interface</h3>
                <ol className="list-decimal list-inside space-y-2">
                  <li>Create an account and login</li>
                  <li>Add credits to your account</li>
                  <li>Choose an AI model</li>
                  <li>Start chatting!</li>
                </ol>
              </div>
              <div>
                <h3 className="text-xl font-semibold text-purple-400 mb-4">🔌 API Integration</h3>
                <ol className="list-decimal list-inside space-y-2">
                  <li>Generate an API key</li>
                  <li>Use our OpenAI-compatible API</li>
                  <li>Integrate with your apps</li>
                  <li>Monitor usage and costs</li>
                </ol>
              </div>
            </div>
          </div>
        </div>

        {/* Features Section */}
        <div className="bg-gray-800/50 rounded-xl p-8 mb-8 border border-gray-700">
          <h2 className="text-2xl font-semibold text-white mb-6 flex items-center">
            <span className="bg-pink-500 rounded-full w-8 h-8 flex items-center justify-center text-sm font-bold mr-3">5</span>
            Key Features
          </h2>
          <div className="grid md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div className="flex items-start space-x-3">
                <span className="text-green-400 text-xl">✅</span>
                <div>
                  <h4 className="text-white font-semibold">Pay-per-use Pricing</h4>
                  <p className="text-gray-400 text-sm">Only pay for what you actually use</p>
                </div>
              </div>
              <div className="flex items-start space-x-3">
                <span className="text-green-400 text-xl">🔒</span>
                <div>
                  <h4 className="text-white font-semibold">Secure & Private</h4>
                  <p className="text-gray-400 text-sm">Your data is protected and encrypted</p>
                </div>
              </div>
              <div className="flex items-start space-x-3">
                <span className="text-green-400 text-xl">📊</span>
                <div>
                  <h4 className="text-white font-semibold">Usage Analytics</h4>
                  <p className="text-gray-400 text-sm">Track your usage and spending</p>
                </div>
              </div>
            </div>
            <div className="space-y-4">
              <div className="flex items-start space-x-3">
                <span className="text-green-400 text-xl">🚀</span>
                <div>
                  <h4 className="text-white font-semibold">Fast Responses</h4>
                  <p className="text-gray-400 text-sm">Get AI responses in seconds</p>
                </div>
              </div>
              <div className="flex items-start space-x-3">
                <span className="text-green-400 text-xl">🔑</span>
                <div>
                  <h4 className="text-white font-semibold">API Access</h4>
                  <p className="text-gray-400 text-sm">Integrate with your own applications</p>
                </div>
              </div>
              <div className="flex items-start space-x-3">
                <span className="text-green-400 text-xl">💬</span>
                <div>
                  <h4 className="text-white font-semibold">Chat History</h4>
                  <p className="text-gray-400 text-sm">Access your previous conversations</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Contact Section */}
        <div className="bg-gradient-to-r from-purple-900/50 to-pink-900/50 rounded-xl p-8 border border-purple-500/30">
          <h2 className="text-2xl font-semibold text-white mb-4 text-center">
            Need Help?
          </h2>
          <p className="text-gray-300 text-center mb-6">
            Have questions or need support? We're here to help!
          </p>
          <div className="flex justify-center space-x-4">
            <button className="px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors">
              Contact Support
            </button>
            <button className="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors">
              View Documentation
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Info;