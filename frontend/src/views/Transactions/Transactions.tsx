import React, { useEffect, useState, useCallback } from 'react';
import styles from './Transactions.module.css';
import Modal from '../../components/Modal/Modal';

interface Account {
  accountId: string;
  name: string;
  isExternal: boolean;
}

interface Category {
  categoryId: string;
  fullName: string;
  expence: boolean;
}

interface Transaction {
  transactionId: string;
  date: string;
  memo: string;
  amount: number;
  categoryId: string;
  accountFrom: string;
  accountTo: string;
}

interface PaginatedResponse {
  transactions: Transaction[];
  totalCount: number;
}

const API_BASE_URL = 'http://localhost:8000/api';
const PAGE_SIZE = 50;

const EditIcon = () => <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" /><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" /></svg>;
const DeleteIcon = () => <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /></svg>;

const Transactions: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedAccountId, setSelectedAccountId] = useState<string>('');
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [totalCount, setTotalCount] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState<Transaction | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    date: new Date().toISOString().split('T')[0],
    memo: '',
    amount: 0,
    categoryId: '',
    accountFrom: '',
    accountTo: '',
    isTransfer: false
  });

  const fetchTransactions = useCallback(async () => {
    setLoading(true);
    try {
      const offset = currentPage * PAGE_SIZE;
      let url = `${API_BASE_URL}/transaction?limit=${PAGE_SIZE}&offset=${offset}`;
      if (selectedAccountId) {
        url += `&accountFrom=${selectedAccountId}`;
      }
      
      const response = await fetch(url, { credentials: 'include' });
      if (!response.ok) throw new Error('Failed to fetch transactions');
      const data: PaginatedResponse = await response.json();
      setTransactions(data.transactions || []);
      setTotalCount(data.totalCount || 0);
    } catch (err) {
      console.error('Error fetching transactions:', err);
      setError('Failed to load transactions.');
    } finally {
      setLoading(false);
    }
  }, [selectedAccountId, currentPage]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [accRes, catRes] = await Promise.all([
          fetch(`${API_BASE_URL}/account`, { credentials: 'include' }),
          fetch(`${API_BASE_URL}/category`, { credentials: 'include' })
        ]);

        if (!accRes.ok || !catRes.ok) throw new Error('Failed to fetch auxiliary data');
        
        const accData: Account[] = await accRes.json();
        const catData: Category[] = await catRes.json();
        
        setAccounts(accData);
        setCategories(catData);
      } catch (err) {
        console.error('Error fetching auxiliary data:', err);
        setError('Error connecting to API. Please ensure the backend is running.');
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    fetchTransactions();
  }, [fetchTransactions]);

  const handleSave = async () => {
    try {
      const method = editingTransaction ? 'PUT' : 'POST';
      const url = editingTransaction 
        ? `${API_BASE_URL}/transaction/${editingTransaction.transactionId}` 
        : `${API_BASE_URL}/transaction`;
      
      const category = categories.find(c => c.categoryId === formData.categoryId);
      let amount = Math.abs(parseFloat(formData.amount.toString()));
      if (category && category.expence) {
        amount = -amount;
      }
      
      const payload = {
        ...formData,
        amount,
        date: new Date(formData.date).toISOString(),
        accountTo: formData.accountTo
      };

      const response = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
        credentials: 'include'
      });

      if (!response.ok) throw new Error('Failed to save transaction');

      setIsModalOpen(false);
      fetchTransactions();
    } catch (err) {
      console.error('Error saving transaction:', err);
      alert('Failed to save transaction.');
    }
  };

  const handleDelete = async () => {
    if (!deletingId) return;
    try {
      const response = await fetch(`${API_BASE_URL}/transaction/${deletingId}`, {
        method: 'DELETE',
        credentials: 'include'
      });
      if (!response.ok) throw new Error('Failed to delete transaction');
      setIsDeleteModalOpen(false);
      fetchTransactions();
    } catch (err) {
      console.error('Error deleting transaction:', err);
      alert('Failed to delete transaction.');
    }
  };

  const openAddModal = () => {
    setEditingTransaction(null);
    setFormData({
      date: new Date().toISOString().split('T')[0],
      memo: '',
      amount: 0,
      categoryId: categories[0]?.categoryId || '',
      accountFrom: selectedAccountId || accounts.find(a => !a.isExternal)?.accountId || '',
      accountTo: '',
      isTransfer: false
    });
    setIsModalOpen(true);
  };

  const openEditModal = (tx: Transaction) => {
    setEditingTransaction(tx);
    const targetAccount = accounts.find(a => a.accountId === tx.accountTo);
    setFormData({
      date: tx.date ? tx.date.split('T')[0] : new Date().toISOString().split('T')[0],
      memo: tx.memo,
      amount: Math.abs(tx.amount),
      categoryId: tx.categoryId,
      accountFrom: tx.accountFrom,
      accountTo: tx.accountTo,
      isTransfer: targetAccount ? !targetAccount.isExternal : false
    });
    setIsModalOpen(true);
  };

  const totalPages = Math.ceil(totalCount / PAGE_SIZE);

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
            onChange={(e) => {
              setSelectedAccountId(e.target.value);
              setCurrentPage(0);
            }}
          >
            <option value="">All Accounts</option>
            {accounts.filter(a => !a.isExternal).map(account => (
              <option key={account.accountId} value={account.accountId}>
                {account.name}
              </option>
            ))}
          </select>
        </div>
        <button className={styles.addButton} onClick={openAddModal}>+ Add Transaction</button>
      </div>

      <div className={styles.tableContainer}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Date</th>
              <th>Account</th>
              <th style={{ textAlign: 'right' }}>Amount</th>
              <th>Recipient</th>
              <th>Category</th>
              <th>Memo</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {!loading && transactions.length > 0 ? (
              transactions.map((tx) => (
                <tr key={tx.transactionId}>
                  <td>{tx.date ? new Date(tx.date).toLocaleDateString() : 'N/A'}</td>
                  <td>{accounts.find(a => a.accountId === tx.accountFrom)?.name || 'Unknown'}</td>
                  <td style={{ 
                    textAlign: 'right', 
                    color: tx.amount < 0 ? '#ef4444' : 'var(--accent)',
                    fontWeight: 600
                  }}>
                    {tx.amount < 0 ? `- $${Math.abs(tx.amount).toFixed(2)}` : `+ $${tx.amount.toFixed(2)}`}
                  </td>
                  <td>{accounts.find(a => a.accountId === tx.accountTo)?.name || 'Unknown'}</td>
                  <td>{categories.find(c => c.categoryId === tx.categoryId)?.fullName || 'Unknown'}</td>
                  <td>{tx.memo || 'No memo'}</td>
                  <td className={styles.actionCell}>
                    <button className={styles.iconButton} onClick={() => openEditModal(tx)}><EditIcon /></button>
                    <button className={styles.iconButton} onClick={() => { setDeletingId(tx.transactionId); setIsDeleteModalOpen(true); }}><DeleteIcon /></button>
                  </td>
                </tr>
              ))
            ) : !loading ? (
              <tr>
                <td colSpan={7} style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-secondary)' }}>
                  No transactions found.
                </td>
              </tr>
            ) : (
              <tr>
                <td colSpan={7} style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-secondary)' }}>
                  Loading...
                </td>
              </tr>
            )}
          </tbody>
        </table>
        
        {totalPages > 1 && (
          <div className={styles.pagination}>
            <span className={styles.pageInfo}>
              Page {currentPage + 1} of {totalPages} ({totalCount} total)
            </span>
            <button 
              className={styles.pageButton} 
              disabled={currentPage === 0}
              onClick={() => setCurrentPage(prev => prev - 1)}
            >
              Previous
            </button>
            <button 
              className={styles.pageButton} 
              disabled={currentPage >= totalPages - 1}
              onClick={() => setCurrentPage(prev => prev + 1)}
            >
              Next
            </button>
          </div>
        )}
      </div>

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingTransaction ? "Edit Transaction" : "Add Transaction"}>
        <div className={styles.form}>
          <input type="date" value={formData.date} onChange={e => setFormData({...formData, date: e.target.value})} />
          <input type="number" step="0.01" value={formData.amount} onChange={e => setFormData({...formData, amount: parseFloat(e.target.value) || 0})} placeholder="Amount" />
          
          <select value={formData.accountFrom} onChange={e => setFormData({...formData, accountFrom: e.target.value})}>
            <option value="" disabled>Select Source Account</option>
            {accounts.filter(a => !a.isExternal).map(a => <option key={a.accountId} value={a.accountId}>{a.name}</option>)}
          </select>

          <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
            <input 
              type="checkbox" 
              checked={formData.isTransfer} 
              onChange={e => setFormData({...formData, isTransfer: e.target.checked})} 
            />
            Transfer to internal account
          </label>

          <select 
            value={formData.accountTo} 
            onChange={e => setFormData({...formData, accountTo: e.target.value})}
          >
            <option value="" disabled>Select Target</option>
            {accounts
              .filter(a => formData.isTransfer ? !a.isExternal : a.isExternal)
              .map(a => <option key={a.accountId} value={a.accountId}>{a.name}</option>)
            }
          </select>

          <select value={formData.categoryId} onChange={e => setFormData({...formData, categoryId: e.target.value})}>
            <option value="" disabled>Select Category</option>
            {categories.map(c => <option key={c.categoryId} value={c.categoryId}>{c.fullName}</option>)}
          </select>

          <textarea value={formData.memo} onChange={e => setFormData({...formData, memo: e.target.value})} placeholder="Memo" />
          
          <button onClick={handleSave}>Save</button>
        </div>
      </Modal>

      <Modal isOpen={isDeleteModalOpen} onClose={() => setIsDeleteModalOpen(false)} title="Confirm Delete">
        <p>Are you sure you want to delete this transaction?</p>
        <div className={styles.confirmModalActions}>
          <button className={styles.btnCancel} onClick={() => setIsDeleteModalOpen(false)}>Cancel</button>
          <button className={styles.btnDelete} onClick={handleDelete}>Delete</button>
        </div>
      </Modal>
    </div>
  );
};

export default Transactions;
