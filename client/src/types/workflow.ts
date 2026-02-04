export interface WorkflowNodeDisplay {
  id: string;
  position: {
    x: number;
    y: number;
  };
  type: 'node';
  data: {
    service: string; //gmail, drive, etc...
    type: 'listener' | 'action' | 'transformer';
    isAuthenticated?: boolean;
    taskName: string;
  };
}

export interface WorkflowEdgeDisplay {
  id: string;
  source: string;
  target: string;
}

export type WorkflowNodeDisplaySelector = Omit<
  WorkflowNodeDisplay['data'],
  'isAuthenticated'
>;

export interface CreateWorkflowNode {
  displayId: string;
  position: {
    x: number;
    y: number;
  };
  service: string; //gmail, drive, etc...
  isAuthenticated?: boolean;
  taskName: string;
}

export interface CreateWorkflowEdge {
  from: string;
  to: string;
}

export interface CreateWorkflowPayload {
  workflow: {
    name: string;
  };
  nodes: CreateWorkflowNode[];
  edges: CreateWorkflowEdge[];
}
