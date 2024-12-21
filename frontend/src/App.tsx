import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import Navbar from './components/layout/Navbar';
import Portfolio from './pages/Portfolio';
import CoinDetail from './pages/CoinDetail';
import Profile from './pages/Profile';

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-black">
        <Toaster
          position="top-center"
          toastOptions={{
            style: {
              background: '#1C1C1E',
              color: '#fff',
              borderRadius: '12px',
            },
          }}
        />
        <main className="pb-16">
          <Routes>
            <Route path="/" element={<Portfolio />} />
            <Route path="/coins/:symbol" element={<CoinDetail />} />
            <Route path="/profile" element={<Profile />} />
          </Routes>
        </main>
        <Navbar />
      </div>
    </Router>
  );
}

export default App;
