import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import InstanceList from './components/InstanceList';
import LaunchInstance from './components/LaunchInstance';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<InstanceList />} />
        <Route path="/launch-instance" element={<LaunchInstance />} />
      </Routes>
    </Router>
  );
}

export default App;
