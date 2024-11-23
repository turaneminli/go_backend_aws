import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Chart as ChartJS, LineElement, LinearScale, CategoryScale, PointElement, Tooltip, Legend } from 'chart.js';
import { Line } from 'react-chartjs-2';
import styles from './css/Dashboard.module.css';
import { toast } from 'react-toastify';  // Import toast from react-toastify
import 'react-toastify/dist/ReactToastify.css';  // Import the CSS for Toastify

// Register necessary Chart.js components
ChartJS.register(LineElement, LinearScale, CategoryScale, PointElement, Tooltip, Legend);

function Dashboard() {
  const [metrics, setMetrics] = useState([]);
  const [selectedInstanceId, setSelectedInstanceId] = useState('');
  const [instances, setInstances] = useState([]);
  const [loading, setLoading] = useState(false);

  // Fetch instances for dropdown
  useEffect(() => {
    axios.get('http://localhost:8080/instances/status')
      .then((response) => {

        const data = response.data.map((instance) => ({
          name: instance.name,
          id: instance.id,
        }));
        console.log("RESPONSE: ", response.data)

        setInstances(data);
      })
      .catch((error) => {
        toast.error('Error fetching instances:', error);
      });
  }, []);

  // Fetch metrics from the backend
  const fetchMetrics = () => {
    if (!selectedInstanceId) return;

    setLoading(true);
    axios.get(`http://localhost:8080/cloudwatch/metrics?instanceId=${selectedInstanceId}`)
      .then((response) => {
        if (!response.data) {
            throw new Error("No data for this instance!")
        }
        setMetrics(response.data);
        setLoading(false);
      })
      .catch((error) => {
        const errorMessage = error.response?.data?.message || error.message || 'An unknown error occurred';
        toast.error(`Error fetching metrics: ${errorMessage}`);
        setLoading(false);
        setLoading(false);
      });
  };

  // Process metrics into chart-ready format
  const processData = (metricName) => {
    const filteredMetrics = metrics.filter((metric) => metric.metric_name === metricName);
    return {
      labels: filteredMetrics.map((metric) => new Date(metric.timestamp).toLocaleTimeString()),
      datasets: [
        {
          label: metricName,
          data: filteredMetrics.map((metric) => metric.value),
          borderColor: metricName === 'NetworkIn' ? '#28a745' : metricName === 'NetworkOut' ? '#dc3545' : '#007bff',
          backgroundColor: metricName === 'NetworkIn' ? 'rgba(40, 167, 69, 0.2)' : metricName === 'NetworkOut' ? 'rgba(220, 53, 69, 0.2)' : 'rgba(0, 123, 255, 0.2)',
          fill: true,
          tension: 0.4,
        },
      ],
    };
  };
  

  return (
    <div className={styles.dashboard}>
      <h1>CloudWatch Metrics Dashboard</h1>

      <div className={styles.controls}>
        {/* Dropdown for selecting instance */}
        <select
          value={selectedInstanceId}
          onChange={(e) => setSelectedInstanceId(e.target.value)}
          className={styles.dropdown}
        >
          <option value="" disabled>Select Instance</option>
          {instances.map((instance) => (
            <option key={instance.id} value={instance.id}>
              {instance.name}
            </option>
          ))}
        </select>

        {/* Fetch Metrics Button */}
        <button onClick={fetchMetrics} className={styles.fetchButton}>
          {loading ? 'Loading...' : 'Fetch Metrics'}
        </button>
      </div>

      {/* Render Graphs */}
      {metrics.length > 0 ? (
        <div className={styles.charts}>
          <div className={styles.chartCard}>
            <h3>CPU Utilization</h3>
            <Line data={processData('CPUUtilization')} />
          </div>
          <div className={styles.chartCard}>
            <h3>Network In</h3>
            <Line data={processData('NetworkIn')} />
            </div>
            <div className={styles.chartCard}>
            <h3>Network Out</h3>
            <Line data={processData('NetworkOut')} />
            </div>

          
        </div>
      ) : (
        <p>Select an instance and fetch metrics to view graphs.</p>
      )}
    </div>
  );
}

export default Dashboard;
