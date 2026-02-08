import type {
  CreateWorkflowEdge,
  CreateWorkflowNode,
  CreateWorkflowPayload,
  WorkflowEdgeDisplay,
  WorkflowNodeDisplay,
} from '../types/workflow';

export function toCreateWorkflowDto(
  nodes: WorkflowNodeDisplay[],
  edges: WorkflowEdgeDisplay[],
  name: string
): CreateWorkflowPayload {
  const workflow = {
    name,
  };

  const nodesToSave: CreateWorkflowNode[] = nodes.map(node => ({
    displayId: node.id,
    serviceName: node.data.serviceName,
    taskName: node.data.taskName,
    type: node.data.type,
    config: JSON.stringify(node.data.config) || '{}',
    position: JSON.stringify(node.position),
    credential_id: node.data.credentialId,
  }));

  const edgesToSave: CreateWorkflowEdge[] = edges.map(edge => ({
    from: edge.source,
    to: edge.target,
    displayId: edge.id,
  }));

  return {
    workflow,
    nodes: nodesToSave,
    edges: edgesToSave,
  };
}
