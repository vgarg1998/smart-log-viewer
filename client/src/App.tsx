import { useState, useRef } from 'react';
import { useWebSocket } from './hooks/useWebSocket';
import { ConnectionState } from './types';

function App() {
  const {
    logs,
    isConnected,
    connectionStatus,
    sendMessage,
    clearLogs
  } = useWebSocket();
  
  const [isPaused, setIsPaused] = useState(false);
  const [selectedLogLevels, setSelectedLogLevels] = useState<Set<string>>(new Set(['all']));
  const logContainerRef = useRef<HTMLDivElement>(null);

  const handlePauseResume = () => {
    const newPausedState = !isPaused;
    console.log(`CLIENT: Pause state changing from ${isPaused} to ${newPausedState}`);
    
    // Send pause/resume message to server
    if (newPausedState) {
      // Send pause message
      console.log('CLIENT: Sending PAUSE message to server');
      sendMessage({ type: 'pause', data: null });
    } else {
      // Send resume message
              console.log('CLIENT: Sending RESUME message to server');
      sendMessage({ type: 'resume', data: null });
    }
    
    // Update UI state
    setIsPaused(newPausedState);
    console.log(`CLIENT: UI state updated to: ${newPausedState}`);
  };

  const handleClear = () => {
    clearLogs();
  };







  const toggleLogLevel = (level: string) => {
    setSelectedLogLevels(prev => {
      console.log(`FILTER: Toggling level "${level}", current state:`, Array.from(prev));
      
      let newSet: Set<string>;
      
      if (level === 'all') {
        // If "all" is clicked and it's already selected, deselect everything
        if (prev.has('all')) {
          newSet = new Set();
          console.log('FILTER: "all" was selected, now deselecting everything');
        } else {
          // If "all" is clicked and it's not selected, select only "all" (deselect others)
          newSet = new Set(['all']);
          console.log('FILTER: "all" was not selected, now selecting only "all"');
        }
      } else {
        // If a specific level is clicked
        if (prev.has('all')) {
          // If "all" was selected, remove it and select only this level
          newSet = new Set([level]);
          console.log(`FILTER: "all" was selected, now selecting only "${level}"`);
        } else {
          // Toggle this specific level
          newSet = new Set(prev);
          if (newSet.has(level)) {
            newSet.delete(level);
            console.log(`FILTER: "${level}" was selected, now deselecting it`);
          } else {
            newSet.add(level);
            console.log(`FILTER: "${level}" was not selected, now selecting it`);
          }
          
          // If no levels selected, select "all"
          if (newSet.size === 0) {
            newSet = new Set(['all']);
            console.log('FILTER: No levels selected, now selecting "all"');
          }
        }
      }
      
      console.log('FILTER: New state:', Array.from(newSet));
      
      // Force a re-render by ensuring we return a completely new Set
      return new Set(newSet);
    });
  };

  const getFilteredLogs = () => {
    if (selectedLogLevels.has('all')) {
      return logs;
    }
    return logs.filter(log => selectedLogLevels.has(log.level));
  };

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString();
  };

  const getLogLevelClass = (level: string) => {
    return level.toLowerCase();
  };



  return (
    <div className="container">
      {/* Header */}
      <div className="header">
        <h1>Smart Log Viewer</h1>
      </div>

      {/* Connection Status */}
      <div className="connection-status">
        <span className={isConnected ? 'connected' : 'disconnected'}>
          {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
        </span>
        {connectionStatus.state === ConnectionState.CONNECTING && (
          <span style={{ marginLeft: '10px', color: '#ffff00' }}>
            Connecting...
          </span>
        )}
        {connectionStatus.state === ConnectionState.ERROR && (
          <span style={{ marginLeft: '10px', color: '#ff6666' }}>
            Error: {connectionStatus.lastError}
          </span>
        )}
        {connectionStatus.reconnectAttempts > 0 && (
          <span style={{ marginLeft: '10px', color: '#ffaa00' }}>
            Retry {connectionStatus.reconnectAttempts}/{10}
          </span>
        )}
      </div>

      {/* Control Buttons */}
      <div className="controls">
        <button 
          className={`btn ${isPaused ? 'btn-primary' : 'btn-secondary'}`}
          onClick={handlePauseResume}
        >
          {isPaused ? 'Resume' : 'Pause'}
        </button>
        <button 
          className="btn btn-danger"
          onClick={handleClear}
        >
          Clear
        </button>



      </div>

      {/* Log Level Filter */}
      <div className="log-filter">
        <span className="filter-label">Filter Logs:</span>
        {['all', 'INFO', 'WARN', 'ERROR'].map(level => (
          <button
            key={`${level}-${Array.from(selectedLogLevels).join('-')}`}
            className={`filter-btn ${selectedLogLevels.has(level) ? 'active' : ''}`}
            onClick={() => toggleLogLevel(level)}
            style={{
              backgroundColor: selectedLogLevels.has(level) ? '#007bff' : '#f8f9fa',
              color: selectedLogLevels.has(level) ? 'white' : '#212529',
              border: selectedLogLevels.has(level) ? '2px solid #0056b3' : '1px solid #dee2e6',
              fontWeight: selectedLogLevels.has(level) ? 'bold' : 'normal'
            }}
          >
            {level === 'all' ? '📋 All' : level}
          </button>
        ))}
        <span className="filter-count">
          Showing {getFilteredLogs().length} of {logs.length} logs
        </span>
        <div style={{ fontSize: '10px', color: '#666', marginTop: '5px' }}>
          Active filters: {Array.from(selectedLogLevels).join(', ')}
        </div>
        <div style={{ fontSize: '10px', color: '#999', marginTop: '2px' }}>
          Debug - Set size: {selectedLogLevels.size}, Has 'all': {selectedLogLevels.has('all').toString()}
        </div>
        <div style={{ fontSize: '10px', color: '#999', marginTop: '2px' }}>
          Button States: All({selectedLogLevels.has('all') ? '✓' : '✗'}) INFO({selectedLogLevels.has('INFO') ? '✓' : '✗'}) WARN({selectedLogLevels.has('WARN') ? '✓' : '✗'}) ERROR({selectedLogLevels.has('ERROR') ? '✓' : '✗'})
        </div>
      </div>

      {/* Log Display Area */}
      <div 
        ref={logContainerRef}
        className="log-container"
      >
        {logs.length === 0 ? (
          <div style={{ textAlign: 'center', color: '#666', marginTop: '50px' }}>
            <p>No logs yet. Connect to the server to see real-time logs.</p>
            <p style={{ fontSize: '12px', marginTop: '10px' }}>
              Server: ws://localhost:8080/ws
            </p>
            <p style={{ fontSize: '12px', marginTop: '10px' }}>
              Click "Test Log" to add a sample log entry
            </p>
          </div>
        ) : (
          <>
            {getFilteredLogs().slice().reverse().map((log, index) => (
              <div key={`${log.timestamp}-${index}`} className={`log-entry ${getLogLevelClass(log.level)}`}>
                <span className="timestamp">{formatTimestamp(log.timestamp)}</span>
                <span className={`level ${getLogLevelClass(log.level)}`}>
                  {log.level}
                </span>
                <span className="message">{log.message}</span>
              </div>
            ))}
            <div style={{ textAlign: 'center', color: '#666', marginTop: '20px', fontSize: '12px' }}>
              {isPaused ? 'Log display is paused' : `Showing ${getFilteredLogs().length} of ${logs.length} log entries`}
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default App;
