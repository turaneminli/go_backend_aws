import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import styles from './css/Navbar.module.css';

function Navbar() {
  const location = useLocation();

  return (
    <nav className={styles.navbar}>
      <div className={styles.logo}>
        <Link to="/">Turandrei AWS Management</Link> {/* Replace with your app name */}
      </div>
      <ul className={styles.navLinks}>
        <li className={location.pathname === '/' ? styles.active : ''}>
          <Link to="/">Instance List</Link>
        </li>
        <li className={location.pathname === '/launch-instance' ? styles.active : ''}>
          <Link to="/launch-instance">Launch Instance</Link>
        </li>
      </ul>
    </nav>
  );
}

export default Navbar;
