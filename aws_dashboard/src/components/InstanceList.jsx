import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import ConfirmationModal from './ConfirmationModal';
import 'react-toastify/dist/ReactToastify.css';
import styles from './css/InstanceList.module.css';

function InstanceList() {
  const [instances, setInstances] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modal, setModal] = useState({ isOpen: false, action: null, instanceId: null });
  const navigate = useNavigate();

  useEffect(() => {
    axios
      .get('http://localhost:8080/instances/status')
      .then((response) => {
        setInstances(response.data);
        setLoading(false);
      })
      .catch((error) => {
        console.error('Error fetching instance statuses:', error);
        setLoading(false);
        toast.error('Failed to load instances');
      });
  }, []);

  const handleAction = () => {
    const { action, instanceId } = modal;

    if (action === 'stop') {
      axios
        .post(`http://localhost:8080/instances/stop?instanceId=${instanceId}`)
        .then(() => {
          setInstances(instances.filter((instance) => instance.id !== instanceId));
          toast.success('Instance stopped successfully!');
        })
        .catch((error) => {
          console.error('Error stopping instance:', error);
          toast.error('Failed to stop the instance.');
        });
    } else if (action === 'start') {
      axios
        .post(`http://localhost:8080/instances/start?instanceId=${instanceId}`)
        .then(() => {
          setInstances(
            instances.map((instance) =>
              instance.id === instanceId ? { ...instance, state: 'running' } : instance
            )
          );
          toast.success('Instance started successfully!');
        })
        .catch((error) => {
          console.error('Error starting instance:', error);
          toast.error('Failed to start the instance.');
        });
    }

    setModal({ isOpen: false, action: null, instanceId: null });
  };

  const openModal = (action, id) => {
    setModal({ isOpen: true, action, instanceId: id });
  };

  const closeModal = () => {
    setModal({ isOpen: false, action: null, instanceId: null });
  };

  const handleInstanceClick = (instanceId) => {
    navigate(`/instance-details/${instanceId}`);
  };

  if (loading) return <p>Loading instances...</p>;

  return (
    <div className={styles.container}>
      <h2 className={styles.header}>EC2 Instances</h2>
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
              <td
                onClick={() => handleInstanceClick(instance.id)}
                className={styles.instanceLink}
              >
                {instance.name}
              </td>
              <td>{instance.id}</td>
              <td>{instance.state}</td>
              <td>{instance.public_ip}</td>
              <td>{instance.private_ip}</td>
              <td>
                <button
                  onClick={() => openModal('stop', instance.id)}
                  className={styles.stopButton}
                  disabled={instance.state !== 'running'}
                >
                  Stop Instance
                </button>
                <button
                  onClick={() => openModal('start', instance.id)}
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

      <ConfirmationModal
        isOpen={modal.isOpen}
        message={`Are you sure you want to ${modal.action} this instance?`}
        onConfirm={handleAction}
        onCancel={closeModal}
      />
    </div>
  );
}

export default InstanceList;
