import React, { useState, useEffect } from 'react';
import { X, Save, RefreshCw, Key, Link as LinkIcon, Building, Sliders, Clock, Check } from 'lucide-react';
import { fetchProviderConfig, saveProviderConfig } from '../../api/chat';

export default function ProviderSettingsModal({ isOpen, provider, onClose, onSaveSuccess }) {
  const [config, setConfig] = useState({
    apiKey: '',
    endpoint: '',
    temperature: 0.2,
    topP: 0.9,
    maxTokens: 400,
    stream: false,
    timeout: 30,
    organization: '',
  });

  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (isOpen && provider) {
      loadConfig();
    }
  }, [isOpen, provider]);

  const loadConfig = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await fetchProviderConfig(provider);
      if (data) {
        setConfig({
          apiKey: data.apiKey || '',
          endpoint: data.endpoint || '',
          temperature: data.temperature !== undefined ? data.temperature : 0.2,
          topP: data.topP !== undefined ? data.topP : 0.9,
          maxTokens: data.maxTokens !== undefined ? data.maxTokens : 400,
          stream: Boolean(data.stream),
          timeout: data.timeout !== undefined ? data.timeout : 30,
          organization: data.organization || '',
        });
      } else {
        // Configuration does not exist: display empty fields
        setConfig({
          apiKey: '',
          endpoint: '',
          temperature: 0.2,
          topP: 0.9,
          maxTokens: 400,
          stream: false,
          timeout: 30,
          organization: '',
        });
      }
    } catch (err) {
      console.warn('Failed to load config, displaying empty fields:', err);
      setConfig({
        apiKey: '',
        endpoint: '',
        temperature: 0.2,
        topP: 0.9,
        maxTokens: 400,
        stream: false,
        timeout: 30,
        organization: '',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError('');
    setSuccess(false);

    try {
      const payload = {
        provider,
        apiKey: config.apiKey,
        endpoint: config.endpoint,
        temperature: parseFloat(config.temperature),
        topP: parseFloat(config.topP),
        maxTokens: parseInt(config.maxTokens, 10),
        stream: Boolean(config.stream),
        timeout: parseInt(config.timeout, 10),
        organization: config.organization,
      };

      await saveProviderConfig(provider, payload);
      setSuccess(true);
      if (onSaveSuccess) onSaveSuccess();
      setTimeout(() => {
        setSuccess(false);
        onClose();
      }, 1000);
    } catch (err) {
      setError(err.message || 'Failed to save configuration');
    } finally {
      setSaving(false);
    }
  };

  if (!isOpen) return null;

  const capitalizedProvider = provider ? provider.charAt(0).toUpperCase() + provider.slice(1) : '';

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/70 backdrop-blur-sm transition-opacity" 
        onClick={onClose}
      />

      {/* Modal Box */}
      <div className="relative w-full max-w-lg bg-gray-900/95 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl shadow-purple-500/10 overflow-hidden flex flex-col transform transition-all">
        
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-white/10 bg-white/5">
          <div className="flex items-center gap-2">
            <Sliders className="w-5 h-5 text-purple-400" />
            <h2 className="text-sm font-semibold text-white tracking-wider uppercase">
              Configure {capitalizedProvider}
            </h2>
          </div>
          <button 
            onClick={onClose}
            className="p-1 rounded-full text-gray-400 hover:text-white hover:bg-white/10 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Form Body */}
        <div className="p-6 overflow-y-auto max-h-[75vh] custom-scrollbar">
          {loading ? (
            <div className="flex flex-col items-center justify-center py-12 gap-3 text-gray-400">
              <RefreshCw className="w-6 h-6 animate-spin text-purple-400" />
              <span className="text-xs font-mono">Loading {capitalizedProvider} settings...</span>
            </div>
          ) : (
            <form id="provider-config-form" onSubmit={handleSave} className="space-y-4 text-xs">
              
              {/* API Key */}
              <div className="space-y-1">
                <label className="text-gray-300 font-medium flex items-center gap-1.5">
                  <Key className="w-3.5 h-3.5 text-purple-400" /> API Key
                </label>
                <input 
                  type="password"
                  value={config.apiKey}
                  onChange={(e) => setConfig({ ...config, apiKey: e.target.value })}
                  placeholder={provider === 'ollama' ? 'Optional for local Ollama' : 'Enter API Key (e.g. sk-...)'}
                  className="w-full bg-black/50 border border-white/10 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500/50 transition-all font-mono"
                />
              </div>

              {/* Endpoint */}
              <div className="space-y-1">
                <label className="text-gray-300 font-medium flex items-center gap-1.5">
                  <LinkIcon className="w-3.5 h-3.5 text-purple-400" /> Endpoint
                </label>
                <input 
                  type="text"
                  value={config.endpoint}
                  onChange={(e) => setConfig({ ...config, endpoint: e.target.value })}
                  placeholder={provider === 'ollama' ? 'http://localhost:11434' : 'Custom endpoint URL (optional)'}
                  className="w-full bg-black/50 border border-white/10 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500/50 transition-all font-mono"
                />
              </div>

              {/* Organization */}
              <div className="space-y-1">
                <label className="text-gray-300 font-medium flex items-center gap-1.5">
                  <Building className="w-3.5 h-3.5 text-purple-400" /> Organization (Optional)
                </label>
                <input 
                  type="text"
                  value={config.organization}
                  onChange={(e) => setConfig({ ...config, organization: e.target.value })}
                  placeholder="org-..."
                  className="w-full bg-black/50 border border-white/10 rounded-lg px-3 py-2 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500/50 transition-all font-mono"
                />
              </div>

              {/* Numeric Parameters */}
              <div className="pt-2 border-t border-white/10 grid grid-cols-3 gap-3">
                <div className="space-y-1">
                  <label className="text-gray-400 text-[10px] uppercase">Temperature</label>
                  <input 
                    type="number"
                    step="0.1"
                    min="0"
                    max="2"
                    value={config.temperature}
                    onChange={(e) => setConfig({ ...config, temperature: e.target.value })}
                    className="w-full bg-black/50 border border-white/10 rounded-lg px-2.5 py-1.5 text-white font-mono focus:outline-none focus:border-purple-500/50"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-gray-400 text-[10px] uppercase">Top P</label>
                  <input 
                    type="number"
                    step="0.1"
                    min="0"
                    max="1"
                    value={config.topP}
                    onChange={(e) => setConfig({ ...config, topP: e.target.value })}
                    className="w-full bg-black/50 border border-white/10 rounded-lg px-2.5 py-1.5 text-white font-mono focus:outline-none focus:border-purple-500/50"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-gray-400 text-[10px] uppercase">Max Tokens</label>
                  <input 
                    type="number"
                    step="1"
                    value={config.maxTokens}
                    onChange={(e) => setConfig({ ...config, maxTokens: e.target.value })}
                    className="w-full bg-black/50 border border-white/10 rounded-lg px-2.5 py-1.5 text-white font-mono focus:outline-none focus:border-purple-500/50"
                  />
                </div>
              </div>

              {/* Stream & Timeout */}
              <div className="grid grid-cols-2 gap-3 pt-2">
                <div className="space-y-1">
                  <label className="text-gray-300 font-medium flex items-center gap-1.5">
                    <Clock className="w-3.5 h-3.5 text-purple-400" /> Timeout (sec)
                  </label>
                  <input 
                    type="number"
                    value={config.timeout}
                    onChange={(e) => setConfig({ ...config, timeout: e.target.value })}
                    className="w-full bg-black/50 border border-white/10 rounded-lg px-3 py-1.5 text-white font-mono focus:outline-none focus:border-purple-500/50"
                  />
                </div>

                <div className="flex items-center pt-5">
                  <label className="flex items-center gap-2 text-gray-300 cursor-pointer select-none">
                    <input 
                      type="checkbox"
                      checked={config.stream}
                      onChange={(e) => setConfig({ ...config, stream: e.target.checked })}
                      className="rounded bg-black/50 border-white/10 text-purple-500 focus:ring-purple-500/50 h-4 w-4"
                    />
                    <span>Stream Response</span>
                  </label>
                </div>
              </div>

              {error && (
                <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-xs">
                  {error}
                </div>
              )}
            </form>
          )}
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-white/10 bg-black/40 flex justify-end gap-3">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 rounded-lg text-xs font-medium text-gray-400 hover:text-white hover:bg-white/5 transition-colors"
          >
            Cancel
          </button>
          <button
            form="provider-config-form"
            type="submit"
            disabled={loading || saving}
            className="flex items-center gap-2 px-5 py-2 rounded-lg text-xs font-medium bg-purple-600 hover:bg-purple-500 text-white shadow-lg shadow-purple-500/20 transition-all disabled:opacity-50 active:scale-95"
          >
            {saving ? (
              <RefreshCw className="w-3.5 h-3.5 animate-spin" />
            ) : success ? (
              <Check className="w-3.5 h-3.5" />
            ) : (
              <Save className="w-3.5 h-3.5" />
            )}
            {saving ? 'Saving...' : success ? 'Saved!' : 'Save Settings'}
          </button>
        </div>
      </div>
    </div>
  );
}
