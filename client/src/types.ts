// Log entry structure matching the Go server's model.Log
export interface LogEntry {
  level: 'INFO' | 'WARN' | 'ERROR';
  message: string;
  timestamp: string; // ISO 8601 timestamp string
}

// WebSocket connection states
export enum ConnectionState {
  CONNECTING = 'connecting',
  CONNECTED = 'connected',
  DISCONNECTED = 'disconnected',
  ERROR = 'error'
}

// WebSocket message types with better validation
export interface WebSocketMessage {
  type: 'log' | 'status' | 'error' | 'ping' | 'pong' | 'pause' | 'resume';
  data: LogEntry | string | any;
}

// Validate WebSocket message
export function isValidWebSocketMessage(message: any): message is WebSocketMessage {
  return (
    typeof message === 'object' &&
    message !== null &&
    typeof message.type === 'string' &&
    ['log', 'status', 'error', 'ping', 'pong', 'pause', 'resume'].includes(message.type) &&
    'data' in message
  );
}

// Validate log entry
export function isValidLogEntry(data: any): data is LogEntry {
  return (
    typeof data === 'object' &&
    data !== null &&
    typeof data.level === 'string' &&
    ['INFO', 'WARN', 'ERROR'].includes(data.level) &&
    typeof data.message === 'string' &&
    typeof data.timestamp === 'string'
  );
}

// Connection status information
export interface ConnectionStatus {
  state: ConnectionState;
  lastConnected?: Date;
  lastError?: string;
  reconnectAttempts: number;
}

// App state interface
export interface AppState {
  logs: LogEntry[];
  isPaused: boolean;
  connectionStatus: ConnectionStatus;
  autoScroll: boolean;
  maxLogs: number; // Maximum number of logs to keep in memory
}

// Log level colors for UI
export const LOG_LEVEL_COLORS = {
  INFO: '#00ff00',   // Green
  WARN: '#ffff00',   // Yellow
  ERROR: '#ff0000'   // Red
} as const;

// WebSocket configuration
export interface WebSocketConfig {
  url: string;
  reconnectInterval: number;
  maxReconnectAttempts: number;
  heartbeatInterval: number;
}

// Default WebSocket configuration
export const DEFAULT_WS_CONFIG: WebSocketConfig = {
  url: 'ws://localhost:8080/ws',
  reconnectInterval: 5000, // 5 seconds
  maxReconnectAttempts: 10,
  heartbeatInterval: 30000 // 30 seconds
};

// Log filtering options
export interface LogFilter {
  levels: ('INFO' | 'WARN' | 'ERROR')[];
  searchText: string;
  startTime?: Date;
  endTime?: Date;
}

// Export all types
export type { LogEntry as Log };
export type { ConnectionStatus as Status };
export type { AppState as State };
export type { WebSocketConfig as WSConfig };
export type { LogFilter as Filter };
