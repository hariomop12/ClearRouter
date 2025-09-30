# ClearRouter API Collection for Postman

## 🔑 Authentication
Replace `YOUR_API_KEY` with your actual API key: `HNNACLIs9au8CoUPeAIPyN--w-UZd8ASAKV61xRYy9I=`

---

## 📋 **Models API**

### 1. Get All Available Models
```bash
curl --location 'http://localhost:8080/models' \
--header 'Content-Type: application/json'
```

**Expected Response:**
```json
{
  "data": [
    {
      "id": "gpt-4o-mini",
      "name": "GPT-4o Mini",
      "family": "openai",
      "providers": [...],
      "json_output": true
    },
    {
      "id": "gemini-2.5-flash",
      "name": "Gemini 2.5 Flash", 
      "family": "google",
      "providers": [...],
      "json_output": true
    }
  ]
}
```

---

## 💬 **Chat Completions API**

### 2. OpenAI GPT-4o Mini
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gpt-4o-mini",
    "messages": [
        {
            "role": "user",
            "content": "Hello! How are you?"
        }
    ]
}'
```

### 3. OpenAI GPT-4o
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gpt-4o",
    "messages": [
        {
            "role": "user",
            "content": "Explain quantum computing in simple terms"
        }
    ]
}'
```

### 4. OpenAI GPT-4
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gpt-4",
    "messages": [
        {
            "role": "system",
            "content": "You are a helpful assistant."
        },
        {
            "role": "user",
            "content": "Write a Python function to calculate fibonacci numbers"
        }
    ]
}'
```

### 5. OpenAI GPT-3.5 Turbo
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gpt-3.5-turbo",
    "messages": [
        {
            "role": "user",
            "content": "What is the capital of France?"
        }
    ]
}'
```

### 6. Google Gemini 2.5 Flash
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gemini-2.5-flash",
    "messages": [
        {
            "role": "user",
            "content": "Explain the difference between AI and Machine Learning"
        }
    ]
}'
```

### 7. Google Gemini 2.5 Pro
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gemini-2.5-pro",
    "messages": [
        {
            "role": "user",
            "content": "Write a detailed analysis of renewable energy trends"
        }
    ]
}'
```

### 8. Google Gemini Flash Latest
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gemini-flash-latest",
    "messages": [
        {
            "role": "user",
            "content": "Create a JSON schema for a user profile"
        }
    ]
}'
```

### 9. Legacy Model (Auto-mapped)
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gemini-1.5-flash",
    "messages": [
        {
            "role": "user",
            "content": "This should auto-map to gemini-2.5-flash"
        }
    ]
}'
```

### 10. Chat with Max Tokens
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gpt-4o-mini",
    "messages": [
        {
            "role": "user",
            "content": "Write a short story about space exploration"
        }
    ],
    "max_tokens": 500
}'
```

### 11. Multi-turn Conversation
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gemini-2.5-flash",
    "messages": [
        {
            "role": "user",
            "content": "What is React?"
        },
        {
            "role": "assistant", 
            "content": "React is a JavaScript library for building user interfaces."
        },
        {
            "role": "user",
            "content": "Can you give me a simple example?"
        }
    ]
}'
```

---

## 📊 **Analytics API**

### 12. Get Usage Statistics (API Key Auth)
```bash
curl --location 'http://localhost:8080/api/usage' \
--header 'Authorization: Bearer YOUR_API_KEY'
```

### 13. Get Usage Statistics with Date Range
```bash
curl --location 'http://localhost:8080/api/usage?days=7' \
--header 'Authorization: Bearer YOUR_API_KEY'
```

### 14. Get Usage Statistics (Last 24 hours)
```bash
curl --location 'http://localhost:8080/api/usage?days=1' \
--header 'Authorization: Bearer YOUR_API_KEY'
```

---

## 💳 **Credits & Payments API**

### 15. Get Current Credits Balance
```bash
curl --location 'http://localhost:8080/credits' \
--header 'Authorization: Bearer JWT_TOKEN'
```

### 16. Create Payment Order
```bash
curl --location 'http://localhost:8080/credits/order' \
--header 'Authorization: Bearer JWT_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
    "amount": 10.00
}'
```

### 17. Razorpay Webhook (Test)
```bash
curl --location 'http://localhost:8080/credits/add' \
--header 'Content-Type: application/json' \
--header 'X-Razorpay-Signature: GENERATED_SIGNATURE' \
--data '{
    "event": "payment.captured",
    "payload": {
        "payment": {
            "entity": {
                "id": "rzp_test_123456789",
                "amount": 10000,
                "status": "captured",
                "order_id": "order_123456789"
            }
        },
        "order": {
            "entity": {
                "id": "order_123456789",
                "amount": 10000,
                "notes": {
                    "user_id": "YOUR_USER_ID"
                }
            }
        }
    }
}'
```

---

## 🔐 **Authentication API**

### 18. User Signup
```bash
curl --location 'http://localhost:8080/auth/signup' \
--header 'Content-Type: application/json' \
--data '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
}'
```

### 19. User Login
```bash
curl --location 'http://localhost:8080/auth/login' \
--header 'Content-Type: application/json' \
--data '{
    "email": "test@example.com",
    "password": "password123"
}'
```

### 20. Email Verification
```bash
curl --location 'http://localhost:8080/auth/verify?token=VERIFICATION_TOKEN'
```

---

## 🔑 **API Keys Management**

### 21. Create API Key
```bash
curl --location 'http://localhost:8080/keys/create' \
--header 'Authorization: Bearer JWT_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
    "name": "My API Key"
}'
```

### 22. List API Keys
```bash
curl --location 'http://localhost:8080/keys' \
--header 'Authorization: Bearer JWT_TOKEN'
```

---

## 💬 **Chat History API**

### 23. Create New Chat
```bash
curl --location 'http://localhost:8080/newchat' \
--header 'Authorization: Bearer JWT_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
    "title": "My Chat Session"
}'
```

### 24. Get Chat History
```bash
curl --location 'http://localhost:8080/chathistory' \
--header 'Authorization: Bearer JWT_TOKEN'
```

### 25. Get Specific Chat
```bash
curl --location 'http://localhost:8080/chathistory/CHAT_ID' \
--header 'Authorization: Bearer JWT_TOKEN'
```

### 26. Delete Chat
```bash
curl --location 'http://localhost:8080/chathistory/CHAT_ID' \
--header 'Authorization: Bearer JWT_TOKEN' \
--request DELETE
```

---

## 🏥 **Health Check**

### 27. API Health Check
```bash
curl --location 'http://localhost:8080/'
```

**Expected Response:**
```json
{
    "status": "ok",
    "message": "ClearRouter API"
}
```

---

## 🧪 **Test Scenarios**

### Cost Calculation Test
```bash
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gemini-2.5-flash",
    "messages": [
        {
            "role": "user",
            "content": "Calculate the cost: Input tokens: 1000, Output tokens: 500. Model: gemini-2.5-flash. Input price: $0.30 per 1M tokens, Output price: $2.50 per 1M tokens. Show the calculation."
        }
    ]
}'
```

### Model Comparison Test
```bash
# Test same prompt with different models
curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "Explain photosynthesis"}]
}'

curl --location 'http://localhost:8080/v1/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gemini-2.5-flash", 
    "messages": [{"role": "user", "content": "Explain photosynthesis"}]
}'
```

---

## 📝 **Notes for Postman**

1. **Environment Variables:**
   - `BASE_URL`: `http://localhost:8080`
   - `API_KEY`: `HNNACLIs9au8CoUPeAIPyN--w-UZd8ASAKV61xRYy9I=`
   - `JWT_TOKEN`: Get from login response

2. **Headers to Set:**
   - `Content-Type`: `application/json`
   - `Authorization`: `Bearer {{API_KEY}}` or `Bearer {{JWT_TOKEN}}`

3. **Expected Status Codes:**
   - `200`: Success
   - `400`: Bad Request (invalid model, etc.)
   - `401`: Unauthorized (invalid API key/token)
   - `402`: Payment Required (insufficient credits)
   - `500`: Internal Server Error

4. **Cost Tracking:**
   After each chat completion, check `/api/usage` to see updated costs and token usage.

---

**Happy Testing! 🚀**
