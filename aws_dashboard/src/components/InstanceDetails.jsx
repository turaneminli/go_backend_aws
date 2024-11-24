import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Line } from 'react-chartjs-2';  // Import the chart component
import 'react-toastify/dist/ReactToastify.css';
import styles from './css/InstanceDetails.module.css';
import { Chart as ChartJS, LineElement, LinearScale, CategoryScale, PointElement, Tooltip, Legend } from 'chart.js';

// Register necessary Chart.js components
ChartJS.register(LineElement, LinearScale, CategoryScale, PointElement, Tooltip, Legend);

function InstanceDetails() {
  const { instanceId } = useParams();
  const [instance, setInstance] = useState(null);
  const [metrics, setMetrics] = useState([]);  // New state for storing CloudWatch metrics
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    // Fetch instance details
    axios
      .get(`http://localhost:8080/instances/detail?instanceId=${instanceId}`)
      .then((response) => {
        setInstance(response.data);
        setLoading(false);
      })
      .catch((error) => {
        console.error('Error fetching instance details:', error);
        setLoading(false);
        toast.error('Failed to load instance details');
      });

    // Fetch CloudWatch metrics for CPU utilization
    axios
      .get(`http://localhost:8080/cloudwatch/metrics?instanceId=${instanceId}`)
      .then((response) => {
        if (!response.data) {
          throw new Error("No CloudWatch data available");
        }
        setMetrics(response.data);
      })
      .catch((error) => {
        console.error('Error fetching CloudWatch metrics:', error);
        toast.error('Failed to load CloudWatch metrics');
      });
  }, [instanceId]);

  const handleAction = (action) => {
    axios
      .post(`http://localhost:8080/instances/${action}?instanceId=${instanceId}`)
      .then(() => {
        toast.success(`Instance ${action} successfully!`);
        setInstance({ ...instance, state: action === 'start' ? 'running' : 'stopped' });
      })
      .catch((error) => {
        console.error(`Error performing action ${action}:`, error);
        toast.error(`Failed to ${action} the instance.`);
      });
  };

  const processData = (metricName) => {
    const filteredMetrics = metrics.filter((metric) => metric.metric_name === metricName);
    return {
      labels: filteredMetrics.map((metric) => new Date(metric.timestamp).toLocaleTimeString()),
      datasets: [
        {
          label: metricName,
          data: filteredMetrics.map((metric) => metric.value),
          borderColor: '#007bff',
          backgroundColor: 'rgba(0, 123, 255, 0.2)',
          fill: true,
          tension: 0.4,
        },
      ],
    };
  };

  if (loading) return <div className={styles.loading}>Loading instance details...</div>;

  return (
    <div className={styles.container}>
      <h2 className={styles.header}>Instance Details - {instance.name}</h2>

      <div className={styles.details}>
        <div className={styles.card}>
          <h3>Basic Information</h3>
          <p><strong>Instance ID:</strong> {instance.id}</p>
          <p><strong>State:</strong> {instance.state}</p>
          <p><strong>Launch Time:</strong> {instance.launchTime}</p>
          <p><strong>Type:</strong> {instance.instanceType}</p>
          <p><strong>Public IP:</strong> {instance.publicIp || 'N/A'}</p>
          <p><strong>Private IP:</strong> {instance.privateIp || 'N/A'}</p>
        </div>

        <div className={styles.card}>
          <h3>Actions</h3>
          <div className={styles.actions}>
            <button onClick={() => handleAction('stop')} disabled={instance.state !== 'running'} className={styles.actionButton}>
              Stop Instance
            </button>
            <button onClick={() => handleAction('start')} disabled={instance.state !== 'stopped'} className={styles.actionButton}>
              Start Instance
            </button>
            <button onClick={() => handleAction('reboot')} className={styles.actionButton}>
              Reboot Instance
            </button>
            <button onClick={() => handleAction('terminate')} className={styles.actionButton}>
              Terminate Instance
            </button>
          </div>
        </div>

        <div className={styles.card}>
          <h3>CloudWatch Metrics - CPU Utilization</h3>
          {metrics.length > 0 ? (
            <Line data={processData('CPUUtilization')} />
          ) : (
            <p>Loading CPU metrics...</p>
          )}
        </div>

        <div className={styles.card}>
          <h3>Security Group</h3>
          <div className={styles.listContainer}>
            {instance.securityGroups.length > 0 ? (
              instance.securityGroups.map((group, index) => (
                <div key={index} className={styles.listItem}>
                  <p>{group}</p>
                </div>
              ))
            ) : (
              <p>No security groups available.</p>
            )}
          </div>
        </div>

        <div className={styles.card}>
          <h3>Attached Volumes</h3>
          <div className={styles.listContainer}>
            {instance.volumes.length > 0 ? (
              instance.volumes.map((volume, index) => (
                <div key={index} className={styles.listItem}>
                  <p>{volume}</p>
                </div>
              ))
            ) : (
              <p>No attached volumes.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default InstanceDetails;
