import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import axios from 'axios';
import { Moon, Sun } from 'lucide-react';
import AnimeList from './components/AnimeList';
import Watchlist from './components/Watchlist';
import LoadingSpinner from './components/LoadingSpinner';

axios.defaults.baseURL = window.location.origin;

function App() {
  const [loading, setLoading] = useState(false);
  const [syncLoading, setSyncLoading] = useState(false);
  const [darkMode, setDarkMode] = useState(() => {
    const saved = localStorage.getItem('darkMode');
    return saved ? JSON.parse(saved) : false;
  });

  useEffect(() => {
    localStorage.setItem('darkMode', JSON.stringify(darkMode));
    if (darkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [darkMode]);

  const handleSync = async () => {
    setSyncLoading(true);
    try {
      await axios.post('/api/sync');
      window.location.reload();
    } catch (error) {
      console.error('Sync failed:', error);
      alert('Failed to sync anime data');
    } finally {
      setSyncLoading(false);
    }
  };

  const NavLink = ({ to, children, className = "" }) => (
    <a
      href={to}
      className={`px-4 py-2 rounded-lg transition-colors ${
        window.location.pathname === to
          ? 'bg-blue-600 text-white shadow-md'
          : 'text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700'
      } ${className}`}
    >
      {children}
    </a>
  );

  return (
    <Router>
      <div className={`min-h-screen transition-colors duration-200 ${
        darkMode 
          ? 'bg-gray-900 text-white' 
          : 'bg-gray-50 text-gray-900'
      }`}>
        <header className={`shadow-sm border-b transition-colors duration-200 ${
          darkMode 
            ? 'bg-gray-800 border-gray-700' 
            : 'bg-white border-gray-200'
        }`}>
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              <h1 className={`text-2xl font-bold ${
                darkMode ? 'text-white' : 'text-gray-900'
              }`}>
                Anime Watchlist
              </h1>
              <div className="flex items-center gap-4">
                <button
                  onClick={() => setDarkMode(!darkMode)}
                  className={`p-2 rounded-lg transition-colors ${
                    darkMode 
                      ? 'bg-gray-700 text-yellow-400 hover:bg-gray-600' 
                      : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                  }`}
                >
                  {darkMode ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
                </button>
                <button
                  onClick={handleSync}
                  disabled={syncLoading}
                  className="btn btn-primary flex items-center gap-2"
                >
                  {syncLoading ? (
                    <LoadingSpinner size="sm" />
                  ) : (
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                  )}
                  Sync Data
                </button>
              </div>
            </div>
          </div>
        </header>

        <nav className={`border-b transition-colors duration-200 ${
          darkMode 
            ? 'bg-gray-800 border-gray-700' 
            : 'bg-white border-gray-200'
        }`}>
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex space-x-8">
              <NavLink to="/">All Anime</NavLink>
              <NavLink to="/watchlist">My Watchlist</NavLink>
            </div>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Routes>
            <Route path="/" element={<AnimeList darkMode={darkMode} />} />
            <Route path="/watchlist" element={<Watchlist darkMode={darkMode} />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App; 