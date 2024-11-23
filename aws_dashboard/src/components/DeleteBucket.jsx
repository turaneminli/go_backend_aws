import React, { useState } from 'react';
import axios from 'axios';

const DeleteBucket = () => {
  const [bucketName, setBucketName] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    axios.delete(`http://localhost:8080/s3/buckets/delete?bucketName=${bucketName}`)
      .then(response => {
        setMessage(response.data.message);
        setError('');
      })
      .catch(err => {
        setError('Error deleting bucket');
        setMessage('');
      });
  };

  return (
    <div>
      <h2>Delete S3 Bucket</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Bucket Name:</label>
          <input 
            type="text" 
            value={bucketName} 
            onChange={(e) => setBucketName(e.target.value)} 
            required 
          />
        </div>
        <button type="submit">Delete Bucket</button>
      </form>
      {message && <p>{message}</p>}
      {error && <p>{error}</p>}
    </div>
  );
};

export default DeleteBucket;
