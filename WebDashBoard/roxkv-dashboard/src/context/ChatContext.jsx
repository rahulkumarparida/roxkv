import { createContext, useState, useCallback, useRef } from 'react';
import chatService from '../services/ChatService';
import { formatTimestamp } from '../utils/formatTime';

export const ChatContext = createContext(null);

const WELCOME_MESSAGE = {
  id: 'msg_welcome',
  role: 'agent',
  content: 'Welcome to **RoxKV AI Master Agent**. I can help you monitor and manage your database. Try asking about storage health, TTL stats, snapshots, or active connections.',
  timestamp: formatTimestamp(new Date()),
  isStreaming: false,
};

/**
 * ChatProvider manages chat messages, generation state, and message flow.
 */
export function ChatProvider({ children }) {
  const [messages, setMessages] = useState([WELCOME_MESSAGE]);
  const [isGenerating, setIsGenerating] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState('connected');
  const streamingRef = useRef(null);

  const sendMessage = useCallback(async (text) => {
    if (!text.trim() || isGenerating) return;

    const userMsg = {
      id: `msg_user_${Date.now()}`, role: 'user', content: text.trim(),
      timestamp: formatTimestamp(new Date()), isStreaming: false,
    };
    setMessages((prev) => [...prev, userMsg]);
    setIsGenerating(true);

    try {
      // Show tool execution shimmer briefly
      const toolMsg = {
        id: `msg_tool_${Date.now()}`, role: 'agent', content: '',
        timestamp: formatTimestamp(new Date()), isStreaming: true, isToolExecution: true,
      };
      setMessages((prev) => [...prev, toolMsg]);
      await new Promise((r) => setTimeout(r, 800));

      // Stream the actual response
      await chatService.sendMessage(text, (update) => {
        streamingRef.current = update;
        setMessages((prev) => {
          const filtered = prev.filter((m) => m.id !== toolMsg.id && m.id !== update.id);
          return [...filtered, update];
        });
      });
    } catch (err) {
      console.error('[Chat] Error:', err);
      const errorMsg = {
        id: `msg_error_${Date.now()}`, role: 'agent',
        content: 'An error occurred while processing your request. Please try again.',
        timestamp: formatTimestamp(new Date()), isStreaming: false,
      };
      setMessages((prev) => [...prev.filter((m) => !m.isToolExecution), errorMsg]);
    } finally {
      setIsGenerating(false);
      streamingRef.current = null;
    }
  }, [isGenerating]);

  const clearHistory = useCallback(() => {
    setMessages([WELCOME_MESSAGE]);
  }, []);

  return (
    <ChatContext.Provider value={{ messages, isGenerating, connectionStatus, sendMessage, clearHistory }}>
      {children}
    </ChatContext.Provider>
  );
}
