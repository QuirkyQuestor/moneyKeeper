import React, { useEffect, useState } from 'react';
import styles from './Accounts.module.css';
import Modal from '../../components/Modal/Modal';

interface AccountType {
  typeId: string;
  name: string;
}

interface Account {
  accountId: string;
  typeId: string;
  name: string;
  description: string;
  active: boolean;
}

const API_BASE_URL = 'http://localhost:8000/api';

const EditIcon = () => <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" /><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" /></svg>;
const DeleteIcon = () => <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /></svg>;

const Accounts: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [accountTypes, setAccountTypes] = useState<AccountType[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [editingAccount, setEditingAccount] = useState<Account | null>(null);
  const [formData, setFormData] = useState({ name: '', description: '', typeId: '', active: true });
  const [deletingId, setDeletingId] = useState<string | null>(null);

  useEffect(() => {
    fetchAccounts();
    fetchAccountTypes();
  }, []);

  const fetchAccounts = async () => {
    const response = await fetch(`${API_BASE_URL}/account`, { credentials: 'include' });
    const data = await response.json();
    setAccounts(data || []);
  };

  const fetchAccountTypes = async () => {
    const response = await fetch(`${API_BASE_URL}/account_type`, { credentials: 'include' });
    const data = await response.json();
    setAccountTypes(data || []);
  };

  const handleSave = async () => {
    const method = editingAccount ? 'PUT' : 'POST';
    const url = editingAccount ? `${API_BASE_URL}/account/${editingAccount.accountId}` : `${API_BASE_URL}/account`;
    await fetch(url, { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(formData), credentials: 'include' });
    setIsModalOpen(false);
    fetchAccounts();
  };

  const handleDelete = async () => {
    if (!deletingId) return;
    await fetch(`${API_BASE_URL}/account/${deletingId}`, { method: 'DELETE', credentials: 'include' });
    setIsDeleteModalOpen(false);
    fetchAccounts();
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2>Accounts</h2>
        <button className={styles.addButton} onClick={() => { setEditingAccount(null); setFormData({name: '', description: '', typeId: accountTypes[0]?.typeId || '', active: true}); setIsModalOpen(true); }}>+ Add Account</button>
      </div>

      <div className={styles.tableContainer}>
        <table className={styles.table}>
          <thead><tr><th>Name</th><th>Type</th><th>Description</th><th>Active</th><th>Actions</th></tr></thead>
          <tbody>
            {accounts.map(a => (
              <tr key={a.accountId}>
                <td>{a.name}</td>
                <td>{accountTypes.find(t => t.typeId === a.typeId)?.name || 'Unknown'}</td>
                <td>{a.description}</td>
                <td>{a.active ? 'Yes' : 'No'}</td>
                <td className={styles.actionCell}>
                  <button className={styles.iconButton} onClick={() => { setEditingAccount(a); setFormData({name: a.name, description: a.description, typeId: a.typeId, active: a.active}); setIsModalOpen(true); }}><EditIcon /></button>
                  <button className={styles.iconButton} onClick={() => { setDeletingId(a.accountId); setIsDeleteModalOpen(true); }}><DeleteIcon /></button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingAccount ? "Edit Account" : "Add Account"}>
        <div className={styles.form}>
          <input value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} placeholder="Name" />
          <input value={formData.description} onChange={e => setFormData({...formData, description: e.target.value})} placeholder="Description" />
          <select value={formData.typeId} onChange={e => setFormData({...formData, typeId: e.target.value})}>
            {accountTypes.map(t => <option key={t.typeId} value={t.typeId}>{t.name}</option>)}
          </select>
          <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
            <input 
              type="checkbox" 
              checked={formData.active} 
              onChange={e => setFormData({...formData, active: e.target.checked})} 
            />
            Active
          </label>
          <button onClick={handleSave}>Save</button>
        </div>
      </Modal>

      <Modal isOpen={isDeleteModalOpen} onClose={() => setIsDeleteModalOpen(false)} title="Confirm Delete">
        <p>Are you sure you want to delete this account?</p>
        <div className={styles.confirmModalActions}>
          <button className={styles.btnCancel} onClick={() => setIsDeleteModalOpen(false)}>Cancel</button>
          <button className={styles.btnDelete} onClick={handleDelete}>Delete</button>
        </div>
      </Modal>
    </div>
  );
};

export default Accounts;
