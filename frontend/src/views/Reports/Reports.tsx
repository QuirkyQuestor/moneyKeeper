import React, { useEffect, useState } from 'react';
import styles from './Reports.module.css';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement, LineElement } from 'chart.js';
import { Pie, Bar, Line } from 'react-chartjs-2';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement, PointElement, LineElement);

const API_BASE_URL = 'http://localhost:8000/api';

const Reports: React.FC = () => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API_BASE_URL}/reports`, { credentials: 'include' })
      .then(res => res.json())
      .then(setData)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className={styles.container}>Loading reports...</div>;

  const pieData = {
    labels: data.expensesByCategory.map((c: any) => c.categoryName),
    datasets: [{
      data: data.expensesByCategory.map((c: any) => c.amount),
      backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF']
    }]
  };

  const barData = {
    labels: data.monthlyComparison.map((m: any) => m.month),
    datasets: [
      { label: 'Income', data: data.monthlyComparison.map((m: any) => m.income), backgroundColor: '#4BC0C0' },
      { label: 'Expenses', data: data.monthlyComparison.map((m: any) => m.expenses), backgroundColor: '#FF6384' }
    ]
  };

  const lineData = {
    labels: data.netWorthTrend.map((p: any) => p.date),
    datasets: [{
      label: 'Balance',
      data: data.netWorthTrend.map((p: any) => p.balance),
      borderColor: '#36A2EB',
      fill: false
    }]
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}><h2>Financial Reports</h2></div>
      <div className={styles.chartsGrid}>
        <div className={styles.chartCard}>
          <h3>Expenses by Category</h3>
          <Pie data={pieData} />
        </div>
        <div className={styles.chartCard}>
          <h3>Income vs Expenses</h3>
          <Bar data={barData} />
        </div>
        <div className={styles.chartCard}>
          <h3>Net Worth Trend</h3>
          <Line data={lineData} />
        </div>
      </div>
    </div>
  );
};

export default Reports;
