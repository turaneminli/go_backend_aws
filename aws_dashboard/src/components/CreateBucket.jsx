import React, { useState } from 'react';
import axios from 'axios';

const CreateBucket = () => {
  const [bucketName, setBucketName] = useState('');
  const [region, setRegion] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    const payload = {
      bucket_name: bucketName,
      region: region,
    };

    axios.post('http://localhost:8080/s3/buckets/create', payload)
      .then(response => {
        setMessage(response.data.message);
        setError('');
      })
      .catch(err => {
        setError('Error creating bucket');
        setMessage('');
      });
  };

  return (
    <div>
      <h2>Create S3 Bucket</h2>
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
        <div>
          <label>Region:</label>
          <input 
            type="text" 
            value={region} 
            onChange={(e) => setRegion(e.target.value)} 
            required 
          />
        </div>
        <button type="submit">Create Bucket</button>
      </form>
      {message && <p>{message}</p>}
      {error && <p>{error}</p>}
    </div>
  );
};

export default CreateBucket;
