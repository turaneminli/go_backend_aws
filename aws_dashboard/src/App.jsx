import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import InstanceList from './components/InstanceList';
import LaunchInstance from './components/LaunchInstance';
import ListBuckets from './components/ListBuckets';
import CreateBucket from './components/CreateBucket';
import DeleteBucket from './components/DeleteBucket';
import ListBucketObjects from './components/ListBucketObjects';

function App() {
  return (
    <Router>
      <Routes>
        {/* Instance-related routes */}
        <Route path="/" element={<InstanceList />} />
        <Route path="/launch-instance" element={<LaunchInstance />} />

        {/* S3-related routes */}
        <Route path="/s3/buckets" element={<ListBuckets />} />
        <Route path="/s3/buckets/create" element={<CreateBucket />} />
        <Route path="/s3/buckets/delete" element={<DeleteBucket />} />
        <Route path="/s3/buckets/objects" element={<ListBucketObjects />} />
      </Routes>
    </Router>
  );
}

export default App;
