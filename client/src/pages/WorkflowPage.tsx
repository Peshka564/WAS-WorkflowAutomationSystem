import {
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  Background,
  BackgroundVariant,
  Controls,
  MarkerType,
  ReactFlow,
  type Connection,
  type Edge,
  type EdgeChange,
  type NodeChange,
  type NodeTypes,
} from '@xyflow/react';
import { useCallback, useEffect, useState } from 'react';
import '@xyflow/react/dist/style.css';
import { useMutation, useQuery } from '@tanstack/react-query';
import {
  Box,
  Button,
  TextField,
  Typography,
  Paper,
  Alert,
  Snackbar,
  CircularProgress,
} from '@mui/material';
import SaveIcon from '@mui/icons-material/Save';
import axios from 'axios';
import type {
  EdgeData,
  NodeData,
  WorkflowData,
  WorkflowEdgeDisplay,
  WorkflowNodeDisplay,
  WorkflowNodeDisplaySelector,
} from '../types/workflow';
import { toCreateWorkflowDto } from '../utils/transform';
import { NodeSelector } from '../components/NodeSelector';
import { Node } from '../components/Node';
import { useNavigate, useParams } from 'react-router-dom';

const nodeTypes: NodeTypes = {
  node: Node,
};

interface Props {
  allNodes: WorkflowNodeDisplaySelector[];
}

export function WorkflowPage({ allNodes }: Props) {
  const { id } = useParams<{ id: string }>();
  const workflowId = id ? Number(id) : undefined;

  const [nodes, setNodes] = useState<WorkflowNodeDisplay[]>([]);
  const [edges, setEdges] = useState<WorkflowEdgeDisplay[]>([]);
  const [workflowName, setWorkflowName] = useState('');

  // UI State for feedback
  const [toast, setToast] = useState<{
    open: boolean;
    message: string;
    severity: 'success' | 'error';
  }>({
    open: false,
    message: '',
    severity: 'success',
  });

  const navigate = useNavigate();

  const { data: existingWorkflow, isLoading } = useQuery<
    unknown,
    Error,
    WorkflowData
  >({
    queryKey: ['workflow', workflowId],
    queryFn: async () => {
      if (!workflowId) return null;
      const token = localStorage.getItem('token');
      const res = await axios.get(
        `http://localhost:3000/api/workflows/${workflowId}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      return res.data;
    },
    enabled: !!workflowId, // Only run if ID is provided
  });

  useEffect(() => {
    if (existingWorkflow) {
      setWorkflowName(existingWorkflow.workflow.name);

      // Map Backend Nodes -> Frontend Nodes
      const loadedNodes = existingWorkflow.nodes.map((node: NodeData) => ({
        id: node.display_id,
        position: JSON.parse(node.position),
        type: 'node' as const,
        data: {
          dbId: node.id,
          serviceName: node.service_name,
          taskName: node.task_name,
          type: node.type,
          config: JSON.parse(node.config),
          credentialId: node.credential_id,
        },
      }));
      console.log(loadedNodes);
      setNodes(loadedNodes);

      console.log(existingWorkflow);
      // Map Backend Edges -> Frontend Edges
      const loadedEdges = existingWorkflow.edges.map((edge: EdgeData) => ({
        id: edge.display_id,
        source: edge.node_from,
        target: edge.node_to,
        markerEnd: {
          type: MarkerType.ArrowClosed,
          width: 20,
          height: 20,
          color: '#666',
        },
        style: { strokeWidth: 2 },
      }));
      console.log(loadedEdges);
      setEdges(loadedEdges);
    }
  }, [existingWorkflow]);

  const mutation = useMutation<{ workflowId: number }, Error, void>({
    mutationFn: async () => {
      if (!workflowName.trim()) throw new Error('Workflow name is required');
      if (nodes.length === 0)
        throw new Error('Workflow must have at least one task');

      const payload = toCreateWorkflowDto(nodes, edges, workflowName);
      const token = localStorage.getItem('token');
      const response = await axios.post(
        'http://localhost:3000/api/workflows',
        payload,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      return response.data;
    },
    onSuccess: ({ workflowId }) => {
      setToast({
        open: true,
        message: 'Workflow saved successfully!',
        severity: 'success',
      });
      navigate(`/workflow/${workflowId}`);
    },
    onError: (error: Error) => {
      setToast({ open: true, message: error.message, severity: 'error' });
    },
  });

  const onNodesChange = useCallback(
    (changes: NodeChange<WorkflowNodeDisplay>[]) =>
      setNodes(nds => applyNodeChanges(changes, nds)),
    []
  );

  const onEdgesChange = useCallback(
    (changes: EdgeChange<Edge>[]) =>
      setEdges(eds => applyEdgeChanges(changes, eds)),
    []
  );

  const onConnect = useCallback((params: Connection) => {
    const newEdge = {
      ...params,
      markerEnd: {
        type: MarkerType.ArrowClosed,
        width: 20,
        height: 20,
        color: '#666',
      },
      style: { strokeWidth: 2 },
    };
    setEdges(eds => addEdge(newEdge, eds));
  }, []);

  if (isLoading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="60vh"
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Paper
        elevation={1}
        sx={{
          p: 2,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          zIndex: 10,
        }}
      >
        <Box display="flex" alignItems="center" gap={2}>
          <Typography variant="h6" fontWeight="bold">
            Create Workflow
          </Typography>
          <TextField
            size="small"
            value={workflowName}
            placeholder="Enter workflow name..."
            onChange={e => setWorkflowName(e.target.value)}
            error={mutation.isError && !workflowName}
            sx={{ width: 300 }}
          />
        </Box>

        <Button
          variant="contained"
          startIcon={
            mutation.isPending ? (
              <CircularProgress size={20} color="inherit" />
            ) : (
              <SaveIcon />
            )
          }
          onClick={() => mutation.mutate()}
          disabled={mutation.isPending}
        >
          {mutation.isPending ? 'Saving...' : 'Save Workflow'}
        </Button>
      </Paper>

      <Box sx={{ display: 'flex', flexGrow: 1, overflow: 'hidden' }}>
        <Box
          sx={{
            width: 250,
            borderRight: '1px solid #e0e0e0',
            bgcolor: '#f5f5f5',
            p: 2,
            overflowY: 'auto',
          }}
        >
          <Typography variant="subtitle2" color="textSecondary" sx={{ mb: 2 }}>
            Available Tasks
          </Typography>
          <NodeSelector
            allNodes={allNodes}
            onNodeAdd={change => onNodesChange([change])}
          />
        </Box>

        <Box sx={{ flexGrow: 1, position: 'relative' }}>
          <ReactFlow
            nodes={nodes}
            edges={edges}
            nodeTypes={nodeTypes}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            fitView
            proOptions={{ hideAttribution: true }}
            snapToGrid
          >
            <Background variant={BackgroundVariant.Dots} gap={20} size={1} />
            <Controls />
          </ReactFlow>
        </Box>
      </Box>

      <Snackbar
        open={toast.open}
        autoHideDuration={6000}
        onClose={() => setToast(prev => ({ ...prev, open: false }))}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert severity={toast.severity} variant="filled">
          {toast.message}
        </Alert>
      </Snackbar>
    </Box>
  );
}
