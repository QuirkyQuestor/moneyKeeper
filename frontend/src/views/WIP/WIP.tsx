import React from 'react';

const WIP: React.FC<{ name: string }> = ({ name }) => {
  return (
    <div style={{ textAlign: 'center', marginTop: '4rem' }}>
      <h2 style={{ fontSize: '2rem', color: 'var(--text-secondary)' }}>Work in progress: {name}</h2>
      <p style={{ color: 'var(--text-secondary)', marginTop: '1rem' }}>We're working hard to bring this feature to you.</p>
    </div>
  );
};

export default WIP;
