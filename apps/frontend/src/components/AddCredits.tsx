import React, { useState, useEffect } from 'react';
import api from '../services/api';

interface CreditPackage {
  amount: number;
  credits: number;
  bonus?: number;
  popular?: boolean;
}

interface Credits {
  total_credits: number;
  used_credits: number;
  remaining_credits: number;
}

const AddCredits: React.FC = () => {
  const [credits, setCredits] = useState<Credits | null>(null);
  const [loading, setLoading] = useState(true);
  const [purchasing, setPurchasing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const creditPackages: CreditPackage[] = [
    { amount: 100, credits: 100 },
    { amount: 500, credits: 550, bonus: 50 },
    { amount: 1000, credits: 1200, bonus: 200, popular: true },
    { amount: 2000, credits: 2500, bonus: 500 },
    { amount: 5000, credits: 6500, bonus: 1500 },
  ];

  const fetchCredits = async () => {
    try {
      setLoading(true);
      const response = await api.get('/credits');
      setCredits(response.data);
    } catch (err) {
      console.error('Error fetching credits:', err);
      setError('Failed to load credit information');
    } finally {
      setLoading(false);
    }
  };

  const purchaseCredits = async (amount: number) => {
    try {
      setPurchasing(true);
      setError(null);

      // Create Razorpay order
      const orderResponse = await api.post('/credits/order', {
        amount: amount * 100, // Convert to paise
        currency: 'INR'
      });

      const { id: order_id, amount: order_amount, currency } = orderResponse.data;

      // Initialize Razorpay
      const options = {
        key: 'rzp_test_RMdWZKqulEGpGi', // Your Razorpay key ID
        amount: order_amount,
        currency: currency,
        name: 'ClearRouter',
        description: `Purchase ${creditPackages.find(pkg => pkg.amount === amount)?.credits} credits`,
        order_id: order_id,
        handler: async function (response: any) {
          try {
            // Verify payment on backend
            await api.post('/credits/add', {
              razorpay_payment_id: response.razorpay_payment_id,
              razorpay_order_id: response.razorpay_order_id,
              razorpay_signature: response.razorpay_signature
            });

            // Refresh credits
            await fetchCredits();
            alert('Credits added successfully!');
          } catch (err) {
            console.error('Payment verification failed:', err);
            setError('Payment verification failed. Please contact support.');
          }
        },
        prefill: {
          name: 'User', // You can get this from auth context
          email: 'user@example.com' // You can get this from auth context
        },
        theme: {
          color: '#8B5CF6'
        },
        modal: {
          ondismiss: function() {
            setPurchasing(false);
          }
        }
      };

      const razorpay = new (window as any).Razorpay(options);
      razorpay.open();
    } catch (err: any) {
      console.error('Error purchasing credits:', err);
      setError(err.response?.data?.error || 'Failed to initiate payment');
      setPurchasing(false);
    }
  };

  useEffect(() => {
    fetchCredits();
    
    // Load Razorpay script
    const script = document.createElement('script');
    script.src = 'https://checkout.razorpay.com/v1/checkout.js';
    script.async = true;
    document.body.appendChild(script);

    return () => {
      document.body.removeChild(script);
    };
  }, []);

  return (
    <div className="max-w-6xl mx-auto">
      <div className="mb-8">
        <h2 className="text-3xl font-bold text-white mb-4">Credits</h2>
        <p className="text-gray-300">Purchase credits to use ClearRouter AI services</p>
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

      {/* Current Credits */}
      <div className="mb-8 bg-white/5 backdrop-blur-sm rounded-2xl border border-gray-800 p-6">
        <h3 className="text-xl font-semibold text-white mb-4">Current Balance</h3>
        {loading ? (
          <div className="animate-pulse">
            <div className="h-4 bg-gray-700 rounded w-1/3 mb-2"></div>
            <div className="h-4 bg-gray-700 rounded w-1/4"></div>
          </div>
        ) : credits ? (
          <div className="grid md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="text-3xl font-bold text-green-400">{credits.total_credits}</div>
              <div className="text-sm text-gray-400">Total Credits</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-red-400">{credits.used_credits}</div>
              <div className="text-sm text-gray-400">Used Credits</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-blue-400">{credits.remaining_credits}</div>
              <div className="text-sm text-gray-400">Available Credits</div>
            </div>
          </div>
        ) : (
          <div className="text-gray-400">Failed to load credit information</div>
        )}
      </div>

      {/* Credit Packages */}
      <div className="mb-8">
        <h3 className="text-xl font-semibold text-white mb-6">Purchase Credits</h3>
        <div className="grid md:grid-cols-3 lg:grid-cols-5 gap-4">
          {creditPackages.map((pkg, index) => (
            <div
              key={index}
              className={`relative bg-white/5 backdrop-blur-sm rounded-xl border p-6 text-center transition-all hover:bg-white/10 ${
                pkg.popular ? 'border-purple-500 ring-2 ring-purple-500/20' : 'border-gray-800'
              }`}
            >
              {pkg.popular && (
                <div className="absolute -top-3 left-1/2 transform -translate-x-1/2">
                  <span className="bg-gradient-to-r from-purple-600 to-pink-600 text-white px-3 py-1 rounded-full text-xs font-medium">
                    Popular
                  </span>
                </div>
              )}
              
              <div className="mb-4">
                <div className="text-2xl font-bold text-white">₹{pkg.amount}</div>
                <div className="text-sm text-gray-400">Indian Rupees</div>
              </div>

              <div className="mb-4">
                <div className="text-xl font-semibold text-green-400">{pkg.credits} Credits</div>
                {pkg.bonus && (
                  <div className="text-sm text-purple-400">+{pkg.bonus} bonus credits</div>
                )}
              </div>

              <button
                onClick={() => purchaseCredits(pkg.amount)}
                disabled={purchasing}
                className={`w-full px-4 py-2 rounded-lg font-medium transition-all ${
                  pkg.popular
                    ? 'bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white'
                    : 'bg-gray-700 hover:bg-gray-600 text-white'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                {purchasing ? 'Processing...' : 'Purchase'}
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Usage Information */}
      <div className="bg-blue-500/10 border border-blue-500/30 rounded-lg p-6">
        <h4 className="text-lg font-semibold text-blue-300 mb-4">How Credits Work</h4>
        <div className="space-y-3 text-blue-200">
          <div className="flex items-start space-x-3">
            <div className="w-2 h-2 bg-blue-400 rounded-full mt-2 flex-shrink-0"></div>
            <p>Credits are consumed based on the AI model you use and the number of tokens processed</p>
          </div>
          <div className="flex items-start space-x-3">
            <div className="w-2 h-2 bg-blue-400 rounded-full mt-2 flex-shrink-0"></div>
            <p>Different models have different pricing - premium models cost more credits per token</p>
          </div>
          <div className="flex items-start space-x-3">
            <div className="w-2 h-2 bg-blue-400 rounded-full mt-2 flex-shrink-0"></div>
            <p>You can monitor your credit usage in real-time through your dashboard</p>
          </div>
          <div className="flex items-start space-x-3">
            <div className="w-2 h-2 bg-blue-400 rounded-full mt-2 flex-shrink-0"></div>
            <p>Credits never expire - use them at your own pace</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AddCredits;