import React from 'react';
import styles from './Navbar.module.css';

function Navbar() {
  return (
    <nav className={styles.navbar}>
      <h1 className={styles.title}>EC2 Metrics Dashboard</h1>
    </nav>
  );
}

export default Navbar;
