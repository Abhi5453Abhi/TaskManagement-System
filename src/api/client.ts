import axios from 'axios'
import { Task, CreateTaskRequest, UpdateTaskRequest } from '../types/task'

const api = axios.create({
  baseURL: '/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

export const taskApi = {
  async getAllTasks(): Promise<Task[]> {
    const response = await api.get('/tasks')
    return response.data
  },

  async getTask(id: number): Promise<Task> {
    const response = await api.get(`/tasks/${id}`)
    return response.data
  },

  async createTask(task: CreateTaskRequest): Promise<Task> {
    const response = await api.post('/tasks', task)
    return response.data
  },

  async updateTask(id: number, updates: UpdateTaskRequest): Promise<Task> {
    const response = await api.patch(`/tasks/${id}`, updates)
    return response.data
  },

  async deleteTask(id: number): Promise<void> {
    await api.delete(`/tasks/${id}`)
  },
}
