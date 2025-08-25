import { useState, useEffect, useRef, useCallback } from 'react';
import {
  LogEntry,
  ConnectionState,
  ConnectionStatus,
  WebSocketConfig,
  DEFAULT_WS_CONFIG,
  WebSocketMessage,
  isValidWebSocketMessage,
  isValidLogEntry
} from '../types';

interface UseWebSocketReturn {
  logs: LogEntry[];
  isConnected: boolean;
  connectionStatus: ConnectionStatus;
  sendMessage: (message: any) => void;
  reconnect: () => void;
  clearLogs: () => void;
}

export function useWebSocket(config: Partial<WebSocketConfig> = {}): UseWebSocketReturn {
  // Merge default config with provided config
  const wsConfig = { ...DEFAULT_WS_CONFIG, ...config };
  
  // State
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>({
    state: ConnectionState.DISCONNECTED,
    reconnectAttempts: 0
  });
  
  // Refs
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number | null>(null);
  const heartbeatIntervalRef = useRef<number | null>(null);
  const wsConfigRef = useRef(wsConfig); // Store config in ref to prevent recreation
  
  // Update config ref when config changes
  useEffect(() => {
    wsConfigRef.current = wsConfig;
  }, [wsConfig]);
  
  // Derived state
  const isConnected = connectionStatus.state === ConnectionState.CONNECTED;
  
  // Clear logs function
  const clearLogs = useCallback(() => {
    setLogs([]);
  }, []);
  
  // Add log function
  const addLog = useCallback((log: LogEntry) => {
    setLogs(prev => {
      const newLogs = [...prev, log];
      // Keep only last 1000 logs to prevent memory issues
      return newLogs.slice(-1000);
    });
  }, []);
  
  // Send message function
  const sendMessage = useCallback((message: any) => {
    if (!wsRef.current) {
      console.error(`WEBSOCKET: No WebSocket connection available`);
      return;
    }
    
    const state = wsRef.current.readyState;
    
    if (state === WebSocket.OPEN) {
      try {
        console.log(`WEBSOCKET: Sending message:`, message);
        wsRef.current.send(JSON.stringify(message));
        console.log(`WEBSOCKET: Message sent successfully`);
      } catch (error) {
        console.error(`WEBSOCKET: Failed to send message:`, error);
      }
    } else {
      const stateNames: Record<number, string> = {
        [WebSocket.CONNECTING]: 'CONNECTING',
        [WebSocket.OPEN]: 'OPEN',
        [WebSocket.CLOSING]: 'CLOSING',
        [WebSocket.CLOSED]: 'CLOSED'
      };
              console.error(`WEBSOCKET: Cannot send message - WebSocket state: ${stateNames[state] || 'UNKNOWN'}`);
      
      // Auto-reconnect if connection is lost
      if (state === WebSocket.CLOSED && connectionStatus.state === ConnectionState.CONNECTED) {
        console.log(`WEBSOCKET: Connection lost, attempting to reconnect...`);
        // Note: reconnect function will be defined later
      }
    }
  }, [connectionStatus.state]);
  
  // Reconnect function
  const reconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
    
    setConnectionStatus(prev => ({
      ...prev,
      state: ConnectionState.CONNECTING,
      reconnectAttempts: prev.reconnectAttempts + 1
    }));
    
    connect();
  }, []); // No dependencies needed since connect is stable
  
  // Connect function
  const connect = useCallback(() => {
    try {
      // Close existing connection
      if (wsRef.current) {
        wsRef.current.close();
      }
      
      // Create new connection
      const ws = new WebSocket(wsConfigRef.current.url);
      wsRef.current = ws;
      
      // Connection opened
      ws.onopen = () => {
        console.log('WebSocket connected to:', wsConfigRef.current.url);
        setConnectionStatus(prev => ({
          ...prev,
          state: ConnectionState.CONNECTED,
          lastConnected: new Date(),
          reconnectAttempts: 0
        }));
        
        // Start heartbeat
        heartbeatIntervalRef.current = setInterval(() => {
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: 'ping', data: 'heartbeat' }));
          }
        }, wsConfigRef.current.heartbeatInterval);
      };
      
      // Connection closed
      ws.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        setConnectionStatus(prev => ({
          ...prev,
          state: ConnectionState.DISCONNECTED
        }));
        
        // Stop heartbeat
        if (heartbeatIntervalRef.current) {
          clearInterval(heartbeatIntervalRef.current);
        }
        
        // Schedule reconnection if not manually closed
        if (event.code !== 1000) { // 1000 = normal closure
          scheduleReconnect();
        }
      };
      
      // Message received
      ws.onmessage = (event) => {
        try {
          const rawMessage = JSON.parse(event.data);
          
          // Validate message structure
          if (!isValidWebSocketMessage(rawMessage)) {
            console.error('Invalid WebSocket message structure:', rawMessage);
            return;
          }
          
          const message: WebSocketMessage = rawMessage;
          
          if (message.type === 'log') {
            // Validate log entry data
            if (isValidLogEntry(message.data)) {
              addLog(message.data);
            } else {
              console.error('Invalid log entry data:', message.data);
            }
          } else if (message.type === 'ping') {
            // Server ping, respond with pong
            console.log('Ping received from server, sending pong');
            try {
              ws.send(JSON.stringify({ type: 'pong', data: 'heartbeat' }));
            } catch (error) {
              console.error('Failed to send pong response:', error);
            }
          } else if (message.type === 'pong') {
            // Heartbeat response, connection is alive
            console.log('Heartbeat received');
          } else if (message.type === 'pause' || message.type === 'resume') {
            // Handle pause/resume confirmations from server
            console.log(`Server confirmed: ${message.type}`);
          } else {
            console.log('Received message:', message.type, message.data);
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error, 'Raw data:', event.data);
        }
      };
      
      // Connection error
      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionStatus(prev => ({
          ...prev,
          state: ConnectionState.ERROR,
          lastError: 'Connection error occurred'
        }));
      };
      
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      setConnectionStatus(prev => ({
        ...prev,
        state: ConnectionState.ERROR,
        lastError: 'Failed to create connection'
      }));
    }
  }, []); // Remove addLog dependency to prevent recreation
  
  // Schedule reconnection with exponential backoff
  const scheduleReconnect = useCallback(() => {
    if (connectionStatus.reconnectAttempts >= wsConfigRef.current.maxReconnectAttempts) {
      console.log('Max reconnection attempts reached');
      return;
    }
    
    const delay = Math.min(1000 * Math.pow(2, connectionStatus.reconnectAttempts), 30000);
    console.log(`Scheduling reconnection in ${delay}ms (attempt ${connectionStatus.reconnectAttempts + 1})`);
    
    reconnectTimeoutRef.current = setTimeout(() => {
      reconnect();
    }, delay);
  }, [connectionStatus.reconnectAttempts]); // Remove reconnect dependency to prevent recreation
  
  // Initial connection
  useEffect(() => {
    // Only connect if no connection exists
    if (!wsRef.current || wsRef.current.readyState === WebSocket.CLOSED) {
      connect();
    }
    
    // Cleanup on unmount
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (heartbeatIntervalRef.current) {
        clearInterval(heartbeatIntervalRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, []); // Empty dependency array - only run once on mount
  
  return {
    logs,
    isConnected,
    connectionStatus,
    sendMessage,
    reconnect,
    clearLogs,
  };
}
