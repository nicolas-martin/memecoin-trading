import { Link, useLocation } from 'react-router-dom';

const Navbar = () => {
  const location = useLocation();

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-[#1C1C1E] border-t border-gray-800 safe-area-bottom">
      <div className="max-w-md mx-auto px-4">
        <div className="flex justify-around py-2">
          <Link 
            to="/" 
            className={`flex flex-col items-center ${
              location.pathname === '/' ? 'text-white' : 'text-gray-500'
            }`}
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            <span className="text-xs mt-1">Home</span>
          </Link>

          <Link 
            to="/market" 
            className={`flex flex-col items-center ${
              location.pathname === '/market' ? 'text-white' : 'text-gray-500'
            }`}
          >
            <svg className="w-7 h-7" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"/>
              <path d="M15.5 9.5l-3.5 3.5-3.5-3.5" stroke="black" strokeWidth="1.5"/>
            </svg>
            <span className="text-xs mt-1">Market</span>
          </Link>

          <Link 
            to="/profile" 
            className={`flex flex-col items-center ${
              location.pathname === '/profile' ? 'text-white' : 'text-gray-500'
            }`}
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            <span className="text-xs mt-1">Profile</span>
          </Link>
        </div>
      </div>
    </nav>
  );
};

export default Navbar; 