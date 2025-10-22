import React, { useState } from 'react'
import { CreateTaskRequest } from '../types/task'

interface TaskFormProps {
  onSubmit: (task: CreateTaskRequest) => void
}

const TaskForm: React.FC<TaskFormProps> = ({ onSubmit }) => {
  const [formData, setFormData] = useState<CreateTaskRequest>({
    title: '',
    description: '',
    priority: 'medium',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (formData.title.trim()) {
      onSubmit(formData)
      setFormData({ title: '', description: '', priority: 'medium' })
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
  }

  return (
    <form className="task-form" onSubmit={handleSubmit}>
      <h2>Add New Task</h2>
      
      <div className="form-group">
        <label htmlFor="title">Title *</label>
        <input
          type="text"
          id="title"
          name="title"
          value={formData.title}
          onChange={handleChange}
          required
          placeholder="Enter task title"
        />
      </div>

      <div className="form-group">
        <label htmlFor="description">Description</label>
        <textarea
          id="description"
          name="description"
          value={formData.description}
          onChange={handleChange}
          placeholder="Enter task description"
          rows={3}
        />
      </div>

      <div className="form-group">
        <label htmlFor="priority">Priority</label>
        <select
          id="priority"
          name="priority"
          value={formData.priority}
          onChange={handleChange}
        >
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
          <option value="critical">Critical</option>
        </select>
      </div>

      <button type="submit" className="btn btn-primary">
        Add Task
      </button>
    </form>
  )
}

export default TaskForm
