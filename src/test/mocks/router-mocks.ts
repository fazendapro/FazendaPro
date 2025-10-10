import { vi } from 'vitest'
import React from 'react'

export const mockRouter = {
  useNavigate: () => vi.fn(),
  useLocation: () => ({
    pathname: '/',
    search: '',
    hash: '',
    state: null,
  }),
  useParams: () => ({}),
  useSearchParams: () => [new URLSearchParams(), vi.fn()],
  BrowserRouter: ({ children }: { children: React.ReactNode }) => children,
  Routes: ({ children }: { children: React.ReactNode }) => children,
  Route: ({ children }: { children: React.ReactNode }) => children,
  Navigate: () => null,
  Link: ({ children, to, ...props }: { children: React.ReactNode; to: string; [key: string]: unknown }) => 
    React.createElement('a', { href: to, ...props }, children),
  NavLink: ({ children, to, ...props }: { children: React.ReactNode; to: string; [key: string]: unknown }) => 
    React.createElement('a', { href: to, ...props }, children),
}

export const createMockNavigate = () => {
  const navigate = vi.fn()
  return {
    navigate,
    mockNavigate: navigate,
  }
}

export const createMockLocation = (pathname = '/', search = '', hash = '', state = null) => ({
  pathname,
  search,
  hash,
  state,
})

export const createMockParams = (params: Record<string, string> = {}) => params

export const createMockSearchParams = (params: Record<string, string> = {}) => {
  const searchParams = new URLSearchParams(params)
  const setSearchParams = vi.fn()
  return [searchParams, setSearchParams]
}
