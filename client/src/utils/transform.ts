import type {
  CreateWorkflowEdge,
  CreateWorkflowNode,
  CreateWorkflowPayload,
  WorkflowEdgeDisplay,
  WorkflowNodeDisplay,
} from '../types/workflow';

export function toPayload(
  nodes: WorkflowNodeDisplay[],
  edges: WorkflowEdgeDisplay[],
  name: string
): CreateWorkflowPayload {
  const workflow = {
    name,
  };
  const nodesToSave: CreateWorkflowNode[] = nodes.map(node => ({
    displayId: node.id,
    position: node.position,
    service: node.data.service,
    taskName: node.data.taskName,
  }));
  const edgesToSave: CreateWorkflowEdge[] = edges.map(edge => ({
    from: edge.source,
    to: edge.target,
  }));
  return {
    workflow,
    nodes: nodesToSave,
    edges: edgesToSave,
  };
}
