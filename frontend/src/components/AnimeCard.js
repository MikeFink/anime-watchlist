import React, { useState } from 'react';
import { Heart, HeartOff, Star, Play, Calendar, Clock } from 'lucide-react';
import axios from 'axios';

function AnimeCard({ anime, onWatchlistChange, viewMode = 'grid', darkMode = false }) {
  const [loading, setLoading] = useState(false);

  const handleWatchlistToggle = async () => {
    setLoading(true);
    try {
      if (anime.is_watching) {
        await axios.delete(`/api/anime/${anime.id}/watch`);
      } else {
        await axios.post(`/api/anime/${anime.id}/watch`);
      }
      onWatchlistChange(anime.id, !anime.is_watching);
    } catch (error) {
      console.error('Failed to update watchlist:', error);
      alert('Failed to update watchlist. Please try again.');
    } finally {
      setLoading(false);
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

  if (viewMode === 'list') {
    return (
      <div className={`flex space-x-4 p-4 rounded-lg shadow-md transition-all duration-300 hover:shadow-lg ${
        darkMode ? 'bg-gray-800' : 'bg-white'
      }`}>
        {/* Cover Image */}
        <div className="relative w-24 h-32 flex-shrink-0">
          <img
            src={anime.cover_image || '/placeholder-anime.jpg'}
            alt={anime.title}
            className="w-full h-full object-cover rounded"
            onError={(e) => {
              e.target.src = '/placeholder-anime.jpg';
            }}
          />
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between">
            <div className="flex-1 min-w-0">
              <h3 className={`font-semibold text-lg mb-1 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
                {anime.title}
              </h3>
              
              {anime.title_english && anime.title_english !== anime.title && (
                <p className={`text-sm mb-2 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                  {anime.title_english}
                </p>
              )}

              {anime.description && (
                <p className={`text-sm mb-3 line-clamp-2 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                  {anime.description}
                </p>
              )}

              {/* Details */}
              <div className={`flex flex-wrap gap-4 text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                {anime.episodes && (
                  <div className="flex items-center space-x-1">
                    <Play className="h-4 w-4" />
                    <span>{anime.episodes} episodes</span>
                  </div>
                )}
                
                {anime.duration && (
                  <div className="flex items-center space-x-1">
                    <Clock className="h-4 w-4" />
                    <span>{formatDuration(anime.duration)}</span>
                  </div>
                )}
                
                {anime.season && anime.season_year && (
                  <div className="flex items-center space-x-1">
                    <Calendar className="h-4 w-4" />
                    <span>{anime.season} {anime.season_year}</span>
                  </div>
                )}
              </div>

              {/* Genres */}
              {anime.genres && (
                <div className="mt-2 flex flex-wrap gap-1">
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

            {/* Right side - Score and Watchlist */}
            <div className="flex flex-col items-end space-y-2 ml-4">
              {anime.score && (
                <div className="flex items-center space-x-1 bg-black/70 text-white px-2 py-1 rounded-full">
                  <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                  <span className="text-xs font-medium">{formatScore(anime.score)}</span>
                </div>
              )}

              {anime.status && (
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(anime.status)}`}>
                  {getStatusText(anime.status)}
                </span>
              )}

              <button
                onClick={handleWatchlistToggle}
                disabled={loading}
                className={`p-2 rounded-full shadow-lg transition-all duration-200 ${
                  anime.is_watching
                    ? 'bg-red-500 hover:bg-red-600 text-white'
                    : 'bg-white/90 hover:bg-white text-gray-600'
                }`}
              >
                {loading ? (
                  <div className="loading-spinner h-4 w-4" />
                ) : anime.is_watching ? (
                  <Heart className="h-4 w-4 fill-current" />
                ) : (
                  <HeartOff className="h-4 w-4" />
                )}
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={`group relative overflow-hidden rounded-lg shadow-lg transition-all duration-300 hover:shadow-xl ${
      darkMode ? 'bg-gray-800' : 'bg-white'
    }`}>
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

        {/* Watchlist Button */}
        <button
          onClick={handleWatchlistToggle}
          disabled={loading}
          className="absolute bottom-2 right-2 p-2 bg-white/90 hover:bg-white rounded-full shadow-lg transition-all duration-200"
        >
          {loading ? (
            <div className="loading-spinner h-4 w-4" />
          ) : anime.is_watching ? (
            <Heart className="h-4 w-4 fill-red-500 text-red-500" />
          ) : (
            <HeartOff className="h-4 w-4 text-gray-600" />
          )}
        </button>
      </div>

      {/* Content */}
      <div className={`p-4 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
        <h3 className="font-semibold text-lg mb-2 line-clamp-2">{anime.title}</h3>
        
        {anime.title_english && anime.title_english !== anime.title && (
          <p className={`text-sm mb-2 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
            {anime.title_english}
          </p>
        )}

        {anime.description && (
          <p className={`text-sm mb-3 line-clamp-3 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
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
  );
}

export default AnimeCard; 