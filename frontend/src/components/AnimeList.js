import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import { Search, Filter, Grid, List } from 'lucide-react';
import AnimeCard from './AnimeCard';
import LoadingSpinner from './LoadingSpinner';

function AnimeList({ darkMode = false }) {
  const [anime, setAnime] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [viewMode, setViewMode] = useState('grid');
  const searchInputRef = useRef(null);

  useEffect(() => {
    fetchAnime();
  }, []);

  const fetchAnime = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/anime/search');
      setAnime(response.data);
    } catch (error) {
      console.error('Failed to fetch anime:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleWatchlistChange = (animeId, isWatching) => {
    setAnime(prevAnime =>
      prevAnime.map(item =>
        item.anilist_id === animeId ? { ...item, is_watching: isWatching } : item
      )
    );
  };

  const handleSearch = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams();
      if (searchTerm) params.append('search', searchTerm);
      if (statusFilter !== 'all') params.append('status', statusFilter);
      
      const response = await axios.get(`/api/anime/search?${params.toString()}`);
      setAnime(response.data);
    } catch (error) {
      console.error('Failed to search anime:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (searchTerm || statusFilter !== 'all') {
        handleSearch();
      } else {
        fetchAnime();
      }
    }, 500);

    return () => clearTimeout(timeoutId);
  }, [searchTerm, statusFilter]);

  const handleSearchInputChange = (e) => {
    setSearchTerm(e.target.value);
  };

  const handleStatusFilterChange = (e) => {
    setStatusFilter(e.target.value);
  };

  return (
    <div className="space-y-6">
      {/* Search and Filters - Fixed Position */}
      <div className="sticky top-0 z-10 bg-white dark:bg-gray-900 py-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex flex-col sm:flex-row gap-4">
          {/* Search */}
          <div className="flex-1 relative">
            <Search className={`absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 ${
              darkMode ? 'text-gray-400' : 'text-gray-400'
            }`} />
            <input
              ref={searchInputRef}
              type="text"
              placeholder="Search anime..."
              value={searchTerm}
              onChange={handleSearchInputChange}
              className={`w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors ${
                darkMode 
                  ? 'bg-gray-800 border-gray-600 text-white placeholder-gray-400' 
                  : 'bg-white border-gray-300 text-gray-900 placeholder-gray-500'
              }`}
            />
          </div>

          {/* Status Filter */}
          <div className="relative">
            <Filter className={`absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 ${
              darkMode ? 'text-gray-400' : 'text-gray-400'
            }`} />
            <select
              value={statusFilter}
              onChange={handleStatusFilterChange}
              className={`pl-10 pr-8 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors ${
                darkMode 
                  ? 'bg-gray-800 border-gray-600 text-white' 
                  : 'bg-white border-gray-300 text-gray-900'
              }`}
            >
              <option value="all">All Status</option>
              <option value="RELEASING">Airing</option>
              <option value="FINISHED">Completed</option>
              <option value="NOT_YET_RELEASED">Upcoming</option>
            </select>
          </div>
        </div>
      </div>

      {/* Results Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center space-y-4 sm:space-y-0">
        <div>
          <h2 className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-900'}`}>
            Search Results
          </h2>
          <p className={`mt-1 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
            {anime.length} anime found
          </p>
        </div>
        
        {/* View Mode Toggle */}
        <div className="flex items-center space-x-2">
          <button
            onClick={() => setViewMode('grid')}
            className={`p-2 rounded-lg transition-colors ${
              viewMode === 'grid' 
                ? 'bg-blue-100 text-blue-600 dark:bg-blue-900 dark:text-blue-400' 
                : 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
            }`}
          >
            <Grid className="h-4 w-4" />
          </button>
          <button
            onClick={() => setViewMode('list')}
            className={`p-2 rounded-lg transition-colors ${
              viewMode === 'list' 
                ? 'bg-blue-100 text-blue-600 dark:bg-blue-900 dark:text-blue-400' 
                : 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
            }`}
          >
            <List className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* Loading State */}
      {loading && (
        <div className="flex justify-center items-center h-64">
          <LoadingSpinner size="lg" />
        </div>
      )}

      {/* Anime Grid/List */}
      {!loading && anime.length === 0 ? (
        <div className={`text-center py-12 ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>
          <p className="text-lg">No anime found matching your criteria.</p>
        </div>
      ) : !loading && (
        <div className={viewMode === 'grid' ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6' : 'space-y-4'}>
          {anime.map(anime => (
            <AnimeCard 
              key={anime.anilist_id} 
              anime={anime} 
              onWatchlistChange={handleWatchlistChange}
              viewMode={viewMode}
              darkMode={darkMode}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export default AnimeList; 