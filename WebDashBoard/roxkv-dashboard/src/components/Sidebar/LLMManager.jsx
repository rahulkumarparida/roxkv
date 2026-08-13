import React, { useState, useEffect } from 'react';
import { Bot, Settings, RefreshCw } from 'lucide-react';
import GlassCard from '../Common/GlassCard';
import StatusBadge from '../Common/StatusBadge';
import ProviderSettingsModal from './ProviderSettingsModal';
import {
  fetchProviders,
  fetchProviderModels,
  fetchProviderStats,
  selectActiveProvider,
} from '../../api/chat';

export default function LLMManager() {
  const [providers, setProviders] = useState(['ollama', 'openai', 'gemini', 'claude', 'groq']);
  const [models, setModels] = useState([]);

  const [currentProvider, setCurrentProvider] = useState('ollama');
  const [currentModel, setCurrentModel] = useState('llama3.2:3b');
  const [connectionStatus, setConnectionStatus] = useState('connected');
  const [totalQueries, setTotalQueries] = useState(0);

  const [selectedProviderInput, setSelectedProviderInput] = useState('ollama');
  const [selectedModelInput, setSelectedModelInput] = useState('llama3.2:3b');

  const [loadingModels, setLoadingModels] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);


  const loadProvidersList = async () => {
    try {
      const list = await fetchProviders();
      if (Array.isArray(list) && list.length > 0) {
        setProviders(list);
      }
    } catch (err) {
      console.warn('Failed to load providers list, using defaults:', err);
    }
  };

  const refreshActiveStats = async (providerToFetch) => {
    try {
      const target = providerToFetch || currentProvider || 'ollama';
      console.log("Target :",target);
      
      const stats = await fetchProviderStats(target);
      if (stats) {
        if (stats.currentProvider && stats.currentProvider !== currentProvider) {
          setCurrentProvider(stats.currentProvider);
          setSelectedProviderInput(stats.currentProvider);
          loadModelsForProvider(stats.currentProvider);
        } else if (stats.currentProvider) {
          setCurrentProvider(stats.currentProvider);
        }

        if (stats.currentModel) {
          setCurrentModel(stats.currentModel);
          setSelectedModelInput(stats.currentModel);
        }

        if (stats.connectionStatus) {
          setConnectionStatus(stats.connectionStatus);
        }

        if (stats.totalQueries !== undefined) {
          setTotalQueries(stats.totalQueries);
        }
      }
    } catch (err) {
      console.warn('Failed to fetch provider stats:', err);
    }
  };

  const loadModelsForProvider = async (providerName) => {
    setLoadingModels(true);
    try {
      const mList = await fetchProviderModels(providerName);
      if (Array.isArray(mList) && mList.length > 0) {
        setModels(mList);
        return mList;
      } else {
        setModels([]);
        return [];
      }
    } catch (err) {
      console.warn(`Failed to fetch models for ${providerName}:`, err);
      setModels([]);
      return [];
    } finally {
      setLoadingModels(false);
    }
  };

  const handleProviderChange = async (newProvider) => {
    
    setSelectedProviderInput(newProvider);
    const fetchedModels = await loadModelsForProvider(newProvider);
    console.log("Fetch:",fetchedModels);
    
    const newModel = fetchedModels.length > 0 ? fetchedModels[0] : '';
    setSelectedModelInput(newModel);

    try {
      await selectActiveProvider(newProvider, newModel);
      await refreshActiveStats(newProvider);
      await loadModelsForProvider(newProvider)
    } catch (err) {
      console.error('Failed to set active provider:', err);
    }
  };

  const handleModelChange = async (newModel) => {
    setSelectedModelInput(newModel);
    try {
      await selectActiveProvider(selectedProviderInput, newModel);
      await refreshActiveStats(selectedProviderInput);
       await loadModelsForProvider(selectedProviderInput)
    } catch (err) {
      console.error('Failed to set active model:', err);
    }
  };




  // Initial load and sync listeners
  useEffect(() => {
    loadProvidersList();
    loadModelsForProvider(currentProvider)
    const handleProviderChange = () => {
      refreshActiveStats();
    };

    if (typeof window !== 'undefined') {
      window.addEventListener('llm_provider_changed', handleProviderChange);
    }

    return () => {
      if (typeof window !== 'undefined') {
        window.removeEventListener('llm_provider_changed', handleProviderChange);
      }

    };
  }, []);

  return (
    <GlassCard className="p-4" hover>
      {/* Header */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <Bot className="h-4 w-4 text-purple-400" />
          <h3 className="text-xs font-semibold tracking-wider text-text-primary uppercase">
            LLM Manager
          </h3>
        </div>
        <StatusBadge
          status={connectionStatus === 'connected' ? 'live' : 'offline'}
          label={connectionStatus.toUpperCase()}
        />
      </div>

      {/* Active Display */}
      <div className="bg-white/5 border border-white/5 rounded-xl p-3 mb-3 space-y-2">
        <div className="flex justify-between items-center text-xs">
          <span className="text-text-muted">Current Provider</span>
          <span className="font-mono text-purple-300 font-semibold uppercase">
            {currentProvider}
          </span>
        </div>
        <div className="flex justify-between items-center text-xs">
          <span className="text-text-muted">Current Model</span>
          <span className="font-mono text-text-primary truncate max-w-[140px]" title={currentModel}>
            {currentModel || 'None'}
          </span>
        </div>
      </div>

      {/* Selectors */}
      <div className="space-y-3 mb-4">
        <div>
          <label className="text-[10px] text-text-muted uppercase tracking-wider block mb-1">
            Select Provider
          </label>
          <select
            value={selectedProviderInput}
            onChange={(e) => handleProviderChange(e.target.value)}
            className="w-full bg-black/40 border border-white/10 rounded-lg px-3 py-1.5 text-xs text-text-primary focus:outline-none focus:border-purple-500/50"
          >
            {providers.map((p) => (
              <option key={p} value={p} className="bg-gray-900 text-white">
                {p.charAt(0).toUpperCase() + p.slice(1)}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="text-[10px] text-text-muted uppercase tracking-wider block mb-1">
            Select Model
          </label>
          <div className="relative">
            <select
              value={selectedModelInput}
              onChange={(e) => handleModelChange(e.target.value)}
              disabled={loadingModels || models.length === 0}
              className="w-full bg-black/40 border border-white/10 rounded-lg px-3 py-1.5 text-xs text-text-primary focus:outline-none focus:border-purple-500/50 disabled:opacity-50 appearance-none pr-7"
            >
              {models.length === 0 ? (
                <option value="" className="bg-gray-900 text-white">
                  No models available
                </option>
              ) : (
                models.map((m) => (
                  <option key={m} value={m} className="bg-gray-900 text-white">
                    {m}
                  </option>
                ))
              )}
            </select>
            {loadingModels && (
              <RefreshCw className="w-3 h-3 animate-spin absolute right-2.5 top-2.5 text-purple-400" />
            )}
          </div>
        </div>
      </div>

      {/* Footer / Buttons */}
      <div className="flex items-center justify-between pt-3 border-t border-white/5">
        <button
          onClick={() => setIsModalOpen(true)}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-purple-500/20 text-purple-300 text-xs hover:bg-purple-500/30 transition-colors border border-purple-500/30 active:scale-95"
        >
          <Settings className="h-3.5 w-3.5" /> Configure Provider
        </button>

        <div className="text-right">
          <div className="text-[10px] text-text-muted uppercase tracking-wider mb-0.5">Queries</div>
          <div className="text-xs font-mono text-purple-300 font-bold">{totalQueries}</div>
        </div>
      </div>

      {/* Modal */}
      <ProviderSettingsModal
        isOpen={isModalOpen}
        provider={selectedProviderInput}
        onClose={() => setIsModalOpen(false)}
        onSaveSuccess={() => refreshActiveStats(selectedProviderInput)}
      />
    </GlassCard>
  );
}
