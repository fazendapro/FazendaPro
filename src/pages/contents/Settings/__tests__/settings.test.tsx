import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { Settings } from '../presentation/settings'
import { useSelectedFarm } from '../../../../hooks/useSelectedFarm'
import { useFarm } from '../../../../hooks/useFarm'
import { UpdateFarmFactory, GetFarmFactory } from '../factories'

vi.mock('../../../../hooks/useSelectedFarm')
vi.mock('../../../../hooks/useFarm')
vi.mock('../factories')
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd')
  return {
    ...actual,
    message: {
      success: vi.fn(),
      error: vi.fn(),
    },
  }
})

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'title': 'Configurações da Fazenda',
        'basicInfo': 'Informações Básicas',
        'farmName': 'Nome da Fazenda',
        'farmNamePlaceholder': 'Nome da fazenda',
        'farmNameDisabled': 'O nome da fazenda é definido pela empresa e não pode ser alterado aqui.',
        'farmLogo': 'Logo da Fazenda',
        'changeLogo': 'Alterar Logo',
        'addLogo': 'Adicionar Logo',
        'logoFormats': 'Formatos aceitos: JPG, PNG, GIF (máx. 2MB)',
        'noFarmSelected': 'Nenhuma fazenda selecionada',
        'farmUpdatedSuccessfully': 'Fazenda atualizada com sucesso',
        'errorUpdatingFarm': 'Erro ao atualizar fazenda',
        'errorUpdatingLogo': 'Erro ao atualizar logo',
        'logoUpdatedSuccessfully': 'Logo atualizada com sucesso!',
      }
      return translations[key] || key
    },
  }),
}))

const mockFarm = {
  id: 1,
  name: 'Fazenda Teste',
  location: 'Rua Teste, 123',
  created_at: '2021-01-01T00:00:00Z',
  updated_at: '2021-01-01T00:00:00Z'
}

const mockFarmData = {
  id: 1,
  logo: 'data:image/jpeg;base64,test',
  company_id: 1,
  company: {
    id: 1,
    company_name: 'Empresa Teste',
    location: 'Rua Teste, 123',
    farm_cnpj: '12345678000199'
  },
  created_at: '2021-01-01T00:00:00Z',
  updated_at: '2021-01-01T00:00:00Z'
}

const mockUpdateFarmUseCase = {
  update: vi.fn()
}

const mockGetFarmUseCase = {
  get: vi.fn()
}

const renderSettings = () => {
  return render(
    <BrowserRouter>
      <Settings />
    </BrowserRouter>
  )
}

describe('Settings Component', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    
    vi.mocked(useSelectedFarm).mockReturnValue({
      selectedFarm: null,
      setSelectedFarm: vi.fn(),
      clearSelectedFarm: vi.fn(),
      farmId: 1,
      farmName: 'Fazenda Teste'
    })
    
    vi.mocked(useFarm).mockReturnValue({
      farm: mockFarm,
      loading: false,
      error: null,
      refetch: vi.fn()
    })
    
    vi.mocked(UpdateFarmFactory.create).mockReturnValue(mockUpdateFarmUseCase)
    vi.mocked(GetFarmFactory.create).mockReturnValue(mockGetFarmUseCase)
  })

  it('deve renderizar o título da página', () => {
    renderSettings()
    
    expect(screen.getByText('Configurações da Fazenda')).toBeInTheDocument()
  })

  it('deve renderizar as seções de informações básicas e logo', () => {
    renderSettings()
    
    expect(screen.getByText('Informações Básicas')).toBeInTheDocument()
    expect(screen.getByText('Logo da Fazenda')).toBeInTheDocument()
  })

  it('deve mostrar mensagem quando nenhuma fazenda está selecionada', () => {
    vi.mocked(useSelectedFarm).mockReturnValue({
      selectedFarm: null,
      setSelectedFarm: vi.fn(),
      clearSelectedFarm: vi.fn(),
      farmId: null,
      farmName: null
    })
    
    renderSettings()
    
    expect(screen.getByText('Nenhuma fazenda selecionada')).toBeInTheDocument()
  })

  it('deve carregar dados da fazenda ao montar o componente', async () => {
    mockGetFarmUseCase.get.mockResolvedValue({
      data: mockFarmData,
      success: true,
      message: 'Success',
      status: 200
    })
    
    renderSettings()
    
    await waitFor(() => {
      expect(mockGetFarmUseCase.get).toHaveBeenCalledWith(1)
    })
  })

  it('deve exibir o nome da fazenda desabilitado', async () => {
    mockGetFarmUseCase.get.mockResolvedValue({
      data: mockFarmData,
      success: true,
      message: 'Success',
      status: 200
    })
    
    renderSettings()
    
    await waitFor(() => {
      const nameInput = screen.getByDisplayValue('Empresa Teste')
      expect(nameInput).toBeDisabled()
    })
  })

  it('deve exibir o logo da fazenda quando disponível', async () => {
    mockGetFarmUseCase.get.mockResolvedValue({
      data: mockFarmData,
      success: true,
      message: 'Success',
      status: 200
    })
    
    renderSettings()
    
    await waitFor(() => {
      const avatar = screen.getByRole('img')
      expect(avatar).toHaveAttribute('src', 'data:image/jpeg;base64,test')
    })
  })

  it('deve mostrar botão "Adicionar Logo" quando não há logo', async () => {
    const farmDataWithoutLogo = { ...mockFarmData, logo: '' }
    mockGetFarmUseCase.get.mockResolvedValue({
      data: farmDataWithoutLogo,
      success: true,
      message: 'Success',
      status: 200
    })
    
    renderSettings()
    
    await waitFor(() => {
      expect(screen.getByText('Adicionar Logo')).toBeInTheDocument()
    })
  })

  it('deve mostrar botão "Alterar Logo" quando há logo', async () => {
    mockGetFarmUseCase.get.mockResolvedValue({
      data: mockFarmData,
      success: true,
      message: 'Success',
      status: 200
    })
    
    renderSettings()
    
    await waitFor(() => {
      expect(screen.getByText('Alterar Logo')).toBeInTheDocument()
    })
  })

  it('deve exibir informações sobre formatos aceitos', () => {
    renderSettings()
    
    expect(screen.getByText('Formatos aceitos: JPG, PNG, GIF (máx. 2MB)')).toBeInTheDocument()
  })

  it('deve exibir texto explicativo sobre o nome da fazenda', () => {
    renderSettings()
    
    expect(screen.getByText('O nome da fazenda é definido pela empresa e não pode ser alterado aqui.')).toBeInTheDocument()
  })
})
