import React, { useState } from 'react';
import axios from 'axios';

const ListBucketObjects = () => {
  const [bucketName, setBucketName] = useState('');
  const [objects, setObjects] = useState([]);
  const [error, setError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    axios.get(`http://localhost:8080/s3/buckets/objects?bucketName=${bucketName}`)
      .then(response => {
        setObjects(response.data);
        setError('');
      })
      .catch(err => {
        setError('Error fetching objects');
        setObjects([]);
      });
  };

  return (
    <div>
      <h2>List Objects in S3 Bucket</h2>
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
        <button type="submit">List Objects</button>
      </form>
      {error && <p>{error}</p>}
      <ul>
        {objects.length > 0 ? (
          objects.map((object, index) => (
            <li key={index}>{object}</li>
          ))
        ) : (
          <p>No objects found.</p>
        )}
      </ul>
    </div>
  );
};

export default ListBucketObjects;
