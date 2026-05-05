import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import styles from './Login.module.css';

const Login: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  
  const { login, register } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const from = location.state?.from?.pathname || '/';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      if (isLogin) {
        await login(email, password);
        navigate(from, { replace: true });
      } else {
        await register(email, password);
        // After registration, auto-login or switch to login
        await login(email, password);
        navigate(from, { replace: true });
      }
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <h2>{isLogin ? 'Login to moneyKeeper' : 'Create Account'}</h2>
        
        {error && <div className={styles.error}>{error}</div>}
        
        <form className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.inputGroup}>
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              className={styles.input}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className={styles.inputGroup}>
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              className={styles.input}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button type="submit" className={styles.button}>
            {isLogin ? 'Login' : 'Register'}
          </button>
        </form>
        
        <div className={styles.toggle}>
          {isLogin ? "Don't have an account?" : "Already have an account?"}
          <button 
            className={styles.toggleLink}
            onClick={() => setIsLogin(!isLogin)}
          >
            {isLogin ? 'Register now' : 'Login here'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default Login;
