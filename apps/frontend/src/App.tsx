
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import Landing from './components/Landing';
import Login from './components/Login';
import Signup from './components/Signup';
import Dashboard from './components/Dashboard';
import DashboardHome from './components/DashboardHome';
import ApiKeys from './components/ApiKeys';
import Chat from './components/Chat';
import AddCredits from './components/AddCredits';
import Models from './components/Models';
import ProtectedRoute from './components/ProtectedRoute';


function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<Signup />} />
          <Route path="/models" element={<Models />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          >
            <Route index element={<DashboardHome />} />
            <Route path="keys" element={<ApiKeys />} />
            <Route path="chat" element={<Chat />} />
            <Route path="models" element={<Models />} />
            <Route path="credits" element={<AddCredits />} />
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
