import React from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import styles from './Layout.module.css';

const Layout: React.FC = () => {
  const { user, logout, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  const navItems = [
    { name: 'Transactions', path: '/' },
    { name: 'Categories', path: '/categories' },
    { name: 'Accounts', path: '/accounts' },
    { name: 'Account Types', path: '/account-types' },
    { name: 'Reports', path: '/reports' },
  ];

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <div className={styles.layout}>
      <header className={styles.header}>
        <h1>moneyKeeper</h1>
        {isAuthenticated && user && (
          <div className={styles.userSection}>
            <span className={styles.userEmail}>{user.email}</span>
            <button className={styles.logoutButton} onClick={handleLogout}>
              Logout
            </button>
          </div>
        )}
      </header>

      {isAuthenticated && (
        <nav className={styles.nav}>
          <ul className={styles.navList}>
            {navItems.map((item) => (
              <li key={item.name}>
                <NavLink
                  to={item.path}
                  className={({ isActive }) =>
                    isActive ? `${styles.navItem} ${styles.activeNavItem}` : styles.navItem
                  }
                >
                  {item.name}
                </NavLink>
              </li>
            ))}
          </ul>
        </nav>
      )}

      <main className={styles.main}>
        <Outlet />
      </main>
...
      <footer className={styles.footer}>
        <p>&copy; {new Date().getFullYear()} moneyKeeper - Personal Finance Manager</p>
      </footer>
    </div>
  );
};

export default Layout;
