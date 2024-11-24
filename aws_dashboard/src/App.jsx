import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'; // Import Router and Routes for navigation
import LaunchInstance from './components/LaunchInstance'; // Import LaunchInstance
import { ToastContainer } from 'react-toastify';  // Import ToastContainer
import 'react-toastify/dist/ReactToastify.css';  // Import Toastify styles
import InstanceList from './components/InstanceList';
import Navbar from './components/Navbar';
import Dashboard from './components/Dashboard';
import ListBuckets from './components/ListBuckets';


function App() {
  return (
    <Router>
      <div className="App">
        <Navbar />
        <Routes>
          <Route path="/launch-instance" element={<LaunchInstance />} />
          <Route path="/" element={<InstanceList />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/s3/buckets" element={<ListBuckets />} />
        </Routes>
        <ToastContainer />
      </div>
    </Router>
  );
}

export default App;
