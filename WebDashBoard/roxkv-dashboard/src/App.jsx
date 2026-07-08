import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { SSEProvider } from './context/SSEContext';
import { DashboardProvider } from './context/DashboardContext';
import { ChatProvider } from './context/ChatContext';
import DashboardLayout from './layouts/DashboardLayout';
import Dashboard from './pages/Dashboard';

/**
 * Root App component.
 * Wraps everything in context providers and sets up routing.
 */
export default function App() {
  return (
    <BrowserRouter >
      <SSEProvider>
        <DashboardProvider>
          <ChatProvider>
            <DashboardLayout>
              <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="*" element={<Dashboard />} />
              </Routes>
            </DashboardLayout>
          </ChatProvider>
        </DashboardProvider>
      </SSEProvider>
    </BrowserRouter>
  );
}
