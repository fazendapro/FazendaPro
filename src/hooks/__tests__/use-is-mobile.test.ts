import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useIsMobile } from '../use-is-mobile'

describe('useIsMobile', () => {
  const originalInnerWidth = window.innerWidth
  const originalAddEventListener = window.addEventListener
  const originalRemoveEventListener = window.removeEventListener

  beforeEach(() => {
    window.addEventListener = vi.fn()
    window.removeEventListener = vi.fn()
  })

  afterEach(() => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: originalInnerWidth,
    })
    window.addEventListener = originalAddEventListener
    window.removeEventListener = originalRemoveEventListener
  })

  it('deve retornar false quando a largura da tela é maior que 768px', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024,
    })

    const { result } = renderHook(() => useIsMobile())

    expect(result.current).toBe(false)
  })

  it('deve retornar true quando a largura da tela é menor que 768px', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 600,
    })

    const { result } = renderHook(() => useIsMobile())

    expect(result.current).toBe(true)
  })

  it('deve retornar true quando a largura da tela é exatamente 767px', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 767,
    })

    const { result } = renderHook(() => useIsMobile())

    expect(result.current).toBe(true)
  })

  it('deve retornar false quando a largura da tela é exatamente 768px', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 768,
    })

    const { result } = renderHook(() => useIsMobile())

    expect(result.current).toBe(false)
  })

  it('deve adicionar listener de resize', () => {
    renderHook(() => useIsMobile())

    expect(window.addEventListener).toHaveBeenCalledWith('resize', expect.any(Function))
  })

  it('deve remover listener de resize quando o componente é desmontado', () => {
    const { unmount } = renderHook(() => useIsMobile())

    unmount()

    expect(window.removeEventListener).toHaveBeenCalledWith('resize', expect.any(Function))
  })

  it('deve atualizar o estado quando a largura da tela muda', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024,
    })

    const { result } = renderHook(() => useIsMobile())

    expect(result.current).toBe(false)

    act(() => {
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 600,
      })

      const resizeHandler = (window.addEventListener as unknown as { mock: { calls: Array<[string, EventListener]> } }).mock.calls.find(
        (call: [string, EventListener]) => call[0] === 'resize'
      )?.[1]

      if (resizeHandler) {
        resizeHandler()
      }
    })

    expect(result.current).toBe(true)
  })
})
