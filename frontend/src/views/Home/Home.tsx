import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './Home.module.css';

interface AccountBalance {
  accountId: string;
  name: string;
  balance: number;
}

const API_BASE_URL = 'http://localhost:8000/api';

const Home: React.FC = () => {
  const [balances, setBalances] = useState<AccountBalance[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchBalances = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/account/balances`, { credentials: 'include' });
        if (!response.ok) throw new Error('Failed to fetch account balances');
        const data = await response.json();
        setBalances(data || []);
      } catch (err) {
        console.error('Error fetching balances:', err);
        setError('Failed to load account balances.');
      } finally {
        setLoading(false);
      }
    };

    fetchBalances();
  }, []);

  const handleAccountClick = (accountId: string) => {
    navigate('/transactions', { state: { selectedAccountId: accountId } });
  };

  const totalNetWorth = balances.reduce((sum, acc) => sum + acc.balance, 0);

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.netWorthCard}>
          <h3>Total Net Worth</h3>
          <p className={totalNetWorth >= 0 ? styles.positive : styles.negative}>
            ${totalNetWorth.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
          </p>
        </div>
      </header>

      <section className={styles.accountsGrid}>
        {loading ? (
          <p>Loading accounts...</p>
        ) : error ? (
          <p className={styles.error}>{error}</p>
        ) : balances.length > 0 ? (
          balances.map((acc) => (
            <div 
              key={acc.accountId} 
              className={styles.accountCard}
              onClick={() => handleAccountClick(acc.accountId)}
            >
              <h4>{acc.name}</h4>
              <p className={acc.balance >= 0 ? styles.positive : styles.negative}>
                ${acc.balance.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
              </p>
            </div>
          ))
        ) : (
          <p>No internal accounts found. Add one in the Accounts tab!</p>
        )}
      </section>
    </div>
  );
};

export default Home;
