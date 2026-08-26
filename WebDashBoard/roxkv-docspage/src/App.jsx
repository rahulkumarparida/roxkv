import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import DocsLayout from './components/docs/DocsLayout';

// RoxKV Pages
import RoxkvOverview from './pages/docs/roxkv/RoxkvOverview';
import RoxkvArchitecture from './pages/docs/roxkv/RoxkvArchitecture';
import CoreDatabase from './pages/docs/roxkv/CoreDatabase';
import TcpServer from './pages/docs/roxkv/TcpServer';
import RespProtocol from './pages/docs/roxkv/RespProtocol';
import CommandIndex from './pages/docs/roxkv/CommandIndex';
import CommandPage from './pages/docs/roxkv/CommandPage';
import PubSub from './pages/docs/roxkv/PubSub';
import Persistence from './pages/docs/roxkv/Persistence';
import TtlExpiration from './pages/docs/roxkv/TtlExpiration';
import Configuration from './pages/docs/roxkv/Configuration';
import Installation from './pages/docs/roxkv/Installation';
import RedisCompat from './pages/docs/roxkv/RedisCompat';
import UsageExamples from './pages/docs/roxkv/UsageExamples';

// RoxAI Pages
import RoxaiOverview from './pages/docs/roxai/RoxaiOverview';
import RoxaiArchitecture from './pages/docs/roxai/RoxaiArchitecture';
import MasterAgent from './pages/docs/roxai/MasterAgent';
import ToolRegistry from './pages/docs/roxai/ToolRegistry';
import ToolRouting from './pages/docs/roxai/ToolRouting';
import ToolExecution from './pages/docs/roxai/ToolExecution';
import LlmLayer from './pages/docs/roxai/LlmLayer';
import Providers from './pages/docs/roxai/Providers';
import Failover from './pages/docs/roxai/Failover';
import Sessions from './pages/docs/roxai/Sessions';
import RoxaiConfiguration from './pages/docs/roxai/RoxaiConfiguration';
import Dashboard from './pages/docs/roxai/Dashboard';
import RoxkvBridge from './pages/docs/roxai/RoxkvBridge';

// Shared Pages
import HighLevelArchitecture from './pages/docs/HighLevelArchitecture';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/docs" element={<DocsLayout />}>
          <Route index element={<Navigate to="/docs/roxkv/overview" replace />} />
          
          {/* High-Level Architecture */}
          <Route path="architecture" element={<HighLevelArchitecture />} />

          {/* Legacy redirects */}
          <Route path="overview" element={<Navigate to="/docs/roxkv/overview" replace />} />
          <Route path="resp" element={<Navigate to="/docs/roxkv/resp" replace />} />
          <Route path="pubsub" element={<Navigate to="/docs/roxkv/pubsub" replace />} />
          <Route path="persistence" element={<Navigate to="/docs/roxkv/persistence" replace />} />
          <Route path="commands" element={<Navigate to="/docs/roxkv/commands" replace />} />
          <Route path="commands/:cmd" element={<Navigate to="/docs/roxkv/commands" replace />} />
          
          {/* RoxKV Routes */}
          <Route path="roxkv/overview" element={<RoxkvOverview />} />
          <Route path="roxkv/architecture" element={<RoxkvArchitecture />} />
          <Route path="roxkv/core" element={<CoreDatabase />} />
          <Route path="roxkv/tcp-server" element={<TcpServer />} />
          <Route path="roxkv/resp" element={<RespProtocol />} />
          <Route path="roxkv/commands" element={<CommandIndex />} />
          <Route path="roxkv/commands/:cmd" element={<CommandPage />} />
          <Route path="roxkv/pubsub" element={<PubSub />} />
          <Route path="roxkv/persistence" element={<Persistence />} />
          <Route path="roxkv/ttl" element={<TtlExpiration />} />
          <Route path="roxkv/configuration" element={<Configuration />} />
          <Route path="roxkv/installation" element={<Installation />} />
          <Route path="roxkv/redis-compat" element={<RedisCompat />} />
          <Route path="roxkv/usage" element={<UsageExamples />} />

          {/* RoxAI Routes */}
          <Route path="roxai/overview" element={<RoxaiOverview />} />
          <Route path="roxai/architecture" element={<RoxaiArchitecture />} />
          <Route path="roxai/master-agent" element={<MasterAgent />} />
          <Route path="roxai/tool-registry" element={<ToolRegistry />} />
          <Route path="roxai/tool-routing" element={<ToolRouting />} />
          <Route path="roxai/tool-execution" element={<ToolExecution />} />
          <Route path="roxai/llm-layer" element={<LlmLayer />} />
          <Route path="roxai/providers" element={<Providers />} />
          <Route path="roxai/failover" element={<Failover />} />
          <Route path="roxai/sessions" element={<Sessions />} />
          <Route path="roxai/configuration" element={<RoxaiConfiguration />} />
          <Route path="roxai/dashboard" element={<Dashboard />} />
          <Route path="roxai/roxkv-bridge" element={<RoxkvBridge />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
