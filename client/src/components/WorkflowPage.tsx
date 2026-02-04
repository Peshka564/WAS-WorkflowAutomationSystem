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
} from '@xyflow/react';
import { useCallback, useState } from 'react';
import '@xyflow/react/dist/style.css';
import { useMutation } from '@tanstack/react-query';
import { Box, Button, TextField, Typography } from '@mui/material';
import axios from 'axios';
import type {
  WorkflowEdgeDisplay,
  WorkflowNodeDisplay,
  WorkflowNodeDisplaySelector,
} from '../types/workflow';
import { toPayload } from '../utils/transform';
import { NodeSelector } from './NodeSelector';
import { Node } from './Node';

interface Props {
  allNodes: WorkflowNodeDisplaySelector[];
}

export function WorkflowPage({ allNodes }: Props) {
  const [nodes, setNodes] = useState<WorkflowNodeDisplay[]>([]);
  const [edges, setEdges] = useState<WorkflowEdgeDisplay[]>([]);
  const [workflowName, setWorkflowName] = useState('');

  const mutation = useMutation({
    mutationFn: async ({
      nodes,
      edges,
      name,
    }: {
      nodes: WorkflowNodeDisplay[];
      edges: WorkflowEdgeDisplay[];
      name: string;
    }) => {
      console.log(1);
      if (name.length === 0) {
        throw new Error('You must specify workflow name');
      }
      if (nodes.length === 0) {
        throw new Error('You must have at least one task in the workflow');
      }
      console.log(
        await axios.post(
          'http://localhost:3000/workflows/create',
          JSON.stringify(toPayload(nodes, edges, name)),
          {
            headers: {
              'Content-Type': 'application/json',
            },
          }
        )
      );
    },
  });

  const onNodesChange = useCallback(
    (changes: NodeChange<WorkflowNodeDisplay>[]) =>
      setNodes(nodesSnapshot => applyNodeChanges(changes, nodesSnapshot)),
    []
  );
  const onEdgesChange = useCallback((changes: EdgeChange<Edge>[]) => {
    setEdges(edgesSnapshot => applyEdgeChanges(changes, edgesSnapshot));
  }, []);
  const onConnect = useCallback((params: Connection | Edge) => {
    setEdges(edgesSnapshot => addEdge(params, edgesSnapshot));
    setEdges(edgesSnapshot =>
      edgesSnapshot.map(edge => ({
        ...edge,
        markerEnd: {
          type: MarkerType.ArrowClosed,
          width: 20,
          height: 20,
          color: '#gray',
        },
      }))
    );
  }, []);

  return (
    <>
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <Typography fontSize={18}>Create workflow</Typography>
        <TextField
          value={workflowName}
          placeholder="Workflow Name"
          onChange={e => setWorkflowName(e.target.value)}
        />
      </Box>
      <NodeSelector
        allNodes={allNodes}
        onNodeAdd={change => onNodesChange([change])}
      />
      <Box style={{ width: '90%', height: '500px', border: '1px solid black' }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={{ node: Node }}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          fitView
          proOptions={{ hideAttribution: true }}
        >
          <Background variant={BackgroundVariant.Dots} gap={20} size={1} />
          <Controls />
        </ReactFlow>
      </Box>
      <Button
        onClick={() => mutation.mutate({ nodes, edges, name: workflowName })}
      >
        {mutation.isPending ? 'Loading...' : 'Save workflow'}
      </Button>
    </>
  );
}
