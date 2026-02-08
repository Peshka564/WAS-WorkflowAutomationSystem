import { ReactFlowProvider } from '@xyflow/react';
import type { WorkflowNodeDisplaySelector } from './types/workflow';
import '@xyflow/react/dist/style.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { WorkflowPage } from './pages/WorkflowPage';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { RegisterPage } from './pages/RegisterPage';
import Navbar from './components/Navbar';
import type { ReactNode } from 'react';
import { LoginPage } from './pages/LoginPage';
import { AuthGuard } from './components/AuthGuard';
import Dashboard from './pages/Dashboard';

// TODO: Get these from db
const allNodes: WorkflowNodeDisplaySelector[] = [
  {
    taskName: 'get-email',
    serviceName: 'gmail',
    type: 'listener',
  },
  {
    taskName: 'send-email',
    serviceName: 'gmail',
    type: 'action',
  },
  {
    taskName: 'open-pr',
    serviceName: 'github',
    type: 'action',
  },
  {
    taskName: 'get-issue',
    serviceName: 'github',
    type: 'listener',
  },
  {
    taskName: 'store-file',
    serviceName: 'drive',
    type: 'action',
  },
];

const queryClient = new QueryClient();

const Layout = ({ children }: { children: ReactNode }) => (
  <>
    <Navbar />
    <div style={{ padding: '20px' }}>{children}</div>
  </>
);

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />

          <Route
            path="/dashboard"
            element={
              <AuthGuard>
                <Layout>
                  <Dashboard />
                </Layout>
              </AuthGuard>
            }
          />
          <Route
            path="/create-workflow"
            element={
              <AuthGuard>
                <ReactFlowProvider>
                  <Layout>
                    <WorkflowPage allNodes={allNodes} />
                  </Layout>
                </ReactFlowProvider>
              </AuthGuard>
            }
          />
          <Route
            path="/workflow/:id"
            element={
              <AuthGuard>
                <ReactFlowProvider>
                  <Layout>
                    <WorkflowPage allNodes={allNodes} />
                  </Layout>
                </ReactFlowProvider>
              </AuthGuard>
            }
          />

          {/* <Route
            path="/community"
            element={
              <Layout>
                <Community />
              </Layout>
            }
          /> */}

          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
