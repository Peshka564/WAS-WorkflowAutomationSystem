import { Box, Button } from '@mui/material';
import type {
  WorkflowNodeDisplay,
  WorkflowNodeDisplaySelector,
} from '../types/workflow';
import type { NodeChange } from '@xyflow/react';
import { v4 } from 'uuid';
import { NodeBase } from './NodeBase';

interface Props {
  allNodes: WorkflowNodeDisplaySelector[];
  onNodeAdd: (change: NodeChange<WorkflowNodeDisplay>) => void;
}

export function NodeSelector({ allNodes, onNodeAdd }: Props) {
  return (
    <Box>
      {allNodes.map((node, idx) => (
        <Button
          key={idx}
          onClick={() =>
            onNodeAdd({
              item: {
                id: v4(),
                position: {
                  x: 0,
                  y: 0,
                },
                type: 'node',
                data: {
                  config: {},
                  dbId: '',
                  serviceName: node.serviceName,
                  taskName: node.taskName,
                  type: node.type,
                },
              },
              type: 'add',
            })
          }
        >
          <NodeBase data={node} />
        </Button>
      ))}
    </Box>
  );
}
