import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import styles from './InstanceList.module.css';

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
      });
  }, []);

  // Stop an instance by ID
  const stopInstance = (id) => {
    axios
      .post('http://localhost:8080/instances/stop', { id }) // POST request to stop the instance
      .then(() => {
        setInstances(instances.filter((instance) => instance.id !== id)); // Remove the instance from the list
        alert('Instance stopped successfully!');
      })
      .catch((error) => {
        console.error('Error stopping instance:', error);
        alert('Failed to stop the instance.');
      });
  };

  // Launch a new instance (redirect to the launch page)
  const handleLaunchInstance = () => {
    navigate('/launch-instance'); // Assuming '/launch-instance' is the route for launching an instance
  };

  if (loading) return <p>Loading instances...</p>;

  return (
    <div className={styles.container}>
      <h2 className={styles.header}>EC2 Instances</h2>
      
      {/* Button to launch a new instance */}
      <button className={styles.launchButton} onClick={handleLaunchInstance}>Launch New Instance</button>

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
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default InstanceList;
