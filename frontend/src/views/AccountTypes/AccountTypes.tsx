import React, { useEffect, useState } from 'react';
import styles from './AccountTypes.module.css';
import Modal from '../../components/Modal/Modal';
import { useAuth } from '../../context/AuthContext';

interface AccountType {
  typeId: string;
  name: string;
  description: string;
}

const API_BASE_URL = 'http://localhost:8000/api';

const AccountTypes: React.FC = () => {
  const [accountTypes, setAccountTypes] = useState<AccountType[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingType, setEditingType] = useState<AccountType | null>(null);
  const [formData, setFormData] = useState({ name: '', description: '' });

  useEffect(() => {
    fetchAccountTypes();
  }, []);

  const fetchAccountTypes = async () => {
    const response = await fetch(`${API_BASE_URL}/account_type`, { credentials: 'include' });
    const data = await response.json();
    setAccountTypes(data || []);
  };

  const handleSave = async () => {
    const method = editingType ? 'PUT' : 'POST';
    const url = editingType ? `${API_BASE_URL}/account_type/${editingType.typeId}` : `${API_BASE_URL}/account_type`;
    
    await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formData),
      credentials: 'include',
    });
    
    setIsModalOpen(false);
    fetchAccountTypes();
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure?')) return;
    await fetch(`${API_BASE_URL}/account_type/${id}`, { method: 'DELETE', credentials: 'include' });
    fetchAccountTypes();
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2>Account Types</h2>
        <button className={styles.addButton} onClick={() => { setEditingType(null); setFormData({name: '', description: ''}); setIsModalOpen(true); }}>+ Add Type</button>
      </div>

      <table className={styles.table}>
        <thead>
          <tr><th>Name</th><th>Description</th><th>Actions</th></tr>
        </thead>
        <tbody>
          {accountTypes.map(t => (
            <tr key={t.typeId}>
              <td>{t.name}</td>
              <td>{t.description}</td>
              <td>
                <button onClick={() => { setEditingType(t); setFormData({name: t.name, description: t.description}); setIsModalOpen(true); }}>Edit</button>
                <button onClick={() => handleDelete(t.typeId)}>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingType ? "Edit Type" : "Add Type"}>
        <div className={styles.form}>
          <input value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} placeholder="Name" />
          <input value={formData.description} onChange={e => setFormData({...formData, description: e.target.value})} placeholder="Description" />
          <button onClick={handleSave}>Save</button>
        </div>
      </Modal>
    </div>
  );
};

export default AccountTypes;
