import { describe, it, expect } from 'vitest'
import { FarmData, BackendFarmData, UpdateFarmParams } from '../types/farm-types'

describe('Farm Types', () => {
  describe('FarmData', () => {
    it('deve ter a estrutura correta', () => {
      const farmData: FarmData = {
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

      expect(farmData.id).toBe(1)
      expect(farmData.logo).toBe('data:image/jpeg;base64,test')
      expect(farmData.company_id).toBe(1)
      expect(farmData.company).toBeDefined()
      expect(farmData.company?.company_name).toBe('Empresa Teste')
      expect(farmData.created_at).toBe('2021-01-01T00:00:00Z')
      expect(farmData.updated_at).toBe('2021-01-01T00:00:00Z')
    })

    it('deve permitir company opcional', () => {
      const farmData: FarmData = {
        id: 1,
        logo: '',
        company_id: 1,
        created_at: '2021-01-01T00:00:00Z',
        updated_at: '2021-01-01T00:00:00Z'
      }

      expect(farmData.company).toBeUndefined()
    })
  })

  describe('BackendFarmData', () => {
    it('deve ter a estrutura correta com campos em maiúscula', () => {
      const backendFarmData: BackendFarmData = {
        ID: 1,
        Logo: 'data:image/jpeg;base64,test',
        CompanyID: 1,
        Company: {
          ID: 1,
          CompanyName: 'Empresa Teste',
          Location: 'Rua Teste, 123',
          FarmCNPJ: '12345678000199'
        },
        CreatedAt: '2021-01-01T00:00:00Z',
        UpdatedAt: '2021-01-01T00:00:00Z'
      }

      expect(backendFarmData.ID).toBe(1)
      expect(backendFarmData.Logo).toBe('data:image/jpeg;base64,test')
      expect(backendFarmData.CompanyID).toBe(1)
      expect(backendFarmData.Company).toBeDefined()
      expect(backendFarmData.Company?.CompanyName).toBe('Empresa Teste')
      expect(backendFarmData.CreatedAt).toBe('2021-01-01T00:00:00Z')
      expect(backendFarmData.UpdatedAt).toBe('2021-01-01T00:00:00Z')
    })

    it('deve permitir Company opcional', () => {
      const backendFarmData: BackendFarmData = {
        ID: 1,
        Logo: '',
        CompanyID: 1,
        CreatedAt: '2021-01-01T00:00:00Z',
        UpdatedAt: '2021-01-01T00:00:00Z'
      }

      expect(backendFarmData.Company).toBeUndefined()
    })
  })

  describe('UpdateFarmParams', () => {
    it('deve ter apenas o campo logo', () => {
      const updateParams: UpdateFarmParams = {
        logo: 'data:image/jpeg;base64,test'
      }

      expect(updateParams.logo).toBe('data:image/jpeg;base64,test')
    })

    it('deve permitir logo vazio', () => {
      const updateParams: UpdateFarmParams = {
        logo: ''
      }

      expect(updateParams.logo).toBe('')
    })
  })
})
