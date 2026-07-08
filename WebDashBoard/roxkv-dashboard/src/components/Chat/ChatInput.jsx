import { useState, useRef, useCallback } from 'react';
import { Send } from 'lucide-react';
import { useChat } from '../../hooks/useChat';
import { cn } from '../../utils/classNames';

export default function ChatInput() {
  const [text, setText] = useState('');
  const { sendMessage, isGenerating } = useChat();
  const textareaRef = useRef(null);

  const adjustHeight = useCallback(() => {
    const el = textareaRef.current;
    if (el) {
      el.style.height = 'auto';
      el.style.height = Math.min(el.scrollHeight, 120) + 'px';
    }
  }, []);

  const handleSend = useCallback(() => {
    if (text.trim() && !isGenerating) {
      sendMessage(text);
      setText('');
      if (textareaRef.current) textareaRef.current.style.height = 'auto';
    }
  }, [text, isGenerating, sendMessage]);

  const handleKeyDown = useCallback((e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  }, [handleSend]);

  return (
    <div className="px-4 pb-3 pt-2">
      <div className={cn(
        'relative flex items-end gap-2',
        'bg-white/5 border border-white/10 rounded-2xl',
        'px-4 py-3',
        'transition-all duration-300',
        'focus-within:border-purple-500/30 focus-within:shadow-[0_0_10px_rgba(168,85,247,0.15)]'
      )}>
        <textarea
          ref={textareaRef}
          value={text}
          onChange={(e) => { setText(e.target.value); adjustHeight(); }}
          onKeyDown={handleKeyDown}
          placeholder="Ask the master agent anything about the database..."
          rows={1}
          className={cn(
            'flex-1 bg-transparent text-sm text-text-primary resize-none',
            'placeholder:text-text-muted',
            'focus:outline-none',
            'max-h-[120px]'
          )}
          disabled={isGenerating}
        />
        <button
          onClick={handleSend}
          disabled={!text.trim() || isGenerating}
          className={cn(
            'shrink-0 h-9 w-9 rounded-xl',
            'flex items-center justify-center',
            'transition-all duration-200',
            'active:scale-90',
            text.trim() && !isGenerating
              ? 'bg-purple-500 text-white hover:bg-purple-400 hover:shadow-[0_0_15px_rgba(168,85,247,0.5)]'
              : 'bg-white/5 text-text-muted cursor-not-allowed'
          )}
        >
          <Send className="h-4 w-4" />
        </button>
      </div>
      <p className="text-center text-[10px] text-text-muted mt-1.5">
        Press Enter to send {'\u2022'} Shift+Enter for new line
      </p>
    </div>
  );
}
