import React, { ReactElement } from 'react'
import { render, RenderOptions, screen, fireEvent, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import { ThemeProvider } from '../styles/context/theme/theme-provider'
import { AntConfigWrapper } from '../styles/config/ant-design-config-wrapper'
import { vi } from 'vitest'

const mockTheme = {
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 6,
  },
}

const AllTheProviders = ({ children }: { children: React.ReactNode }) => {
  return (
    <BrowserRouter>
      <ThemeProvider>
        <AntConfigWrapper>
          <ConfigProvider theme={mockTheme}>
            {children}
          </ConfigProvider>
        </AntConfigWrapper>
      </ThemeProvider>
    </BrowserRouter>
  )
}

const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) => render(ui, { wrapper: AllTheProviders, ...options })

export * from '@testing-library/react'

export { customRender as render }

export const createMockUser = () => ({
  id: 1,
  name: 'Test User',
  email: 'test@example.com',
  farm_id: 1,
})

export const createMockFarm = () => ({
  id: 1,
  name: 'Test Farm',
  domain: 'test-farm',
})

export const createMockAnimal = () => ({
  id: 1,
  name: 'Test Animal',
  ear_tag: '001',
  breed: 'Holandesa',
  birth_date: '2020-01-01',
  farm_id: 1,
})

export const createMockReproduction = () => ({
  id: 1,
  animal_id: 1,
  animal_name: 'Test Animal',
  ear_tag: '001',
  current_phase: 1,
  insemination_date: '2024-01-01',
  pregnancy_date: '2024-01-15',
  expected_birth_date: '2024-10-15',
  veterinary_confirmation: true,
  farm_id: 1,
})

export const mockUseAuth = () => ({
  user: createMockUser(),
  farm: createMockFarm(),
  login: vi.fn(),
  logout: vi.fn(),
  isLoading: false,
  error: null,
})

export const mockUseFarm = () => ({
  farm: createMockFarm(),
  setFarm: vi.fn(),
})

export const mockUseReproduction = () => ({
  getReproductionsByFarm: vi.fn().mockResolvedValue([createMockReproduction()]),
  createReproduction: vi.fn().mockResolvedValue(createMockReproduction()),
  updateReproduction: vi.fn().mockResolvedValue(true),
  deleteReproduction: vi.fn().mockResolvedValue(true),
  loading: false,
  error: null,
})

export const mockTranslation = {
  t: (key: string) => key,
  i18n: {
    changeLanguage: vi.fn(),
    language: 'pt',
  },
}
