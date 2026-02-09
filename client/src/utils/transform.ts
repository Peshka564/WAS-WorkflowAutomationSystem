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
  name: string,
  workflowId?: number
): CreateWorkflowPayload {
  const workflow = {
    id: workflowId,
    name,
  };

  const nodesToSave: CreateWorkflowNode[] = nodes.map(node => ({
    id: node.data.dbId,
    displayId: node.id,
    serviceName: node.data.serviceName,
    taskName: node.data.taskName,
    type: node.data.type,
    config: JSON.stringify(node.data.config) || '{}',
    position: JSON.stringify(node.position),
    credential_id: node.data.credentialId,
  }));

  const edgesToSave = edges.map(edge => ({
    from: edge.source,
    to: edge.target,
    displayId: edge.id,
    id: edge.data?.dbId,
  })) satisfies CreateWorkflowEdge[];

  return {
    workflow,
    nodes: nodesToSave,
    edges: edgesToSave,
  };
}
