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

// TODO: Get these from db
const allNodes: WorkflowNodeDisplaySelector[] = [
  {
    taskName: 'get-email',
    service: 'gmail',
    type: 'listener',
  },
  {
    taskName: 'send-email',
    service: 'gmail',
    type: 'action',
  },
  {
    taskName: 'open-pr',
    service: 'github',
    type: 'action',
  },
  {
    taskName: 'get-issue',
    service: 'github',
    type: 'listener',
  },
  {
    taskName: 'store-file',
    service: 'drive',
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
          {/* Public Routes */}
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/login" element={<LoginPage />} />

          {/* Community Page (Public) */}
          {/* <Route
            path="/community"
            element={
              <Layout>
                <Community />
              </Layout>
            }
          /> */}

          {/* Protected Routes (Workflows) */}
          <Route
            path="/dashboard"
            element={
              <AuthGuard>
                <Layout>
                  <WorkflowPage allNodes={allNodes} />
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

          {/* Dynamic Route for specific workflow */}
          {/* <Route
            path="/workflow/:id"
            element={
              <Layout>
                <WorkflowDetail />
              </Layout>
            }
          /> */}

          {/* Default Redirect: Send unknown routes to Login or Community */}
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
