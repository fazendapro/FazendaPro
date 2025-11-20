import { describe, it, expect } from 'vitest'
import { renderHook } from '@testing-library/react'
import { useCsrfTokenContext, CsrfTokenContext } from '../csrf-token-context'
import { ReactNode } from 'react'

describe('useCsrfTokenContext', () => {
  it('deve retornar csrfToken do contexto', () => {
    const mockCsrfToken = 'test-csrf-token'
    
    const wrapper = ({ children }: { children: ReactNode }) => (
      <CsrfTokenContext.Provider value={{ csrfToken: mockCsrfToken }}>
        {children}
      </CsrfTokenContext.Provider>
    )

    const { result } = renderHook(() => useCsrfTokenContext(), { wrapper })

    expect(result.current.csrfToken).toBe(mockCsrfToken)
  })

  it('deve retornar undefined quando csrfToken não está definido', () => {
    const wrapper = ({ children }: { children: ReactNode }) => (
      <CsrfTokenContext.Provider value={{ csrfToken: undefined }}>
        {children}
      </CsrfTokenContext.Provider>
    )

    const { result } = renderHook(() => useCsrfTokenContext(), { wrapper })

    expect(result.current.csrfToken).toBeUndefined()
  })
})

