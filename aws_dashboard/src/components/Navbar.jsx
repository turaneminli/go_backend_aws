import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import styles from './css/Navbar.module.css';

function Navbar() {
  const location = useLocation();

  return (
    <nav className={styles.navbar}>
      <div className={styles.logo}>
        <Link to="/">AWS Management</Link> {/* Replace with your app name */}
      </div>
      <ul className={styles.navLinks}>
        <li className={location.pathname === '/' ? styles.active : ''}>
          <Link to="/">Instance List</Link>
        </li>
        <li className={location.pathname === '/launch-instance' ? styles.active : ''}>
          <Link to="/launch-instance">Launch Instance</Link>
        </li>
        <li className={location.pathname === '/dashboard' ? styles.active : ''}>
          <Link to="/dashboard">Dashboard</Link>
        </li>
        <li className={location.pathname === '/s3/buckets' ? styles.active : ''}>
          <Link to="/s3/buckets">List S3 Buckets</Link>
        </li>
      </ul>
    </nav>
  );
}

export default Navbar;
