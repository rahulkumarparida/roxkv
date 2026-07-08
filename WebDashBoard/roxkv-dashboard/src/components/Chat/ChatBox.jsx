import { useEffect, useRef } from 'react';
import { Wifi } from 'lucide-react';
import ChatMessage from './ChatMessage';
import ChatInput from './ChatInput';
import QuickActions from './QuickActions';
import { useChat } from '../../hooks/useChat';
import { APP_META } from '../../utils/constants';

export default function ChatBox() {
  const { messages } = useChat();
  const scrollRef = useRef(null);

  // Auto-scroll to bottom on new messages
  useEffect(() => {
    if (scrollRef.current) {
      const el = scrollRef.current;
      el.scrollTo({ top: el.scrollHeight, behavior: 'smooth' });
    }
  }, [messages]);

  return (
    <div className="flex flex-col h-full bg-white/[0.02] backdrop-blur-md border border-white/10 rounded-2xl overflow-hidden">
      {/* Header */}
      <div className="px-5 py-3 border-b border-white/10 shrink-0">
        <div className="flex items-center justify-between">
          <h2 className="text-sm font-bold tracking-wider text-text-primary uppercase">
            ROXX AI Master Agent
          </h2>
          <div className="flex items-center gap-1.5">
            <Wifi className="h-3 w-3 text-green-400" />
            <span className="text-xs text-text-secondary">
              Connected to <span className="text-green-400 font-mono">{APP_META.connectionHost}</span>
            </span>
          </div>
        </div>
      </div>

      {/* Chat Messages */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto px-5 py-2 space-y-1">
        {messages.map((msg) => (
          <ChatMessage key={msg.id} message={msg} />
        ))}
      </div>

      {/* Quick Actions */}
      <QuickActions />

      {/* Input */}
      <ChatInput />
    </div>
  );
}
