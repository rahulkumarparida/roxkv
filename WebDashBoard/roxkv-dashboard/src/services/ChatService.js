/**
 * @fileoverview Chat service layer between API and Context.
 * Orchestrates chat flow: send → receive → stream token-by-token.
 */

import { sendChatMessage } from '../api/chat';

/**
 * @typedef {Object} ChatMessage
 * @property {string} id - Unique message ID
 * @property {'user'|'agent'} role - Message sender role
 * @property {string} content - Message content (markdown for agent)
 * @property {string} timestamp - Display timestamp
 * @property {boolean} [isStreaming] - Whether the message is still being streamed
 * @property {boolean} [isToolExecution] - Whether agent is executing a tool
 */

class ChatService {
  /**
   * Send a user message and get a response.
   * @param {string} message - The user's message text
   * @param {function(ChatMessage): void} onUpdate - Called with streaming updates
   * @returns {Promise<ChatMessage>} The final agent response message
   */
  async sendMessage(message, onUpdate) {
    const response = await sendChatMessage({ message, context_keys: [] });
    const fullText = response.response;

    const agentMessage = {
      id: `msg_agent_${Date.now()}`,
      role: 'agent',
      content: '',
      timestamp: new Date().toLocaleTimeString('en-US', {
        hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit',
      }),
      isStreaming: true,
      isToolExecution: false,
    };

    // Stream the response character-by-character in chunks
    await this._streamResponse(fullText, agentMessage, onUpdate);

    agentMessage.isStreaming = false;
    agentMessage.content = fullText;
    onUpdate({ ...agentMessage });

    return agentMessage;
  }

  /**
   * Simulates token-by-token streaming for realistic LLM output.
   * @param {string} fullText
   * @param {ChatMessage} message
   * @param {function(ChatMessage): void} onUpdate
   */
  async _streamResponse(fullText, message, onUpdate) {
    const chunkSize = 3 + Math.floor(Math.random() * 5);
    let index = 0;

    while (index < fullText.length) {
      const end = Math.min(index + chunkSize, fullText.length);
      message.content = fullText.slice(0, end);
      onUpdate({ ...message });
      index = end;

      // Vary delay for natural feel
      const char = fullText[index - 1];
      const delay = char === '\n' ? 40 : char === ' ' ? 15 : 10 + Math.random() * 20;
      await new Promise((r) => setTimeout(r, delay));
    }
  }
}

const chatService = new ChatService();
export default chatService;
