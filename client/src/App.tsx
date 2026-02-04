import { ReactFlowProvider } from '@xyflow/react';
import type { WorkflowNodeDisplaySelector } from './types/workflow';
import '@xyflow/react/dist/style.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { WorkflowPage } from './components/WorkflowPage';

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

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ReactFlowProvider>
        <WorkflowPage allNodes={allNodes} />
      </ReactFlowProvider>
    </QueryClientProvider>
  );
}
