export interface Task {
  id: number
  title: string
  description: string
  status: 'todo' | 'doing' | 'done'
  priority: 'low' | 'medium' | 'high' | 'critical'
  created_at: string
  updated_at: string
}

export interface CreateTaskRequest {
  title: string
  description: string
  priority: 'low' | 'medium' | 'high' | 'critical'
}

export interface UpdateTaskRequest {
  title?: string
  description?: string
  status?: 'todo' | 'doing' | 'done'
  priority?: 'low' | 'medium' | 'high' | 'critical'
}
