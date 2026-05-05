import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Layout from './components/Layout/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import Transactions from './views/Transactions/Transactions';
import Login from './views/Login/Login';
import WIP from './views/WIP/WIP';

const App: React.FC = () => {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route path="login" element={<Login />} />
            
            <Route index element={
              <ProtectedRoute>
                <Transactions />
              </ProtectedRoute>
            } />
            
            <Route path="categories" element={
              <ProtectedRoute>
                <WIP name="Categories" />
              </ProtectedRoute>
            } />
            
            <Route path="accounts" element={
              <ProtectedRoute>
                <WIP name="Accounts" />
              </ProtectedRoute>
            } />
            
            <Route path="account-types" element={
              <ProtectedRoute>
                <WIP name="Account Types" />
              </ProtectedRoute>
            } />
            
            <Route path="reports" element={
              <ProtectedRoute>
                <WIP name="Reports" />
              </ProtectedRoute>
            } />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
};

export default App;
