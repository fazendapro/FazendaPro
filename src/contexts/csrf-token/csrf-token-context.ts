import { createContext, useContext } from 'react'

import { CsrfTokenContextProps } from './csrf-token-context-props'

export const CsrfTokenContext = createContext<CsrfTokenContextProps>({
  csrfToken: undefined
})

export function useCsrfTokenContext(): CsrfTokenContextProps {
  const context = useContext<CsrfTokenContextProps>(CsrfTokenContext)
  const { csrfToken } = context
  return {
    csrfToken
  }
}
