import { useChat } from '../../hooks/useChat';
import { QUICK_ACTIONS } from '../../utils/constants';
import { cn } from '../../utils/classNames';

export default function QuickActions() {
  const { sendMessage, isGenerating } = useChat();

  return (
    <div className="px-4 py-2 border-t border-white/5">
      <div className="flex gap-2 overflow-x-auto pb-1">
        {QUICK_ACTIONS.map((action) => (
          <button
            key={action.label}
            onClick={() => !isGenerating && sendMessage(action.message)}
            disabled={isGenerating}
            className={cn(
              'shrink-0 px-3 py-1.5 rounded-xl text-xs',
              'bg-white/5 border border-white/10 text-text-secondary',
              'transition-all duration-200',
              'hover:bg-purple-500/10 hover:border-purple-500/20 hover:text-purple-300',
              'hover:shadow-[0_0_10px_rgba(168,85,247,0.2)]',
              'active:scale-95',
              'disabled:opacity-50 disabled:cursor-not-allowed'
            )}
          >
            {action.label}
          </button>
        ))}
      </div>
    </div>
  );
}
