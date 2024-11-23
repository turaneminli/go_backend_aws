import React, { useEffect, useState } from 'react';
import axios from 'axios';
import styles from './ListBuckets.module.css'; // Import the CSS file

const ListBuckets = () => {
  const [buckets, setBuckets] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    // Fetch list of buckets
    axios.get('http://localhost:8080/s3/buckets')
      .then(response => {
        setBuckets(response.data);
      })
      .catch(err => {
        setError('Error fetching buckets');
      });
  }, []);

  return (
    <div className={styles.container}>
      <h2 className={styles.header}>List of S3 Buckets</h2>

      {error && <p className={styles.errorMessage}>{error}</p>}

      <button className={styles.createBucketButton}>Create New Bucket</button>

      {buckets.length > 0 ? (
        <table className={styles.instanceTable}>
          <thead>
            <tr>
              <th>Bucket Name</th>
              <th>Creation Date</th>
            </tr>
          </thead>
          <tbody>
            {buckets.map((bucket, index) => (
              <tr key={index}>
                <td>{bucket.name}</td>
                <td>{bucket.creation_date}</td>
                <td>
                  <button className={styles.stopButton}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <p>No buckets available.</p>
      )}
    </div>
  );
};

export default ListBuckets;
