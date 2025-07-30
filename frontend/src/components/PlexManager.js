import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Server, RefreshCw, Search, MapPin, CheckCircle, XCircle, AlertCircle } from 'lucide-react';
import LoadingSpinner from './LoadingSpinner';

function PlexManager({ darkMode = false }) {
  const [serverStatus, setServerStatus] = useState(null);
  const [showsOnServer, setShowsOnServer] = useState([]);
  const [unmappedShows, setUnmappedShows] = useState([]);
  const [loading, setLoading] = useState(false);
  const [syncing, setSyncing] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState([]);

  useEffect(() => {
    fetchServerStatus();
    fetchShowsOnServer();
    fetchUnmappedShows();
  }, []);

  const fetchServerStatus = async () => {
    try {
      const response = await axios.get('/api/plex/status');
      setServerStatus(response.data);
    } catch (error) {
      console.error('Failed to fetch server status:', error);
    }
  };

  const fetchShowsOnServer = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/plex/shows');
      setShowsOnServer(response.data);
    } catch (error) {
      console.error('Failed to fetch shows on server:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchUnmappedShows = async () => {
    try {
      const response = await axios.get('/api/plex/unmapped');
      setUnmappedShows(response.data);
    } catch (error) {
      console.error('Failed to fetch unmapped shows:', error);
    }
  };

  const handleSyncPlex = async () => {
    try {
      setSyncing(true);
      await axios.post('/api/plex/sync');
      await fetchServerStatus();
      await fetchShowsOnServer();
      await fetchUnmappedShows();
    } catch (error) {
      console.error('Failed to sync plex:', error);
      alert('Failed to sync with Plex server. Please check your configuration.');
    } finally {
      setSyncing(false);
    }
  };

  const handleSearchServer = async () => {
    if (!searchTerm.trim()) {
      setSearchResults([]);
      return;
    }

    try {
      const response = await axios.get(`/api/plex/search?q=${encodeURIComponent(searchTerm)}`);
      setSearchResults(response.data);
    } catch (error) {
      console.error('Failed to search server:', error);
    }
  };

  const handleAutoMap = async (plexID) => {
    try {
      await axios.post('/api/plex/auto-map', { plex_id: plexID });
      await fetchUnmappedShows();
      await fetchShowsOnServer();
    } catch (error) {
      console.error('Failed to auto-map show:', error);
      alert('Failed to map show. Please try manually mapping.');
    }
  };

  const handleCheckShowOnServer = async (anilistID) => {
    try {
      const response = await axios.get(`/api/plex/check?anilist_id=${anilistID}`);
      return response.data.on_server;
    } catch (error) {
      console.error('Failed to check show on server:', error);
      return false;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center space-x-3">
        <Server className="h-8 w-8 text-blue-500" />
        <div>
          <h2 className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-900'}`}>
            Plex Server Manager
          </h2>
          <p className={`mt-1 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
            Manage your Plex server shows and Anilist mappings
          </p>
        </div>
      </div>

      {/* Server Status */}
      {serverStatus && (
        <div className={`grid grid-cols-1 md:grid-cols-4 gap-4 p-4 rounded-lg ${
          darkMode ? 'bg-gray-800' : 'bg-white'
        } shadow-sm`}>
          <div className="text-center">
            <div className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-900'}`}>
              {serverStatus.shows_on_server}
            </div>
            <div className={`text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
              Shows on Server
            </div>
          </div>
          <div className="text-center">
            <div className={`text-2xl font-bold text-green-600`}>
              {serverStatus.mapped_to_anilist}
            </div>
            <div className={`text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
              Mapped to Anilist
            </div>
          </div>
          <div className="text-center">
            <div className={`text-2xl font-bold text-yellow-600`}>
              {serverStatus.unmapped_shows}
            </div>
            <div className={`text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
              Unmapped Shows
            </div>
          </div>
          <div className="text-center">
            <div className={`text-2xl font-bold text-blue-600`}>
              {serverStatus.watchlist_shows}
            </div>
            <div className={`text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
              Watchlist Shows
            </div>
          </div>
        </div>
      )}

      {/* Sync Button */}
      <div className="flex justify-between items-center">
        <button
          onClick={handleSyncPlex}
          disabled={syncing}
          className={`flex items-center space-x-2 px-4 py-2 rounded-lg transition-colors ${
            syncing
              ? 'bg-gray-400 cursor-not-allowed'
              : 'bg-blue-600 hover:bg-blue-700 text-white'
          }`}
        >
          <RefreshCw className={`h-4 w-4 ${syncing ? 'animate-spin' : ''}`} />
          <span>{syncing ? 'Syncing...' : 'Sync with Plex'}</span>
        </button>
      </div>

      {/* Search Server */}
      <div className={`p-4 rounded-lg ${darkMode ? 'bg-gray-800' : 'bg-white'} shadow-sm`}>
        <h3 className={`text-lg font-semibold mb-4 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
          Search Shows on Server
        </h3>
        <div className="flex space-x-2">
          <input
            type="text"
            placeholder="Search shows on server..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className={`flex-1 px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors ${
              darkMode 
                ? 'bg-gray-700 border-gray-600 text-white placeholder-gray-400' 
                : 'bg-white border-gray-300 text-gray-900 placeholder-gray-500'
            }`}
          />
          <button
            onClick={handleSearchServer}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            <Search className="h-4 w-4" />
          </button>
        </div>

        {/* Search Results */}
        {searchResults.length > 0 && (
          <div className="mt-4 space-y-2">
            {searchResults.map(show => (
              <div
                key={show.plex_id}
                className={`flex items-center justify-between p-3 rounded-lg ${
                  darkMode ? 'bg-gray-700' : 'bg-gray-50'
                }`}
              >
                <div>
                  <div className={`font-medium ${darkMode ? 'text-white' : 'text-gray-900'}`}>
                    {show.title}
                  </div>
                  <div className={`text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                    {show.year} • {show.episode_count} episodes
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  {show.anilist_id ? (
                    <CheckCircle className="h-5 w-5 text-green-500" />
                  ) : (
                    <XCircle className="h-5 w-5 text-red-500" />
                  )}
                  {!show.anilist_id && (
                    <button
                      onClick={() => handleAutoMap(show.plex_id)}
                      className="px-3 py-1 bg-yellow-600 hover:bg-yellow-700 text-white rounded text-sm transition-colors"
                    >
                      Map
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Unmapped Shows */}
      {unmappedShows.length > 0 && (
        <div className={`p-4 rounded-lg ${darkMode ? 'bg-gray-800' : 'bg-white'} shadow-sm`}>
          <h3 className={`text-lg font-semibold mb-4 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
            Unmapped Shows ({unmappedShows.length})
          </h3>
          <div className="space-y-2">
            {unmappedShows.slice(0, 10).map(show => (
              <div
                key={show.plex_id}
                className={`flex items-center justify-between p-3 rounded-lg ${
                  darkMode ? 'bg-gray-700' : 'bg-gray-50'
                }`}
              >
                <div>
                  <div className={`font-medium ${darkMode ? 'text-white' : 'text-gray-900'}`}>
                    {show.title}
                  </div>
                  <div className={`text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                    {show.year} • {show.episode_count} episodes
                  </div>
                </div>
                <button
                  onClick={() => handleAutoMap(show.plex_id)}
                  className="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded text-sm transition-colors"
                >
                  <MapPin className="h-4 w-4" />
                </button>
              </div>
            ))}
            {unmappedShows.length > 10 && (
              <div className={`text-center text-sm ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                ... and {unmappedShows.length - 10} more shows
              </div>
            )}
          </div>
        </div>
      )}

      {/* Loading State */}
      {loading && (
        <div className="flex justify-center items-center h-32">
          <LoadingSpinner size="lg" />
        </div>
      )}

      {/* No Plex Configuration */}
      {!serverStatus && (
        <div className={`text-center py-12 ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>
          <AlertCircle className="h-12 w-12 mx-auto mb-4 text-yellow-500" />
          <h3 className={`text-lg font-medium mb-2 ${darkMode ? 'text-white' : 'text-gray-900'}`}>
            Plex Server Not Configured
          </h3>
          <p className={darkMode ? 'text-gray-400' : 'text-gray-600'}>
            Configure your Plex server settings to start managing your shows
          </p>
        </div>
      )}
    </div>
  );
}

export default PlexManager; 