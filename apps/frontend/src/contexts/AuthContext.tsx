import React, { createContext, useContext, useReducer, useEffect, useCallback } from 'react';
import type { ReactNode } from 'react';
import { authAPI, setAuthToken, removeAuthToken } from '../services/api';
import type { User, LoginRequest, SignupRequest } from '../services/api';

// ----------------- Auth state interface -----------------
interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
}

// ----------------- Auth actions -----------------
type AuthAction =
  | { type: 'LOGIN_START' }
  | { type: 'LOGIN_SUCCESS'; payload: { user: User; token: string } }
  | { type: 'LOGIN_FAILURE'; payload: string }
  | { type: 'SIGNUP_START' }
  | { type: 'SIGNUP_SUCCESS' }
  | { type: 'SIGNUP_FAILURE'; payload: string }
  | { type: 'LOGOUT' }
  | { type: 'CLEAR_ERROR' }
  | { type: 'RESTORE_AUTH'; payload: { user: User; token: string } };

// ----------------- Auth context interface -----------------
interface AuthContextType {
  state: AuthState;
  login: (credentials: LoginRequest) => Promise<void>;
  signup: (userData: SignupRequest) => Promise<void>;
  logout: () => void;
  clearError: () => void;
}

// ----------------- Initial state -----------------
const initialState: AuthState = {
  user: null,
  token: null,
  isLoading: false,
  error: null,
};

// ----------------- Reducer -----------------
const authReducer = (state: AuthState, action: AuthAction): AuthState => {
  switch (action.type) {
    case 'LOGIN_START':
    case 'SIGNUP_START':
      return { ...state, isLoading: true, error: null };

    case 'LOGIN_SUCCESS':
      return {
        ...state,
        isLoading: false,
        user: action.payload.user,
        token: action.payload.token,
        error: null,
      };

    case 'LOGIN_FAILURE':
    case 'SIGNUP_FAILURE':
      return { ...state, isLoading: false, error: action.payload };

    case 'SIGNUP_SUCCESS':
      return { ...state, isLoading: false, error: null };

    case 'LOGOUT':
      return { ...state, user: null, token: null, error: null, isLoading: false };

    case 'CLEAR_ERROR':
      return { ...state, error: null };

    case 'RESTORE_AUTH':
      return {
        ...state,
        user: action.payload.user,
        token: action.payload.token,
        isLoading: false,
        error: null,
      };

    default:
      return state;
  }
};

// ----------------- Context -----------------
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// ----------------- Provider -----------------
interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(authReducer, initialState);

  // Restore auth state from localStorage on app start
  useEffect(() => {
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');

    if (token && userStr) {
      try {
        const user = JSON.parse(userStr);
        setAuthToken(token);
        dispatch({ type: 'RESTORE_AUTH', payload: { user, token } });
      } catch {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }
    }
  }, []);

  // ----------------- Login -----------------
  const login = async (credentials: LoginRequest) => {
    dispatch({ type: 'LOGIN_START' });
    try {
      const response = await authAPI.login(credentials);

      // Save to localStorage
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      setAuthToken(response.token);

      dispatch({
        type: 'LOGIN_SUCCESS',
        payload: { user: response.user, token: response.token },
      });
    } catch (error: any) {
      const errorMessage =
        typeof error.response?.data?.error === 'string'
          ? error.response.data.error
          : 'Login failed';
      dispatch({ type: 'LOGIN_FAILURE', payload: errorMessage });
      return Promise.reject(errorMessage);
    }
  };

  // ----------------- Signup -----------------
  const signup = async (userData: SignupRequest) => {
    dispatch({ type: 'SIGNUP_START' });
    try {
      await authAPI.signup(userData);
      dispatch({ type: 'SIGNUP_SUCCESS' });
    } catch (error: any) {
      const errorMessage =
        typeof error.response?.data?.error === 'string'
          ? error.response.data.error
          : 'Signup failed';
      dispatch({ type: 'SIGNUP_FAILURE', payload: errorMessage });
      return Promise.reject(errorMessage);
    }
  };

  // ----------------- Logout -----------------
  const logout = () => {
    removeAuthToken();
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    dispatch({ type: 'LOGOUT' });
  };

  // ----------------- Clear Error -----------------
  const clearError = useCallback(() => {
    dispatch({ type: 'CLEAR_ERROR' });
  }, []);

  // Context value
  const value: AuthContextType = {
    state,
    login,
    signup,
    logout,
    clearError,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// ----------------- Custom hook -----------------
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export default AuthContext;