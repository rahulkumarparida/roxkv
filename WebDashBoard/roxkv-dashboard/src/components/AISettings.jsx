import React, { useState, useEffect } from 'react';
import { X, Save, RefreshCw, Cpu, Key, Link as LinkIcon, Database, Check } from 'lucide-react';

export default function AISettings({ isOpen, onClose }) {
  const [config, setConfig] = useState({
    provider: 'ollama',
    model: '',
    endpoint: '',
    apiKey: '',
    temperature: 0.2,
    topP: 0.9,
    maxTokens: 400,
  });
  
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (isOpen) {
      fetchConfig();
    }
  }, [isOpen]);

  const fetchConfig = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await fetch('/api/chat/config');
      if (!response.ok) throw new Error('Failed to load settings');
      const data = await response.json();
      setConfig({
        provider: data.provider || 'ollama',
        model: data.model || '',
        endpoint: data.endpoint || '',
        apiKey: data.apiKey || '',
        temperature: data.temperature || 0.2,
        topP: data.topP || 0.9,
        maxTokens: data.maxTokens || 400,
      });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError('');
    setSaveSuccess(false);

    try {
      const response = await fetch('/api/chat/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          provider: config.provider,
          model: config.model,
          endpoint: config.endpoint,
          apiKey: config.apiKey,
          temperature: parseFloat(config.temperature),
          topP: parseFloat(config.topP),
          maxTokens: parseInt(config.maxTokens, 10),
        }),
      });

      if (!response.ok) {
        const errData = await response.text();
        throw new Error(errData || 'Failed to save settings');
      }
      
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setSaving(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/60 backdrop-blur-sm transition-opacity" 
        onClick={onClose}
      />
      
      {/* Modal */}
      <div className="relative w-full max-w-lg bg-gray-900/90 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl shadow-purple-500/10 overflow-hidden flex flex-col transform transition-all">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-white/5 bg-white/5">
          <div className="flex items-center gap-2">
            <Cpu className="w-5 h-5 text-purple-400" />
            <h2 className="text-lg font-medium text-white tracking-wide">AI Settings</h2>
          </div>
          <button 
            onClick={onClose}
            className="p-1 rounded-full text-gray-400 hover:text-white hover:bg-white/10 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 overflow-y-auto max-h-[70vh] custom-scrollbar">
          {loading ? (
            <div className="flex flex-col items-center justify-center py-12 gap-3 text-gray-400">
              <RefreshCw className="w-6 h-6 animate-spin text-purple-400" />
              <span>Loading configuration...</span>
            </div>
          ) : (
            <form id="ai-settings-form" onSubmit={handleSave} className="space-y-5 text-sm">
              
              {/* Provider Selection */}
              <div className="space-y-1.5">
                <label className="text-gray-300 font-medium block">LLM Provider</label>
                <div className="grid grid-cols-3 gap-2">
                  {['ollama', 'openai', 'claude', 'gemini', 'groq'].map((p) => (
                    <button
                      key={p}
                      type="button"
                      onClick={() => setConfig({ ...config, provider: p })}
                      className={`py-2 px-3 rounded-lg border text-center transition-all ${
                        config.provider === p 
                          ? 'border-purple-500 bg-purple-500/20 text-purple-200 shadow-[0_0_10px_rgba(168,85,247,0.2)]' 
                          : 'border-white/10 bg-white/5 text-gray-400 hover:bg-white/10'
                      }`}
                    >
                      {p.charAt(0).toUpperCase() + p.slice(1)}
                    </button>
                  ))}
                </div>
              </div>

              {/* Model & Endpoint */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="space-y-1.5">
                  <label className="text-gray-300 font-medium flex items-center gap-1.5">
                    <Database className="w-3.5 h-3.5 text-gray-500" /> Model
                  </label>
                  <input 
                    type="text" 
                    value={config.model}
                    onChange={(e) => setConfig({ ...config, model: e.target.value })}
                    className="w-full bg-black/40 border border-white/10 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500/50 focus:ring-1 focus:ring-purple-500/50 transition-all"
                    placeholder={config.provider === 'ollama' ? 'llama3.2:3b' : 'gpt-4o'}
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-gray-300 font-medium flex items-center gap-1.5">
                    <LinkIcon className="w-3.5 h-3.5 text-gray-500" /> Endpoint API
                  </label>
                  <input 
                    type="text" 
                    value={config.endpoint}
                    onChange={(e) => setConfig({ ...config, endpoint: e.target.value })}
                    className="w-full bg-black/40 border border-white/10 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500/50 focus:ring-1 focus:ring-purple-500/50 transition-all"
                    placeholder="https://api.openai.com/v1"
                  />
                </div>
              </div>

              {/* API Key */}
              <div className="space-y-1.5">
                <label className="text-gray-300 font-medium flex items-center gap-1.5">
                  <Key className="w-3.5 h-3.5 text-gray-500" /> API Key
                </label>
                <input 
                  type="password" 
                  value={config.apiKey}
                  onChange={(e) => setConfig({ ...config, apiKey: e.target.value })}
                  className="w-full bg-black/40 border border-white/10 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500/50 focus:ring-1 focus:ring-purple-500/50 transition-all"
                  placeholder={config.provider === 'ollama' ? 'Leave empty for local Ollama' : 'sk-...'}
                />
              </div>

              {/* Advanced Parameters */}
              <div className="pt-2 border-t border-white/5 space-y-4">
                <h3 className="text-gray-400 font-medium">Advanced Parameters</h3>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-gray-400 text-xs">Temperature</label>
                    <input 
                      type="number" 
                      step="0.1" min="0" max="2"
                      value={config.temperature}
                      onChange={(e) => setConfig({ ...config, temperature: e.target.value })}
                      className="w-full bg-black/40 border border-white/10 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-purple-500/50 transition-all"
                    />
                  </div>
                  <div className="space-y-1.5">
                    <label className="text-gray-400 text-xs">Top P</label>
                    <input 
                      type="number" 
                      step="0.1" min="0" max="1"
                      value={config.topP}
                      onChange={(e) => setConfig({ ...config, topP: e.target.value })}
                      className="w-full bg-black/40 border border-white/10 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-purple-500/50 transition-all"
                    />
                  </div>
                  <div className="space-y-1.5">
                    <label className="text-gray-400 text-xs">Max Tokens</label>
                    <input 
                      type="number" 
                      step="1"
                      value={config.maxTokens}
                      onChange={(e) => setConfig({ ...config, maxTokens: e.target.value })}
                      className="w-full bg-black/40 border border-white/10 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-purple-500/50 transition-all"
                    />
                  </div>
                </div>
              </div>

              {/* Error Message */}
              {error && (
                <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-xs">
                  {error}
                </div>
              )}
            </form>
          )}
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-white/5 bg-black/40 flex justify-end gap-3">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 rounded-lg text-sm font-medium text-gray-300 hover:text-white hover:bg-white/5 transition-colors"
          >
            Cancel
          </button>
          <button
            form="ai-settings-form"
            type="submit"
            disabled={loading || saving}
            className="flex items-center gap-2 px-5 py-2 rounded-lg text-sm font-medium bg-purple-600 hover:bg-purple-500 text-white shadow-lg shadow-purple-500/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed active:scale-95"
          >
            {saving ? (
              <RefreshCw className="w-4 h-4 animate-spin" />
            ) : saveSuccess ? (
              <Check className="w-4 h-4" />
            ) : (
              <Save className="w-4 h-4" />
            )}
            {saving ? 'Saving...' : saveSuccess ? 'Saved!' : 'Save Settings'}
          </button>
        </div>
      </div>
    </div>
  );
}
