import type { WorkflowNodeDisplaySelector } from '../types/workflow';
import { Handle, Position } from '@xyflow/react';
import { NodeBase } from './NodeBase';

export function Node({ data }: { data: WorkflowNodeDisplaySelector }) {
  return (
    <>
      <NodeBase data={data} />
      {data.type !== 'action' && (
        <Handle type="source" position={Position.Right} />
      )}
      {data.type !== 'listener' && (
        <Handle type="target" position={Position.Left} />
      )}
    </>
  );
}
