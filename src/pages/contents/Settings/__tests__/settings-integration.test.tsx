import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { message } from 'antd'
import { Settings } from '../presentation/settings'
import { useSelectedFarm } from '../../../../hooks/useSelectedFarm'
import { useFarm } from '../../../../hooks/useFarm'
import { UpdateFarmFactory, GetFarmFactory } from '../factories'

// Mock dos hooks
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

// Mock do i18n
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

describe('Settings Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    
    vi.mocked(useSelectedFarm).mockReturnValue({
      selectedFarm: null,
      setSelectedFarm: vi.fn(),
      clearSelectedFarm: vi.fn(),
      farmId: 1,
      farmName: 'Fazenda Teste',
      farmLogo: 'logo.png'
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

  describe('Fluxo completo de carregamento de dados', () => {
    it('deve carregar dados da fazenda e exibir corretamente', async () => {
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

      expect(screen.getByDisplayValue('Empresa Teste')).toBeInTheDocument()
      expect(screen.getByRole('img')).toHaveAttribute('src', 'data:image/jpeg;base64,test')
      expect(screen.getByText('Alterar Logo')).toBeInTheDocument()
    })

    it('deve tratar erro no carregamento de dados', async () => {
      mockGetFarmUseCase.get.mockRejectedValue(new Error('Erro na API'))

      renderSettings()

      await waitFor(() => {
        expect(mockGetFarmUseCase.get).toHaveBeenCalledWith(1)
      })

      expect(screen.getByDisplayValue('Fazenda Teste')).toBeInTheDocument()
    })
  })

  describe('Fluxo de upload de logo', () => {
    it('deve fazer upload de logo com sucesso', async () => {
      mockGetFarmUseCase.get.mockResolvedValue({
        data: { ...mockFarmData, logo: '' },
        success: true,
        message: 'Success',
        status: 200
      })

      mockUpdateFarmUseCase.update.mockResolvedValue({
        success: true,
        message: 'Logo atualizada com sucesso',
        data: { ...mockFarmData, logo: 'data:image/jpeg;base64,newlogo' }
      })

      mockGetFarmUseCase.get
        .mockResolvedValueOnce({
          data: { ...mockFarmData, logo: '' },
          success: true,
          message: 'Success',
          status: 200
        })
        .mockResolvedValueOnce({
          data: { ...mockFarmData, logo: 'data:image/jpeg;base64,newlogo' },
          success: true,
          message: 'Success',
          status: 200
        })

      renderSettings()

      await waitFor(() => {
        expect(screen.getByText('Adicionar Logo')).toBeInTheDocument()
      })

      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
      const uploadInput = screen.getByRole('button', { name: /adicionar logo/i })
      
      fireEvent.click(uploadInput)
      
      const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
      if (fileInput) {
        fireEvent.change(fileInput, { target: { files: [file] } })
      }

      await waitFor(() => {
        expect(mockUpdateFarmUseCase.update).toHaveBeenCalledWith(1, {
          logo: expect.stringContaining('data:image/jpeg;base64,')
        })
      })

      expect(message.success).toHaveBeenCalledWith('Logo atualizada com sucesso!')
    })

    it('deve tratar erro no upload de logo', async () => {
      mockGetFarmUseCase.get.mockResolvedValue({
        data: { ...mockFarmData, logo: '' },
        success: true,
        message: 'Success',
        status: 200
      })

      mockUpdateFarmUseCase.update.mockRejectedValue(new Error('Erro no upload'))

      renderSettings()

      await waitFor(() => {
        expect(screen.getByText('Adicionar Logo')).toBeInTheDocument()
      })

      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
      const uploadInput = screen.getByRole('button', { name: /adicionar logo/i })
      
      fireEvent.click(uploadInput)
      
      const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
      if (fileInput) {
        fireEvent.change(fileInput, { target: { files: [file] } })
      }

      await waitFor(() => {
        expect(message.error).toHaveBeenCalledWith('Erro ao atualizar logo')
      })
    })
  })

  describe('Estados de erro e loading', () => {
    it('deve exibir mensagem quando farmId é null', () => {
      vi.mocked(useSelectedFarm).mockReturnValue({
        selectedFarm: null,
        setSelectedFarm: vi.fn(),
        clearSelectedFarm: vi.fn(),
        farmId: null,
        farmName: null,
        farmLogo: ''
      })

      renderSettings()

      expect(screen.getByText('Nenhuma fazenda selecionada')).toBeInTheDocument()
    })

    it('deve exibir dados de fallback quando não consegue carregar dados', async () => {
      mockGetFarmUseCase.get.mockRejectedValue(new Error('Erro na API'))

      renderSettings()

      await waitFor(() => {
        expect(screen.getByDisplayValue('Fazenda Teste')).toBeInTheDocument()
      })
    })
  })
})
