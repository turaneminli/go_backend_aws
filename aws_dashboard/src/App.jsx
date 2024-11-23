import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'; // Import Router and Routes for navigation
import LaunchInstance from './components/LaunchInstance'; // Import LaunchInstance
import { ToastContainer } from 'react-toastify';  // Import ToastContainer
import 'react-toastify/dist/ReactToastify.css';  // Import Toastify styles
import InstanceList from './components/InstanceList';
import Navbar from './components/Navbar';

function App() {
  return (
    <Router>
      <div className="App">
      <Navbar/> 
        <Routes>
          <Route path="/launch-instance" element={<LaunchInstance />} />
          <Route path="/" element={<InstanceList />} />

          {/* You can add more routes for other pages */}
        </Routes>

        {/* ToastContainer component to display the toasts */}
        <ToastContainer />
      </div>
    </Router>
  );
}

export default App;
