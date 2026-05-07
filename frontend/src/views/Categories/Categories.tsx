import React, { useEffect, useState } from 'react';
import styles from './Categories.module.css';
import Modal from '../../components/Modal/Modal';

interface Category {
  categoryId: string;
  parentId: string | null;
  name: string;
  fullName: string;
  description: string;
  expence: boolean;
}

const API_BASE_URL = 'http://localhost:8000/api';

const EditIcon = () => <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" /><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" /></svg>;
const DeleteIcon = () => <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /></svg>;

const Categories: React.FC = () => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [formData, setFormData] = useState({ name: '', description: '', parentId: '', expence: true });
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/category`, { credentials: 'include' });
      const data = await response.json();
      setCategories(data || []);
    } catch (err) {
      console.error('Failed to fetch categories:', err);
    }
  };

  const handleSave = async () => {
    const method = editingCategory ? 'PUT' : 'POST';
    const url = editingCategory ? `${API_BASE_URL}/category/${editingCategory.categoryId}` : `${API_BASE_URL}/category`;
    
    // Ensure parentId is null if empty string
    const payload = {
      ...formData,
      parentId: formData.parentId === '' ? null : formData.parentId
    };

    try {
      const response = await fetch(url, { 
        method, 
        headers: { 'Content-Type': 'application/json' }, 
        body: JSON.stringify(payload), 
        credentials: 'include' 
      });

      if (response.status === 409) {
        setError('A category with this name already exists.');
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to save category');
      }

      setIsModalOpen(false);
      setError(null);
      fetchCategories();
    } catch (err) {
      setError('An error occurred while saving.');
      console.error(err);
    }
  };

  const handleDelete = async () => {
    if (!deletingId) return;
    try {
      const response = await fetch(`${API_BASE_URL}/category/${deletingId}`, { 
        method: 'DELETE', 
        credentials: 'include' 
      });

      if (!response.ok) {
        // In a real app, we'd handle 409/conflict if there are children/transactions
        const data = await response.json();
        alert(data.error || 'Cannot delete category. It might have subcategories or transactions.');
        return;
      }

      setIsDeleteModalOpen(false);
      fetchCategories();
    } catch (err) {
      console.error('Failed to delete category:', err);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2>Categories</h2>
        <button className={styles.addButton} onClick={() => { 
          setEditingCategory(null); 
          setFormData({name: '', description: '', parentId: '', expence: true}); 
          setError(null);
          setIsModalOpen(true); 
        }}>+ Add Category</button>
      </div>

      <div className={styles.tableContainer}>
        <table className={styles.table}>
          <thead><tr><th>Full Name</th><th>Description</th><th>Type</th><th>Actions</th></tr></thead>
          <tbody>
            {categories.map(c => (
              <tr key={c.categoryId}>
                <td>{c.fullName}</td>
                <td>{c.description}</td>
                <td>{c.expence ? 'Expense' : 'Income'}</td>
                <td className={styles.actionCell}>
                  <button className={styles.iconButton} onClick={() => { 
                    setEditingCategory(c); 
                    setFormData({name: c.name, description: c.description, parentId: c.parentId || '', expence: c.expence}); 
                    setError(null);
                    setIsModalOpen(true); 
                  }}><EditIcon /></button>
                  <button className={styles.iconButton} onClick={() => { setDeletingId(c.categoryId); setIsDeleteModalOpen(true); }}><DeleteIcon /></button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={editingCategory ? "Edit Category" : "Add Category"}>
        <div className={styles.form}>
          {error && <p style={{ color: '#ef4444', fontSize: '0.875rem' }}>{error}</p>}
          <input value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} placeholder="Name" />
          <input value={formData.description} onChange={e => setFormData({...formData, description: e.target.value})} placeholder="Description" />
          
          <select value={formData.parentId} onChange={e => setFormData({...formData, parentId: e.target.value})}>
            <option value="">None (Top Level)</option>
            {categories
              .filter(c => !editingCategory || (c.categoryId !== editingCategory.categoryId)) // Basic circular ref prevention
              .map(c => <option key={c.categoryId} value={c.categoryId}>{c.fullName}</option>)
            }
          </select>

          <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
            <input 
              type="checkbox" 
              checked={formData.expence} 
              onChange={e => setFormData({...formData, expence: e.target.checked})} 
            />
            Expense
          </label>
          <button onClick={handleSave}>Save</button>
        </div>
      </Modal>

      <Modal isOpen={isDeleteModalOpen} onClose={() => setIsDeleteModalOpen(false)} title="Confirm Delete">
        <p>Are you sure you want to delete this category?</p>
        <p style={{ fontSize: '0.875rem', color: 'var(--text-secondary)', marginTop: '0.5rem' }}>
          Deleting a category is only possible if it has no subcategories and no transactions.
        </p>
        <div className={styles.confirmModalActions}>
          <button className={styles.btnCancel} onClick={() => setIsDeleteModalOpen(false)}>Cancel</button>
          <button className={styles.btnDelete} onClick={handleDelete}>Delete</button>
        </div>
      </Modal>
    </div>
  );
};

export default Categories;
