import React, { useState } from 'react';
import axios from 'axios';

function LaunchInstance() {
  const [formData, setFormData] = useState({
    region: '',
    instanceType: '',
  });

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    axios
      .post('http://localhost:8080/instances/launch', formData)
      .then((response) => {
        alert('Instance launched successfully!');
      })
      .catch((error) => {
        console.error('Error launching instance:', error);
        alert('Failed to launch instance.');
      });
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>
          Region:
          <input
            type="text"
            name="region"
            value={formData.region}
            onChange={handleChange}
            required
          />
        </label>
      </div>
      <div>
        <label>
          Instance Type:
          <input
            type="text"
            name="instanceType"
            value={formData.instanceType}
            onChange={handleChange}
            required
          />
        </label>
      </div>
      <button type="submit">Launch Instance</button>
    </form>
  );
}

export default LaunchInstance;
