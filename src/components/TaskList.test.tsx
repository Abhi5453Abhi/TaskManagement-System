import React from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import TaskList from '../components/TaskList'
import { Task } from '../types/task'

const mockTasks: Task[] = [
  {
    id: 1,
    title: 'Test Task 1',
    description: 'Description 1',
    status: 'todo',
    priority: 'high',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z'
  },
  {
    id: 2,
    title: 'Test Task 2',
    description: 'Description 2',
    status: 'doing',
    priority: 'medium',
    created_at: '2024-01-02T00:00:00Z',
    updated_at: '2024-01-02T00:00:00Z'
  }
]

describe('TaskList', () => {
  it('renders empty state when no tasks', () => {
    const mockOnUpdate = vi.fn()
    const mockOnDelete = vi.fn()
    
    render(<TaskList tasks={[]} onUpdate={mockOnUpdate} onDelete={mockOnDelete} />)
    
    expect(screen.getByText(/no tasks yet/i)).toBeInTheDocument()
  })

  it('renders tasks correctly', () => {
    const mockOnUpdate = vi.fn()
    const mockOnDelete = vi.fn()
    
    render(<TaskList tasks={mockTasks} onUpdate={mockOnUpdate} onDelete={mockOnDelete} />)
    
    expect(screen.getByText('Test Task 1')).toBeInTheDocument()
    expect(screen.getByText('Test Task 2')).toBeInTheDocument()
    expect(screen.getByText(/tasks \(2\)/i)).toBeInTheDocument()
  })

  it('calls onUpdate when status changes', () => {
    const mockOnUpdate = vi.fn()
    const mockOnDelete = vi.fn()
    
    render(<TaskList tasks={mockTasks} onUpdate={mockOnUpdate} onDelete={mockOnDelete} />)
    
    const statusSelects = screen.getAllByDisplayValue(/todo|doing/i)
    fireEvent.change(statusSelects[0], { target: { value: 'done' } })
    
    expect(mockOnUpdate).toHaveBeenCalledWith(1, { status: 'done' })
  })

  it('calls onUpdate when priority changes', () => {
    const mockOnUpdate = vi.fn()
    const mockOnDelete = vi.fn()
    
    render(<TaskList tasks={mockTasks} onUpdate={mockOnUpdate} onDelete={mockOnDelete} />)
    
    const prioritySelects = screen.getAllByDisplayValue(/high|medium/i)
    fireEvent.change(prioritySelects[0], { target: { value: 'critical' } })
    
    expect(mockOnUpdate).toHaveBeenCalledWith(1, { priority: 'critical' })
  })

  it('calls onDelete when delete button clicked', () => {
    const mockOnUpdate = vi.fn()
    const mockOnDelete = vi.fn()
    
    render(<TaskList tasks={mockTasks} onUpdate={mockOnUpdate} onDelete={mockOnDelete} />)
    
    const deleteButtons = screen.getAllByRole('button', { name: /delete task/i })
    fireEvent.click(deleteButtons[0])
    
    expect(mockOnDelete).toHaveBeenCalledWith(1)
  })

  it('applies correct priority classes', () => {
    const mockOnUpdate = vi.fn()
    const mockOnDelete = vi.fn()
    
    render(<TaskList tasks={mockTasks} onUpdate={mockOnUpdate} onDelete={mockOnDelete} />)
    
    const taskCards = screen.getAllByRole('generic')
    const highPriorityCard = taskCards.find(card => 
      card.textContent?.includes('Test Task 1') && 
      card.className.includes('priority-high')
    )
    
    expect(highPriorityCard).toBeInTheDocument()
  })
})
