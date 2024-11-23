import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';  // Import useNavigate from React Router
import { toast } from 'react-toastify';  // Import toast from react-toastify
import 'react-toastify/dist/ReactToastify.css';  // Import the CSS for Toastify
import './css/LaunchInstance.css';  // Import the CSS file for styling

function LaunchInstance() {
  const [formData, setFormData] = useState({
    region: '',
    instanceType: '',
    instanceName: '',
    ami: '',
    keyPair: '',
    minCount: 1,
    maxCount: 1,
    securityGroups: [],
  });

  const [regions, setRegions] = useState([]);
  const [securityGroups, setSecurityGroups] = useState([]);
  const [loadingRegions, setLoadingRegions] = useState(true);
  const [loadingSecurityGroups, setLoadingSecurityGroups] = useState(true);

  const navigate = useNavigate(); // Initialize the navigate function

  // Fetch available regions from the backend
  useEffect(() => {
    axios
      .get('http://localhost:8080/regions')
      .then((response) => {
        const regionNames = response.data.regions.map(region => region.RegionName);
        setRegions(regionNames);
        setLoadingRegions(false);
      })
      .catch((error) => {
        console.error('Error fetching regions:', error);
        setLoadingRegions(false);
      });
  }, []);

  // Fetch available security groups from the backend
  useEffect(() => {
    axios
      .get('http://localhost:8080/security-groups')
      .then((response) => {
        setSecurityGroups(response.data);
        setLoadingSecurityGroups(false);
      })
      .catch((error) => {
        console.error('Error fetching security groups:', error);
        setLoadingSecurityGroups(false);
      });
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;

    if (name === 'securityGroups') {
      const selectedGroups = Array.from(e.target.selectedOptions, option => option.value);
      setFormData({ ...formData, [name]: selectedGroups });
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    axios
      .post('http://localhost:8080/instances/launch', formData)
      .then((response) => {
        toast.success('Instance launched successfully!'); // Success notification
      })
      .catch((error) => {
        toast.error(`Failed to launch instance. ${error.response.data}`); // Error notification
      });
  };

  return (
    <div className="launch-instance-container">
      {/* Back to Instance List Button */}
      <button 
        className="back-button" 
        onClick={() => navigate('/')}  // Navigate back to the instance list page
      >
        Back to Instance List
      </button>

      <h2>Launch EC2 Instance</h2>
      
      {/* ToastContainer component */}
      <form onSubmit={handleSubmit} className="launch-form">
        <div className="form-group">
          <label htmlFor="region">Region:</label>
          {loadingRegions ? (
            <p>Loading regions...</p>
          ) : (
            <select
              name="region"
              value={formData.region}
              onChange={handleChange}
              required
              className="form-input"
            >
              <option value="" disabled>Select Region</option>
              {regions.map((region) => (
                <option key={region} value={region}>
                  {region}
                </option>
              ))}
            </select>
          )}
        </div>

        <div className="form-group">
          <label htmlFor="instanceName">Instance Name:</label>
          <input
            type="text"
            id="instanceName"
            name="instanceName"
            value={formData.instanceName}
            onChange={handleChange}
            required
            className="form-input"
          />
        </div>

        <div className="form-group">
          <label htmlFor="instanceType">Instance Type:</label>
          <input
            type="text"
            id="instanceType"
            name="instanceType"
            value={formData.instanceType}
            onChange={handleChange}
            required
            className="form-input"
          />
        </div>

        <div className="form-group">
          <label htmlFor="ami">AMI ID:</label>
          <input
            type="text"
            id="ami"
            name="ami"
            value={formData.ami}
            onChange={handleChange}
            required
            className="form-input"
          />
        </div>

        <div className="form-group">
          <label htmlFor="keyPair">Key Pair:</label>
          <input
            type="text"
            id="keyPair"
            name="keyPair"
            value={formData.keyPair}
            onChange={handleChange}
            required
            className="form-input"
          />
        </div>

        <div className="form-group">
          <label htmlFor="minCount">Min Count:</label>
          <input
            type="number"
            id="minCount"
            name="minCount"
            value={formData.minCount}
            onChange={handleChange}
            className="form-input"
          />
        </div>

        <div className="form-group">
          <label htmlFor="maxCount">Max Count:</label>
          <input
            type="number"
            id="maxCount"
            name="maxCount"
            value={formData.maxCount}
            onChange={handleChange}
            className="form-input"
          />
        </div>

        <div className="form-group">
          <label htmlFor="securityGroups">Security Groups:</label>
          {loadingSecurityGroups ? (
            <p>Loading security groups...</p>
          ) : (
            <select
              name="securityGroups"
              value={formData.securityGroups}
              onChange={handleChange}
              required
              className="form-input"
            >
              {securityGroups.map((group) => (
                <option key={group.GroupId} value={group.GroupId}>
                  {group.GroupName}
                </option>
              ))}
            </select>
          )}
        </div>

        <button type="submit" className="launch-button">Launch Instance</button>
      </form>
    </div>
  );
}

export default LaunchInstance;
