import { vi } from 'vitest'

Object.defineProperty(window, 'URL', {
  value: {
    createObjectURL: vi.fn(() => 'mock-object-url'),
    revokeObjectURL: vi.fn(),
  },
  writable: true,
})

global.FileReader = vi.fn().mockImplementation(() => ({
  readAsDataURL: vi.fn(),
  onload: null,
  result: 'data:image/jpeg;base64,test',
})) as unknown as typeof FileReader

global.console = {
  ...console,
  log: vi.fn(),
  error: vi.fn(),
  warn: vi.fn(),
}
