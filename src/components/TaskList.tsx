import React from 'react'
import { Task } from '../types/task'

interface TaskListProps {
  tasks: Task[]
  onUpdate: (id: number, updates: Partial<Task>) => void
  onDelete: (id: number) => void
}

const TaskList: React.FC<TaskListProps> = ({ tasks, onUpdate, onDelete }) => {
  const getPriorityClass = (priority: string) => {
    return `priority-${priority}`
  }

  const getStatusClass = (status: string) => {
    return `status-${status}`
  }

  const handleStatusChange = (id: number, newStatus: 'todo' | 'doing' | 'done') => {
    onUpdate(id, { status: newStatus })
  }

  const handlePriorityChange = (id: number, newPriority: 'low' | 'medium' | 'high' | 'critical') => {
    onUpdate(id, { priority: newPriority })
  }

  if (tasks.length === 0) {
    return (
      <div className="task-list">
        <h2>Tasks</h2>
        <div className="empty-state">
          <p>No tasks yet. Create your first task above!</p>
        </div>
      </div>
    )
  }

  return (
    <div className="task-list">
      <h2>Tasks ({tasks.length})</h2>
      
      <div className="tasks-grid">
        {tasks.map(task => (
          <div key={task.id} className={`task-card ${getPriorityClass(task.priority)}`}>
            <div className="task-header">
              <h3 className="task-title">{task.title}</h3>
              <button 
                className="btn btn-danger btn-sm"
                onClick={() => onDelete(task.id)}
                aria-label={`Delete task: ${task.title}`}
              >
                ×
              </button>
            </div>
            
            {task.description && (
              <p className="task-description">{task.description}</p>
            )}
            
            <div className="task-meta">
              <div className="task-controls">
                <div className="control-group">
                  <label>Status:</label>
                  <select 
                    value={task.status} 
                    onChange={(e) => handleStatusChange(task.id, e.target.value as 'todo' | 'doing' | 'done')}
                    className={`status-select ${getStatusClass(task.status)}`}
                  >
                    <option value="todo">To Do</option>
                    <option value="doing">Doing</option>
                    <option value="done">Done</option>
                  </select>
                </div>
                
                <div className="control-group">
                  <label>Priority:</label>
                  <select 
                    value={task.priority} 
                    onChange={(e) => handlePriorityChange(task.id, e.target.value as 'low' | 'medium' | 'high' | 'critical')}
                    className={`priority-select ${getPriorityClass(task.priority)}`}
                  >
                    <option value="low">Low</option>
                    <option value="medium">Medium</option>
                    <option value="high">High</option>
                    <option value="critical">Critical</option>
                  </select>
                </div>
              </div>
              
              <div className="task-dates">
                <small>Created: {new Date(task.created_at).toLocaleDateString()}</small>
                {task.updated_at !== task.created_at && (
                  <small>Updated: {new Date(task.updated_at).toLocaleDateString()}</small>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export default TaskList
