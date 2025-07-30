import React from 'react';
import { Heart, HeartOff, Star, Play, Clock, Calendar } from 'lucide-react';
import axios from 'axios';

function AnimeCard({ anime, onWatchlistChange, viewMode = 'grid', darkMode = false }) {
  const handleWatchlistToggle = async () => {
    try {
      if (anime.is_watching) {
        await axios.delete(`/api/anime/${anime.anilist_id}/watch`);
      } else {
        await axios.post(`/api/anime/${anime.anilist_id}/watch`);
      }
      onWatchlistChange(anime.anilist_id, !anime.is_watching);
    } catch (error) {
      console.error('Failed to toggle watchlist:', error);
    }
  };

  const formatScore = (score) => {
    return score ? score.toFixed(1) : 'N/A';
  };

  const formatDuration = (duration) => {
    if (!duration) return 'Unknown';
    const hours = Math.floor(duration / 60);
    const minutes = duration % 60;
    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    }
    return `${minutes}m`;
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'RELEASING':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      case 'FINISHED':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case 'NOT_YET_RELEASED':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
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

  const sanitizeHtml = (html) => {
    if (!html) return '';
    
    // Convert HTML tags to proper formatting
    let processed = html
      .replace(/<i>/gi, '<em>')
      .replace(/<\/i>/gi, '</em>')
      .replace(/<b>/gi, '<strong>')
      .replace(/<\/b>/gi, '</strong>')
      .replace(/<br\s*\/?>/gi, '<br>');
    
    // Allow only specific tags and remove others
    processed = processed.replace(/<[^>]*>/g, (match) => {
      const allowedTags = ['em', 'strong', 'br', 'p', 'span'];
      const tagName = match.replace(/[<>]/g, '').split(' ')[0].toLowerCase();
      return allowedTags.includes(tagName) ? match : '';
    });
    
    return processed;
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
          <div className="flex items-start justify-between mb-2">
            <div className="flex-1 min-w-0">
              <h3 className={`font-semibold text-lg mb-1 line-clamp-1 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
                {anime.title}
              </h3>
              
              {anime.title_english && anime.title_english !== anime.title && (
                <p className={`text-sm mb-1 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                  {anime.title_english}
                </p>
              )}
            </div>
            
            <button
              onClick={handleWatchlistToggle}
              className={`p-2 rounded-full transition-colors ${
                anime.is_watching 
                  ? 'bg-red-500 text-white hover:bg-red-600' 
                  : 'bg-gray-200 text-gray-600 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600'
              }`}
            >
              {anime.is_watching ? <Heart className="h-4 w-4 fill-current" /> : <HeartOff className="h-4 w-4" />}
            </button>
          </div>

          {anime.description && (
            <div
              className={`text-sm mb-3 line-clamp-2 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}
              dangerouslySetInnerHTML={{ __html: sanitizeHtml(anime.description) }}
            />
          )}

          {/* Details */}
          <div className={`flex flex-wrap gap-4 text-xs ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
            {anime.status && (
              <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(anime.status)}`}>
                {getStatusText(anime.status)}
              </span>
            )}
            
            {anime.episodes && (
              <div className="flex items-center space-x-1">
                <Play className="h-3 w-3" />
                <span>{anime.episodes} episodes</span>
              </div>
            )}
            
            {anime.duration && (
              <div className="flex items-center space-x-1">
                <Clock className="h-3 w-3" />
                <span>{formatDuration(anime.duration)}</span>
              </div>
            )}
            
            {anime.score && (
              <div className="flex items-center space-x-1">
                <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                <span>{formatScore(anime.score)}</span>
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
      </div>
    );
  }

  return (
    <div className={`group relative overflow-hidden rounded-lg shadow-md transition-all duration-300 hover:shadow-lg ${
      darkMode ? 'bg-gray-800' : 'bg-white'
    }`}>
      {/* Cover Image */}
      <div className="relative aspect-[3/4] overflow-hidden">
        <img
          src={anime.cover_image || '/placeholder-anime.jpg'}
          alt={anime.title}
          className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
          onError={(e) => {
            e.target.src = '/placeholder-anime.jpg';
          }}
        />
        
        {/* Overlay with Status and Score */}
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent">
          {/* Status Badge */}
          {anime.status && (
            <div className="absolute top-2 left-2">
              <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(anime.status)}`}>
                {getStatusText(anime.status)}
              </span>
            </div>
          )}
          
          {/* Score */}
          {anime.score && (
            <div className="absolute top-2 right-2 flex items-center space-x-1 bg-black/70 text-white px-2 py-1 rounded-full">
              <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
              <span className="text-xs font-medium">{formatScore(anime.score)}</span>
            </div>
          )}
          
          {/* Watchlist Button */}
          <div className="absolute bottom-2 right-2">
            <button
              onClick={handleWatchlistToggle}
              className={`p-2 rounded-full transition-colors ${
                anime.is_watching 
                  ? 'bg-red-500 text-white hover:bg-red-600' 
                  : 'bg-black/50 text-white hover:bg-black/70'
              }`}
            >
              {anime.is_watching ? <Heart className="h-4 w-4 fill-current" /> : <HeartOff className="h-4 w-4" />}
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="p-4">
        <h3 className={`font-semibold text-lg mb-1 line-clamp-2 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
          {anime.title}
        </h3>
        
        {anime.title_english && anime.title_english !== anime.title && (
          <p className={`text-sm mb-2 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
            {anime.title_english}
          </p>
        )}

        {anime.description && (
          <div 
            className={`text-sm mb-3 line-clamp-3 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}
            dangerouslySetInnerHTML={{ __html: sanitizeHtml(anime.description) }}
          />
        )}

        {/* Details */}
        <div className={`flex flex-wrap gap-2 text-xs ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
          {anime.episodes && (
            <div className="flex items-center space-x-1">
              <Play className="h-3 w-3" />
              <span>{anime.episodes} episodes</span>
            </div>
          )}
          
          {anime.duration && (
            <div className="flex items-center space-x-1">
              <Clock className="h-3 w-3" />
              <span>{formatDuration(anime.duration)}</span>
            </div>
          )}
        </div>

        {/* Genres */}
        {anime.genres && (
          <div className="mt-2 flex flex-wrap gap-1">
            {anime.genres.split(', ').slice(0, 2).map((genre, index) => (
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