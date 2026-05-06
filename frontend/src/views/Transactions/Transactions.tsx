import React, { useEffect, useState } from 'react';
import styles from './Transactions.module.css';

interface Account {
  accountId: string;
  name: string;
}

interface Transaction {
  transactionId: string;
  date: string;
  memo: string;
  amount: number;
  categoryId: string;
}

const API_BASE_URL = 'http://localhost:8000/api';

const Transactions: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [selectedAccountId, setSelectedAccountId] = useState<string>('');
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const isInitialMount = React.useRef(true);

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/account`, {
          credentials: 'include',
        });
        if (!response.ok) throw new Error('Failed to fetch accounts');
        const data: Account[] = await response.json();
        setAccounts(data);
        
        if (data.length > 0) {
          // If already set by parent/other logic, don't overwrite
          setSelectedAccountId(prev => prev || data[0].accountId);
        }
      } catch (err) {
        setError('Error connecting to API. Please ensure the backend is running.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchAccounts();
  }, []);

  useEffect(() => {
    // Prevent fetch on initial mount if ID is set by account fetcher
    if (isInitialMount.current) {
        isInitialMount.current = false;
        if (!selectedAccountId) return;
    }

    const fetchTransactions = async () => {
      if (!selectedAccountId) return;
      
      try {
        const response = await fetch(`${API_BASE_URL}/transaction?accountFrom=${selectedAccountId}`, {
          credentials: 'include',
        });
        if (!response.ok) throw new Error('Failed to fetch transactions');
        const data: Transaction[] = await response.json();
        setTransactions(data);
      } catch (err) {
        console.error('Error fetching transactions:', err);
      }
    };

    fetchTransactions();
  }, [selectedAccountId]);

  if (loading) return <div className={styles.container}>Loading...</div>;

  return (
    <div className={styles.container}>
      {error && <div style={{ color: '#ef4444', marginBottom: '1rem' }}>{error}</div>}
      
      <div className={styles.controls}>
        <div className={styles.selectGroup}>
          <label htmlFor="account-select">Account:</label>
          <select 
            id="account-select" 
            className={styles.select}
            value={selectedAccountId}
            onChange={(e) => setSelectedAccountId(e.target.value)}
          >
            <option value="">Select an account</option>
            {accounts.map(account => (
              <option key={account.accountId} value={account.accountId}>
                {account.name}
              </option>
            ))}
          </select>
        </div>
        <button className={styles.addButton}>+ Add Transaction</button>
      </div>

      <div className={styles.tableContainer}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Date</th>
              <th>Memo</th>
              <th>Category</th>
              <th style={{ textAlign: 'right' }}>Amount</th>
            </tr>
          </thead>
          <tbody>
            {transactions.length > 0 ? (
              transactions.map((tx) => (
                <tr key={tx.transactionId}>
                  <td>{tx.date ? new Date(tx.date).toLocaleDateString() : 'N/A'}</td>
                  <td>{tx.memo || 'No memo'}</td>
                  <td>{tx.categoryId}</td>
                  <td style={{ 
                    textAlign: 'right', 
                    color: tx.amount < 0 ? '#ef4444' : 'var(--accent)',
                    fontWeight: 600
                  }}>
                    {tx.amount < 0 ? `- $${Math.abs(tx.amount).toFixed(2)}` : `+ $${tx.amount.toFixed(2)}`}
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={4} style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-secondary)' }}>
                  No transactions found for this account.
                </td>
              </tr>
            )}
          </tbody>
        </table>
        
        <div className={styles.pagination}>
          <button className={styles.pageButton}>Previous</button>
          <button className={`${styles.pageButton} ${styles.pageButtonActive}`}>1</button>
          <button className={styles.pageButton}>Next</button>
        </div>
      </div>
    </div>
  );
};

export default Transactions;
