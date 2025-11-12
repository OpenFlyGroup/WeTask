# Using `/api/users/me` Endpoint

The `/api/users/me` endpoint returns the current authenticated user's information.

## Endpoint Details

- **URL**: `GET /api/users/me`
- **Authentication**: Required (JWT Bearer token)
- **Response**: User object with id, email, name, createdAt, updatedAt

## Usage Examples

### 1. Using cURL

```bash
# First, login to get a token
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response will include accessToken:
# {
#   "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   ...
# }

# Then use the token to get current user
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 2. JavaScript/TypeScript with Fetch

```javascript
// Get current user
async function getCurrentUser() {
  const token = localStorage.getItem('accessToken');
  
  const response = await fetch('http://localhost:3000/api/users/me', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
  
  if (!response.ok) {
    if (response.status === 401) {
      // Token expired, try to refresh
      await refreshToken();
      return getCurrentUser(); // Retry
    }
    throw new Error('Failed to get user');
  }
  
  const user = await response.json();
  return user;
}

// Usage
getCurrentUser().then(user => {
  console.log('Current user:', user);
  // {
  //   id: 1,
  //   email: "user@example.com",
  //   name: "John Doe",
  //   createdAt: "2024-01-01T00:00:00Z",
  //   updatedAt: "2024-01-01T00:00:00Z"
  // }
});
```

### 3. JavaScript/TypeScript with Axios

```javascript
import axios from 'axios';

// Create axios instance
const api = axios.create({
  baseURL: 'http://localhost:3000/api',
});

// Add token to all requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token refresh on 401
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const refreshToken = localStorage.getItem('refreshToken');
      if (refreshToken) {
        try {
          const { data } = await axios.post(
            'http://localhost:3000/api/auth/refresh',
            { refreshToken }
          );
          localStorage.setItem('accessToken', data.accessToken);
          localStorage.setItem('refreshToken', data.refreshToken);
          
          // Retry original request
          error.config.headers.Authorization = `Bearer ${data.accessToken}`;
          return axios.request(error.config);
        } catch {
          // Refresh failed, redirect to login
          window.location.href = '/login';
        }
      }
    }
    return Promise.reject(error);
  }
);

// Get current user
async function getCurrentUser() {
  try {
    const { data } = await api.get('/users/me');
    return data;
  } catch (error) {
    console.error('Error getting user:', error);
    throw error;
  }
}

// Usage
getCurrentUser().then(user => {
  console.log('Current user:', user);
});
```

### 4. React Hook Example

```jsx
import { useState, useEffect } from 'react';
import axios from 'axios';

function useCurrentUser() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const token = localStorage.getItem('accessToken');
        if (!token) {
          setLoading(false);
          return;
        }

        const { data } = await axios.get('http://localhost:3000/api/users/me', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        
        setUser(data);
      } catch (err) {
        if (err.response?.status === 401) {
          // Token expired, try refresh
          const refreshToken = localStorage.getItem('refreshToken');
          if (refreshToken) {
            try {
              const { data } = await axios.post(
                'http://localhost:3000/api/auth/refresh',
                { refreshToken }
              );
              localStorage.setItem('accessToken', data.accessToken);
              localStorage.setItem('refreshToken', data.refreshToken);
              
              // Retry
              const { data: userData } = await axios.get(
                'http://localhost:3000/api/users/me',
                {
                  headers: {
                    Authorization: `Bearer ${data.accessToken}`,
                  },
                }
              );
              setUser(userData);
            } catch {
              setError('Session expired');
              localStorage.removeItem('accessToken');
              localStorage.removeItem('refreshToken');
            }
          } else {
            setError('Not authenticated');
          }
        } else {
          setError(err.message);
        }
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, []);

  return { user, loading, error };
}

// Usage in component
function UserProfile() {
  const { user, loading, error } = useCurrentUser();

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!user) return <div>Not logged in</div>;

  return (
    <div>
      <h1>Profile</h1>
      <p>ID: {user.id}</p>
      <p>Email: {user.email}</p>
      <p>Name: {user.name}</p>
      <p>Member since: {new Date(user.createdAt).toLocaleDateString()}</p>
    </div>
  );
}
```

### 5. Complete React Context Example

```jsx
import { createContext, useContext, useState, useEffect } from 'react';
import axios from 'axios';

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const api = axios.create({
    baseURL: 'http://localhost:3000/api',
  });

  // Add auth token to requests
  api.interceptors.request.use((config) => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });

  // Handle token refresh
  api.interceptors.response.use(
    (response) => response,
    async (error) => {
      if (error.response?.status === 401) {
        const refreshToken = localStorage.getItem('refreshToken');
        if (refreshToken) {
          try {
            const { data } = await axios.post(
              'http://localhost:3000/api/auth/refresh',
              { refreshToken }
            );
            localStorage.setItem('accessToken', data.accessToken);
            localStorage.setItem('refreshToken', data.refreshToken);
            error.config.headers.Authorization = `Bearer ${data.accessToken}`;
            return axios.request(error.config);
          } catch {
            logout();
          }
        } else {
          logout();
        }
      }
      return Promise.reject(error);
    }
  );

  const getCurrentUser = async () => {
    try {
      const { data } = await api.get('/users/me');
      setUser(data);
      return data;
    } catch (error) {
      setUser(null);
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    setUser(null);
  };

  useEffect(() => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      getCurrentUser().finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  return (
    <AuthContext.Provider value={{ user, getCurrentUser, logout, loading, api }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}

// Usage
function App() {
  return (
    <AuthProvider>
      <UserProfile />
    </AuthProvider>
  );
}

function UserProfile() {
  const { user, loading } = useAuth();

  if (loading) return <div>Loading...</div>;
  if (!user) return <div>Please login</div>;

  return (
    <div>
      <h1>Welcome, {user.name}!</h1>
      <p>Email: {user.email}</p>
    </div>
  );
}
```

### 6. Go Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

func getCurrentUser(token string) (*User, error) {
    req, err := http.NewRequest("GET", "http://localhost:3000/api/users/me", nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 401 {
        // Token expired, refresh it
        // ... refresh logic ...
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var user User
    if err := json.Unmarshal(body, &user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

type User struct {
    ID        uint   `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
}
```

## Response Format

### Success Response (200 OK)

```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Error Responses

#### 401 Unauthorized (No token or invalid token)

```json
{
  "error": "Invalid token"
}
```

#### 404 Not Found (User doesn't exist)

```json
{
  "error": "User not found"
}
```

#### 500 Internal Server Error

```json
{
  "error": "Internal server error message"
}
```

## Important Notes

1. **Authentication Required**: You must include a valid JWT token in the `Authorization` header
2. **Token Format**: `Bearer <token>`
3. **Token Expiration**: Access tokens expire after 15 minutes - implement token refresh
4. **User ID**: The endpoint automatically extracts the user ID from the JWT token, so you don't need to pass it
5. **CORS**: Make sure your frontend is allowed to make requests (CORS is configured in the API Gateway)

## Quick Test

```bash
# 1. Register or login
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' \
  | jq -r '.accessToken')

# 2. Get current user
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

