import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { toast } from 'react-toastify'; // Import toast from react-toastify
import 'react-toastify/dist/ReactToastify.css'; // Import Toastify styles
import styles from './css/InstanceList.module.css';

const ListBuckets = () => {
  const [buckets, setBuckets] = useState([]);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [sortOrder, setSortOrder] = useState({ name: 'asc', region: 'asc', date: 'asc' }); // State for sort order of different columns
  const [loading, setLoading] = useState(false); // Loading state

  // Fetch the list of buckets from the backend
  const fetchBuckets = () => {
    setLoading(true); // Start loading when the request is initiated
    axios.get('http://localhost:8080/s3/buckets')
      .then(response => {
        setBuckets(response.data);
        setError('');
        toast.success('Buckets fetched successfully!'); // Show success toast
      })
      .catch(err => {
        setError('Error fetching buckets');
        setMessage('');
        toast.error('Failed to fetch buckets!'); // Show error toast
      })
      .finally(() => {
        setLoading(false); // Stop loading when request finishes (success or failure)
      });
  };

  // Function to clean the date string (remove "+0000 UTC")
  const cleanDate = (dateStr) => {
    return dateStr.replace(' +0000 UTC', '');
  };

  // Sort buckets by given field
  const sortBucketsByField = (field) => {
    const sortedBuckets = [...buckets].sort((a, b) => {
      if (field === 'creation_date') {
        const dateA = new Date(cleanDate(a.creation_date));
        const dateB = new Date(cleanDate(b.creation_date));

        if (isNaN(dateA) || isNaN(dateB)) {
          console.error('Invalid date format:', a.creation_date, b.creation_date);
          return 0; // Return no change if the date is invalid
        }

        return sortOrder.date === 'asc' ? dateA - dateB : dateB - dateA;
      } else {
        return sortOrder[field] === 'asc' 
          ? a[field].localeCompare(b[field]) 
          : b[field].localeCompare(a[field]);
      }
    });

    setBuckets(sortedBuckets);
    setSortOrder(prevSortOrder => ({
      ...prevSortOrder,
      [field]: prevSortOrder[field] === 'asc' ? 'desc' : 'asc' // Toggle sort order
    }));
  };

  // Fetch the buckets when the component mounts
  useEffect(() => {
    fetchBuckets();
  }, []);

  return (
    <div className={styles.container}>
      <h2 className={styles.header}>List of S3 Buckets</h2>

      {error && <p className={styles.errorMessage}>{error}</p>}
      {message && <p className={styles.successMessage}>{message}</p>}

      {/* Loading Indicator */}
      {loading && <div className={styles.loading}>Loading...</div>}

      {/* Table to list the buckets */}
      {buckets.length > 0 ? (
        <table className={styles.instanceTable}>
          <thead>
            <tr>
              <th>
                Bucket Name 
                <button 
                  className={styles.sortButton}
                  onClick={() => sortBucketsByField('name')}
                >
                  {sortOrder.name === 'asc' ? 'Sort Desc' : 'Sort Asc'}
                </button>
              </th>
              <th>
                Region
                <button 
                  className={styles.sortButton}
                  onClick={() => sortBucketsByField('region')}
                >
                  {sortOrder.region === 'asc' ? 'Sort Desc' : 'Sort Asc'}
                </button>
              </th>
              <th>
                Creation Date 
                <button 
                  className={styles.sortButton}
                  onClick={() => sortBucketsByField('creation_date')}
                >
                  {sortOrder.date === 'asc' ? 'Sort Desc' : 'Sort Asc'}
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            {buckets.map((bucket, index) => (
              <tr key={index}>
                <td>{bucket.name}</td>
                <td>{bucket.region ? bucket.region : 'N/A'}</td> {/* Display the region or 'N/A' */}
                <td>{cleanDate(bucket.creation_date)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        !loading && <p>No buckets available.</p>
      )}
    </div>
  );
};

export default ListBuckets;
