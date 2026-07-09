import { memo } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Bot, Check } from 'lucide-react';
import TypingIndicator from '../Common/TypingIndicator';
import { cn } from '../../utils/classNames';

/**
 * Chat message component with distinct user/agent styles.
 * Agent messages render full markdown with code blocks and tables.
 */
const ChatMessage = memo(function ChatMessage({ message }) {
  const isUser = message.role === 'user';
  const isToolExecution = message.isToolExecution;

  if (isToolExecution) {
    return (
      <div className="flex items-start gap-3 py-3">
        <div className="shrink-0 h-7 w-7 rounded-lg bg-purple-500/20 flex items-center justify-center">
          <Bot className="h-4 w-4 text-purple-400" />
        </div>
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <span className="text-sm font-medium text-purple-300">Master Agent</span>
            <span className="text-[10px] text-text-muted font-mono">{message.timestamp}</span>
          </div>
          <div className="shimmer rounded-xl py-3 px-4 max-w-md">
            <div className="flex items-center gap-2 text-xs text-text-secondary">
              <div className="h-3 w-3 animate-spin rounded-full border border-purple-500/30 border-t-purple-400" />
              Agent is querying database...
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (isUser) {
    return (
      <div className="flex justify-end py-3">
        <div className="max-w-[75%]">
          <div className={cn(
            'bg-purple-500/15 border border-purple-500/20 rounded-2xl rounded-tr-sm',
            'py-2.5 px-4 text-sm text-text-primary'
          )}>
            {message.content}
          </div>
          <div className="flex items-center justify-end gap-1 mt-1">
            <span className="text-[10px] text-text-muted font-mono">{message.timestamp}</span>
            <Check className="h-3 w-3 text-purple-400" />
          </div>
        </div>
      </div>
    );
  }

  // Agent message
  return (
    <div className="flex items-start gap-3 py-3">
      <div className="shrink-0 h-7 w-7 rounded-lg bg-purple-500/20 flex items-center justify-center mt-0.5">
        <Bot className="h-4 w-4 text-purple-400" />
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <span className="text-sm font-medium text-purple-300">Master Agent</span>
          <span className="text-[10px] text-text-muted font-mono">{message.timestamp}</span>
        </div>
        <div className={cn(
          'bg-white/[0.03] border border-white/5 rounded-2xl rounded-tl-sm',
          'py-3 px-4'
        )}>
          <div className="markdown-body">
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{message.content}</ReactMarkdown>
          </div>
          {message.isStreaming && (
            <div className="mt-2 border-t border-white/5 pt-2">
              <TypingIndicator />
            </div>
          )}
        </div>

      </div>
    </div>
  );
});

export default ChatMessage;
