import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Heart, Trash2, Star, Play, Calendar, Clock } from 'lucide-react';
import LoadingSpinner from './LoadingSpinner';

function Watchlist({ darkMode = false }) {
  const [watchlist, setWatchlist] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchWatchlist();
  }, []);

  const fetchWatchlist = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/anime/watching');
      setWatchlist(response.data);
    } catch (error) {
      console.error('Failed to fetch watchlist:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveFromWatchlist = async (animeId) => {
    try {
      await axios.delete(`/api/anime/${animeId}/watch`);
      setWatchlist(prev => prev.filter(item => item.id !== animeId));
    } catch (error) {
      console.error('Failed to remove from watchlist:', error);
      alert('Failed to remove from watchlist. Please try again.');
    }
  };

  const formatScore = (score) => {
    if (!score) return 'N/A';
    return (score / 10).toFixed(1);
  };

  const formatDuration = (duration) => {
    if (!duration) return 'Unknown';
    return `${duration} min`;
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'RELEASING':
        return darkMode ? 'text-green-400 bg-green-900' : 'text-green-600 bg-green-100';
      case 'FINISHED':
        return darkMode ? 'text-blue-400 bg-blue-900' : 'text-blue-600 bg-blue-100';
      case 'NOT_YET_RELEASED':
        return darkMode ? 'text-yellow-400 bg-yellow-900' : 'text-yellow-600 bg-yellow-100';
      default:
        return darkMode ? 'text-gray-400 bg-gray-700' : 'text-gray-600 bg-gray-100';
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 'RELEASING':
        return 'Airing';
      case 'FINISHED':
        return 'Completed';
      case 'NOT_YET_RELEASED':
        return 'Upcoming';
      default:
        return status || 'Unknown';
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center space-x-3">
        <Heart className="h-8 w-8 text-red-500 fill-red-500" />
        <div>
          <h2 className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-900'}`}>
            My Watchlist
          </h2>
          <p className={`mt-1 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
            {watchlist.length} anime in your watchlist
          </p>
        </div>
      </div>

      {/* Watchlist */}
      {watchlist.length === 0 ? (
        <div className={`text-center py-12 ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>
          <div className="text-gray-400 mb-4">
            <Heart className="h-12 w-12 mx-auto" />
          </div>
          <h3 className={`text-lg font-medium mb-2 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
            Your watchlist is empty
          </h3>
          <p className={darkMode ? 'text-gray-400' : 'text-gray-600'}>
            Start adding anime to your watchlist from the All Anime page
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {watchlist.map(anime => (
            <div
              key={anime.id}
              className={`group relative overflow-hidden rounded-lg shadow-lg transition-all duration-300 hover:shadow-xl ${
                darkMode ? 'bg-gray-800' : 'bg-white'
              }`}
            >
              {/* Cover Image */}
              <div className="relative h-64 overflow-hidden">
                <img
                  src={anime.cover_image || '/placeholder-anime.jpg'}
                  alt={anime.title}
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                  onError={(e) => {
                    e.target.src = '/placeholder-anime.jpg';
                  }}
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
                
                {/* Status Badge */}
                {anime.status && (
                  <div className="absolute top-2 left-2">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(anime.status)}`}>
                      {getStatusText(anime.status)}
                    </span>
                  </div>
                )}

                {/* Score Badge */}
                {anime.score && (
                  <div className="absolute top-2 right-2 flex items-center space-x-1 bg-black/70 text-white px-2 py-1 rounded-full">
                    <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                    <span className="text-xs font-medium">{formatScore(anime.score)}</span>
                  </div>
                )}

                {/* Remove Button */}
                <button
                  onClick={() => handleRemoveFromWatchlist(anime.id)}
                  className="absolute bottom-2 right-2 p-2 bg-red-500 hover:bg-red-600 text-white rounded-full shadow-lg transition-all duration-200"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>

              {/* Content */}
              <div className={`p-4 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
                <h3 className="font-semibold text-lg mb-2 line-clamp-2">{anime.title}</h3>
                
                {anime.description && (
                  <p className={`text-sm mb-3 line-clamp-3 ${
                    darkMode ? 'text-gray-300' : 'text-gray-600'
                  }`}>
                    {anime.description}
                  </p>
                )}

                {/* Details */}
                <div className={`space-y-2 text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                  {anime.episodes && (
                    <div className="flex items-center space-x-2">
                      <Play className="h-4 w-4" />
                      <span>{anime.episodes} episodes</span>
                    </div>
                  )}
                  
                  {anime.duration && (
                    <div className="flex items-center space-x-2">
                      <Clock className="h-4 w-4" />
                      <span>{formatDuration(anime.duration)}</span>
                    </div>
                  )}
                  
                  {anime.season && anime.season_year && (
                    <div className="flex items-center space-x-2">
                      <Calendar className="h-4 w-4" />
                      <span>{anime.season} {anime.season_year}</span>
                    </div>
                  )}
                </div>

                {/* Genres */}
                {anime.genres && (
                  <div className="mt-3 flex flex-wrap gap-1">
                    {anime.genres.split(', ').slice(0, 3).map((genre, index) => (
                      <span
                        key={index}
                        className={`px-2 py-1 rounded-full text-xs font-medium ${
                          darkMode 
                            ? 'bg-gray-700 text-gray-300' 
                            : 'bg-gray-100 text-gray-600'
                        }`}
                      >
                        {genre}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default Watchlist; 