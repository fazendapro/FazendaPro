import { describe, it, expect } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useModal } from '../useModal'

describe('useModal', () => {
  it('deve inicializar com estado fechado por padrão', () => {
    const { result } = renderHook(() => useModal())

    expect(result.current.isOpen).toBe(false)
  })

  it('deve inicializar com estado customizado', () => {
    const { result } = renderHook(() => useModal(true))

    expect(result.current.isOpen).toBe(true)
  })

  it('deve abrir o modal quando onOpen é chamado', () => {
    const { result } = renderHook(() => useModal())

    act(() => {
      result.current.onOpen()
    })

    expect(result.current.isOpen).toBe(true)
  })

  it('deve fechar o modal quando onClose é chamado', () => {
    const { result } = renderHook(() => useModal(true))

    act(() => {
      result.current.onClose()
    })

    expect(result.current.isOpen).toBe(false)
  })

  it('deve alternar o estado quando onToggle é chamado', () => {
    const { result } = renderHook(() => useModal())

    act(() => {
      result.current.onToggle()
    })
    expect(result.current.isOpen).toBe(true)

    act(() => {
      result.current.onToggle()
    })
    expect(result.current.isOpen).toBe(false)
  })

  it('deve manter referências estáveis das funções', () => {
    const { result, rerender } = renderHook(() => useModal())

    const firstOnOpen = result.current.onOpen
    const firstOnClose = result.current.onClose
    const firstOnToggle = result.current.onToggle

    rerender()

    expect(result.current.onOpen).toBe(firstOnOpen)
    expect(result.current.onClose).toBe(firstOnClose)
    expect(result.current.onToggle).toBe(firstOnToggle)
  })

  it('deve funcionar corretamente com múltiplas chamadas', () => {
    const { result } = renderHook(() => useModal())

    act(() => {
      result.current.onOpen()
    })
    expect(result.current.isOpen).toBe(true)

    act(() => {
      result.current.onClose()
    })
    expect(result.current.isOpen).toBe(false)

    act(() => {
      result.current.onToggle()
    })
    expect(result.current.isOpen).toBe(true)

    act(() => {
      result.current.onToggle()
    })
    expect(result.current.isOpen).toBe(false)
  })
})
