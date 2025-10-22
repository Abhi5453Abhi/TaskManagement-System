import React, { useState, useEffect } from 'react'
import TaskList from './components/TaskList'
import TaskForm from './components/TaskForm'
import { Task, CreateTaskRequest } from './types/task'
import { taskApi } from './api/client'

function App() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadTasks()
  }, [])

  const loadTasks = async () => {
    try {
      setLoading(true)
      setError(null)
      const fetchedTasks = await taskApi.getAllTasks()
      setTasks(fetchedTasks)
    } catch (err) {
      setError('Failed to load tasks')
      console.error('Error loading tasks:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateTask = async (taskData: CreateTaskRequest) => {
    try {
      const newTask = await taskApi.createTask(taskData)
      setTasks(prev => [...prev, newTask])
    } catch (err) {
      setError('Failed to create task')
      console.error('Error creating task:', err)
    }
  }

  const handleUpdateTask = async (id: number, updates: Partial<Task>) => {
    try {
      const updatedTask = await taskApi.updateTask(id, updates)
      setTasks(prev => prev.map(task => task.id === id ? updatedTask : task))
    } catch (err) {
      setError('Failed to update task')
      console.error('Error updating task:', err)
    }
  }

  const handleDeleteTask = async (id: number) => {
    try {
      await taskApi.deleteTask(id)
      setTasks(prev => prev.filter(task => task.id !== id))
    } catch (err) {
      setError('Failed to delete task')
      console.error('Error deleting task:', err)
    }
  }

  if (loading) {
    return (
      <div className="container">
        <div className="loading">Loading tasks...</div>
      </div>
    )
  }

  return (
    <div className="container">
      <header className="header">
        <h1>Task Manager</h1>
        {error && <div className="error">{error}</div>}
      </header>
      
      <main className="main">
        <TaskForm onSubmit={handleCreateTask} />
        <TaskList 
          tasks={tasks}
          onUpdate={handleUpdateTask}
          onDelete={handleDeleteTask}
        />
      </main>
    </div>
  )
}

export default App
