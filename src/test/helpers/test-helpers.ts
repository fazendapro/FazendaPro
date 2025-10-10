import { vi } from 'vitest'

export const createMockFunction = <T extends (...args: any[]) => any>(
  returnValue?: ReturnType<T>
) => {
  return vi.fn().mockReturnValue(returnValue)
}

export const createMockAsyncFunction = <T extends (...args: any[]) => Promise<any>>(
  returnValue?: Awaited<ReturnType<T>>
) => {
  return vi.fn().mockResolvedValue(returnValue)
}

export const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

export const createTestData = <T>(data: Partial<T>): T => {
  return data as T
}

export const createKeyboardEvent = (key: string, options: KeyboardEventInit = {}) => {
  return new KeyboardEvent('keydown', { key, ...options })
}

export const createMouseEvent = (type: string, options: MouseEventInit = {}) => {
  return new MouseEvent(type, options)
}

export const mockWindowSize = (width: number, height: number) => {
  Object.defineProperty(window, 'innerWidth', {
    writable: true,
    configurable: true,
    value: width,
  })
  Object.defineProperty(window, 'innerHeight', {
    writable: true,
    configurable: true,
    value: height,
  })
  
  window.dispatchEvent(new Event('resize'))
}

export const mockScroll = (scrollTop: number, scrollLeft: number = 0) => {
  Object.defineProperty(window, 'scrollY', {
    writable: true,
    configurable: true,
    value: scrollTop,
  })
  Object.defineProperty(window, 'scrollX', {
    writable: true,
    configurable: true,
    value: scrollLeft,
  })
  
  window.dispatchEvent(new Event('scroll'))
}

export const mockIntersectionObserver = (_isIntersecting: boolean = true) => {
  const mockIntersectionObserver = vi.fn()
  mockIntersectionObserver.mockReturnValue({
    observe: () => null,
    unobserve: () => null,
    disconnect: () => null,
  })
  window.IntersectionObserver = mockIntersectionObserver
  
  return mockIntersectionObserver
}

export const mockResizeObserver = () => {
  const mockResizeObserver = vi.fn()
  mockResizeObserver.mockReturnValue({
    observe: () => null,
    unobserve: () => null,
    disconnect: () => null,
  })
  window.ResizeObserver = mockResizeObserver
  
  return mockResizeObserver
}

export const mockMatchMedia = (matches: boolean = false) => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
      matches,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    })),
  })
}

export const mockLocalStorage = () => {
  const store: Record<string, string> = {}
  
  const localStorage = {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      Object.keys(store).forEach(key => delete store[key])
    }),
    length: 0,
    key: vi.fn(),
  }
  
  Object.defineProperty(window, 'localStorage', {
    value: localStorage,
    writable: true,
  })
  
  return localStorage
}

export const mockSessionStorage = () => {
  const store: Record<string, string> = {}
  
  const sessionStorage = {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      Object.keys(store).forEach(key => delete store[key])
    }),
    length: 0,
    key: vi.fn(),
  }
  
  Object.defineProperty(window, 'sessionStorage', {
    value: sessionStorage,
    writable: true,
  })
  
  return sessionStorage
}

export const waitForElement = async (selector: string, timeout: number = 1000) => {
  const startTime = Date.now()
  
  while (Date.now() - startTime < timeout) {
    const element = document.querySelector(selector)
    if (element) return element
    await delay(10)
  }
  
  throw new Error(`Element with selector "${selector}" not found within ${timeout}ms`)
}

export const createFileUploadEvent = (files: File[]) => {
  const event = new Event('change', { bubbles: true })
  Object.defineProperty(event, 'target', {
    value: {
      files,
    },
    writable: false,
  })
  return event
}

export const createMockFile = (name: string, type: string, content: string = '') => {
  const file = new File([content], name, { type })
  return file
}

export const createDragEvent = (type: string, dataTransfer: DataTransfer) => {
  const event = new Event(type, { bubbles: true })
  Object.defineProperty(event, 'dataTransfer', {
    value: dataTransfer,
    writable: false,
  })
  return event
}

export const createMockDataTransfer = (files: File[] = []) => {
  return {
    files,
    items: files.map(file => ({ kind: 'file', type: file.type, getAsFile: () => file })),
    types: ['Files'],
    effectAllowed: 'all',
    dropEffect: 'none',
    clearData: vi.fn(),
    getData: vi.fn(),
    setData: vi.fn(),
    setDragImage: vi.fn(),
  } as any
}
