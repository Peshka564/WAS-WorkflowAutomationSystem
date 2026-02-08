export interface WorkflowNodeDisplay {
  id: string;
  position: {
    x: number;
    y: number;
  };
  type: 'node';
  data: {
    dbId: string;
    serviceName: string;
    taskName: string;
    credentialId?: number;
    type: 'listener' | 'action' | 'transformer';
    config: Record<string, unknown>;
  };
}

export interface WorkflowEdgeDisplay {
  id: string;
  source: string;
  target: string;
}

export type WorkflowNodeDisplaySelector = {
  taskName: string;
  serviceName: string;
  type: 'listener' | 'action' | 'transformer';
};

export interface CreateWorkflowNode {
  displayId: string;
  serviceName: string;
  taskName: string;
  type: 'listener' | 'action' | 'transformer';
  position: string;
  config: string;
  credential_id?: number;
}

export interface CreateWorkflowEdge {
  from: string;
  to: string;
  displayId: string;
}

export interface CreateWorkflowPayload {
  workflow: {
    name: string;
  };
  nodes: CreateWorkflowNode[];
  edges: CreateWorkflowEdge[];
}

export interface Workflow {
  id: number;
  created_at: Date;
  updated_at: Date;
  name: string;
  active: boolean;
  user_id: number;
}

export interface NodeData {
  id: string;
  display_id: string;
  service_name: string;
  task_name: string;
  workflow_id: number;
  type: 'listener' | 'action' | 'transformer';
  position: string;
  config: string;
  credential_id?: number;
}

export interface EdgeData {
  id: string;
  display_id: string;
  workflow_id: number;
  node_from: string;
  node_to: string;
}

export interface WorkflowData {
  workflow: Workflow;
  nodes: NodeData[];
  edges: EdgeData[];
}
