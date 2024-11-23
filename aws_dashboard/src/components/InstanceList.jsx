import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';  // Import toast from react-toastify
import 'react-toastify/dist/ReactToastify.css';  // Import Toastify styles
import styles from './css/InstanceList.module.css';

function InstanceList() {
  const [instances, setInstances] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  // Fetch EC2 instance statuses
  useEffect(() => {
    axios
      .get('http://localhost:8080/instances/status') // Replace with your API base URL
      .then((response) => {
        setInstances(response.data); // Assuming the response is an array of instance objects
        setLoading(false);
      })
      .catch((error) => {
        console.error('Error fetching instance statuses:', error);
        setLoading(false);
        toast.error('Failed to load instances');  // Show error toast
      });
  }, []);

  // Stop an instance by ID
  const stopInstance = (id) => {
    axios
      .post(`http://localhost:8080/instances/stop?instanceId=${id}`) // Send id as query parameter in the URL
      .then(() => {
        setInstances(instances.filter((instance) => instance.id !== id)); // Remove the instance from the list
        toast.success('Instance stopped successfully!');  // Success toast
      })
      .catch((error) => {
        console.error('Error stopping instance:', error);
        toast.error('Failed to stop the instance.');  // Error toast
      });
  };

  // Start an instance by ID
  const startInstance = (id) => {
    axios
      .post(`http://localhost:8080/instances/start?instanceId=${id}`) // Send id as query parameter in the URL
      .then(() => {
        setInstances(instances.map((instance) => 
          instance.id === id ? { ...instance, state: 'running' } : instance
        ));
        toast.success('Instance started successfully!');  // Success toast
      })
      .catch((error) => {
        console.error('Error starting instance:', error);
        toast.error('Failed to start the instance.');  // Error toast
      });
  };



  if (loading) return <p>Loading instances...</p>;

  return (
    <div className={styles.container}>
      <h2 className={styles.header}>EC2 Instances</h2>
      

      {/* Table displaying instances */}
      <table className={styles.instanceTable}>
        <thead>
          <tr>
            <th>Name</th>
            <th>ID</th>
            <th>State</th>
            <th>Public IP</th>
            <th>Private IP</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {instances.map((instance) => (
            <tr key={instance.id}>
              <td>{instance.name}</td>
              <td>{instance.id}</td>
              <td>{instance.state}</td>
              <td>{instance.public_ip}</td>
              <td>{instance.private_ip}</td>
              <td>
                {/* Button to stop the instance */}
                <button 
                  onClick={() => stopInstance(instance.id)} 
                  className={styles.stopButton}
                  disabled={instance.state !== 'running'}
                >
                  Stop Instance
                </button>

                {/* Button to start the instance */}
                <button 
                  onClick={() => startInstance(instance.id)} 
                  className={styles.startButton}
                  disabled={instance.state !== 'stopped'}
                >
                  Start Instance
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default InstanceList;
