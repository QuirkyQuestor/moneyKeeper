import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout/Layout';
import Transactions from './views/Transactions/Transactions';
import WIP from './views/WIP/WIP';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Transactions />} />
          <Route path="categories" element={<WIP name="Categories" />} />
          <Route path="accounts" element={<WIP name="Accounts" />} />
          <Route path="account-types" element={<WIP name="Account Types" />} />
          <Route path="reports" element={<WIP name="Reports" />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
