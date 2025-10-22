import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import TaskForm from '../components/TaskForm'
import { CreateTaskRequest } from '../types/task'

describe('TaskForm', () => {
  it('renders form fields correctly', () => {
    const mockOnSubmit = vi.fn()
    render(<TaskForm onSubmit={mockOnSubmit} />)
    
    expect(screen.getByLabelText(/title/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/priority/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /add task/i })).toBeInTheDocument()
  })

  it('submits form with correct data', async () => {
    const mockOnSubmit = vi.fn()
    render(<TaskForm onSubmit={mockOnSubmit} />)
    
    const titleInput = screen.getByLabelText(/title/i)
    const descriptionInput = screen.getByLabelText(/description/i)
    const prioritySelect = screen.getByLabelText(/priority/i)
    const submitButton = screen.getByRole('button', { name: /add task/i })
    
    fireEvent.change(titleInput, { target: { value: 'Test Task' } })
    fireEvent.change(descriptionInput, { target: { value: 'Test Description' } })
    fireEvent.change(prioritySelect, { target: { value: 'high' } })
    
    fireEvent.click(submitButton)
    
    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith({
        title: 'Test Task',
        description: 'Test Description',
        priority: 'high'
      })
    })
  })

  it('resets form after submission', async () => {
    const mockOnSubmit = vi.fn()
    render(<TaskForm onSubmit={mockOnSubmit} />)
    
    const titleInput = screen.getByLabelText(/title/i) as HTMLInputElement
    const submitButton = screen.getByRole('button', { name: /add task/i })
    
    fireEvent.change(titleInput, { target: { value: 'Test Task' } })
    fireEvent.click(submitButton)
    
    await waitFor(() => {
      expect(titleInput.value).toBe('')
    })
  })

  it('does not submit with empty title', () => {
    const mockOnSubmit = vi.fn()
    render(<TaskForm onSubmit={mockOnSubmit} />)
    
    const submitButton = screen.getByRole('button', { name: /add task/i })
    fireEvent.click(submitButton)
    
    expect(mockOnSubmit).not.toHaveBeenCalled()
  })
})
