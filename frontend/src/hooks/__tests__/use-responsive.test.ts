import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { renderHook } from '@testing-library/react'
import { useResponsive } from '../use-responsive'

describe('useResponsive', () => {
  const originalInnerWidth = window.innerWidth
  const originalInnerHeight = window.innerHeight

  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: originalInnerWidth,
    })
    Object.defineProperty(window, 'innerHeight', {
      writable: true,
      configurable: true,
      value: originalInnerHeight,
    })
  })

  it('deve retornar isMobile true quando largura é menor que 768', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 600,
    })

    const { result } = renderHook(() => useResponsive())

    expect(result.current.isMobile).toBe(true)
    expect(result.current.isTablet).toBe(false)
    expect(result.current.isDesktop).toBe(false)
    expect(result.current.isLargeDesktop).toBe(false)
    expect(result.current.screenWidth).toBe(600)
  })

  it('deve retornar isTablet true quando largura está entre 768 e 1024', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 900,
    })

    const { result } = renderHook(() => useResponsive())

    expect(result.current.isMobile).toBe(false)
    expect(result.current.isTablet).toBe(true)
    expect(result.current.isDesktop).toBe(false)
    expect(result.current.isLargeDesktop).toBe(false)
    expect(result.current.screenWidth).toBe(900)
  })

  it('deve retornar isDesktop true quando largura está entre 1024 e 1440', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1200,
    })

    const { result } = renderHook(() => useResponsive())

    expect(result.current.isMobile).toBe(false)
    expect(result.current.isTablet).toBe(false)
    expect(result.current.isDesktop).toBe(true)
    expect(result.current.isLargeDesktop).toBe(false)
    expect(result.current.screenWidth).toBe(1200)
  })

  it('deve retornar isLargeDesktop true quando largura é maior ou igual a 1440', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1920,
    })

    const { result } = renderHook(() => useResponsive())

    expect(result.current.isMobile).toBe(false)
    expect(result.current.isTablet).toBe(false)
    expect(result.current.isDesktop).toBe(false)
    expect(result.current.isLargeDesktop).toBe(true)
    expect(result.current.screenWidth).toBe(1920)
  })

  it('deve atualizar valores quando window é redimensionado', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 600,
    })

    const { result, rerender } = renderHook(() => useResponsive())

    expect(result.current.isMobile).toBe(true)

    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1200,
    })

    window.dispatchEvent(new Event('resize'))
    rerender()

    expect(result.current.isDesktop).toBe(true)
    expect(result.current.screenWidth).toBe(1200)
  })
})

