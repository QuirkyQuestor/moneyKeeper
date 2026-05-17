import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Layout from './components/Layout/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import Transactions from './views/Transactions/Transactions';
import Login from './views/Login/Login';
import AccountTypes from './views/AccountTypes/AccountTypes';
import Accounts from './views/Accounts/Accounts';
import Categories from './views/Categories/Categories';
import Reports from './views/Reports/Reports';

import Home from './views/Home/Home';

const App: React.FC = () => {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route path="login" element={<Login />} />
            
            <Route index element={
              <ProtectedRoute>
                <Home />
              </ProtectedRoute>
            } />

            <Route path="transactions" element={
              <ProtectedRoute>
                <Transactions />
              </ProtectedRoute>
            } />
            
            <Route path="categories" element={
              <ProtectedRoute>
                <Categories />
              </ProtectedRoute>
            } />
            
            <Route path="accounts" element={
              <ProtectedRoute>
                <Accounts />
              </ProtectedRoute>
            } />
            
            <Route path="account-types" element={
              <ProtectedRoute>
                <AccountTypes />
              </ProtectedRoute>
            } />
            
            <Route path="reports" element={
              <ProtectedRoute>
                <Reports />
              </ProtectedRoute>
            } />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
};

export default App;
