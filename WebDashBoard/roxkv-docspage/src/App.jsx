import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import DocsLayout from './components/docs/DocsLayout';
import Overview from './pages/docs/Overview';
import Architecture from './pages/docs/Architecture';
import RespProtocol from './pages/docs/RespProtocol';
import PubSub from './pages/docs/PubSub';
import Persistence from './pages/docs/Persistence';
import CommandIndex from './pages/docs/CommandIndex';
import CommandPage from './pages/docs/CommandPage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/docs" element={<DocsLayout />}>
          <Route index element={<Navigate to="/docs/overview" replace />} />
          <Route path="overview" element={<Overview />} />
          <Route path="architecture" element={<Architecture />} />
          <Route path="resp" element={<RespProtocol />} />
          <Route path="pubsub" element={<PubSub />} />
          <Route path="persistence" element={<Persistence />} />
          <Route path="commands" element={<CommandIndex />} />
          <Route path="commands/:cmd" element={<CommandPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
