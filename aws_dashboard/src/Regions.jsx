import React, { useState, useEffect } from 'react';
import axios from 'axios';

function Regions() {
  const [regions, setRegions] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    axios
      .get('http://localhost:8080/regions') // Replace with your Go API base URL
      .then((response) => {
        setRegions(response.data); // Assuming the API returns an array of regions
        setLoading(false);
      })
      .catch((error) => {
        console.error('Error fetching regions:', error);
        setLoading(false);
      });
  }, []);

  if (loading) return <p>Loading regions...</p>;

  return (
    <div>
      <h2>Available Regions</h2>
      <ul>
        {regions.map((region, index) => (
          <li key={index}>{region}</li>
        ))}
      </ul>
    </div>
  );
}

export default Regions;
