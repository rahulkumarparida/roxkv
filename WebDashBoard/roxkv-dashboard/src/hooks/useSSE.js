import { useContext } from 'react';
import { SSEContext } from '../context/SSEContext';

export function useSSE() {
  const context = useContext(SSEContext);
  if (!context) throw new Error('useSSE must be used within an SSEProvider');
  return context;
}
