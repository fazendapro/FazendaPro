import { describe, it, expect, vi, beforeEach } from 'vitest'
import { UpdateFarmFactory, GetFarmFactory } from '../factories/usecases/update-farm-factory'
import { RemoteUpdateFarm, RemoteGetFarm } from '../data/usecases/remote-update-farm'

vi.mock('../data/usecases/remote-update-farm', () => ({
  RemoteUpdateFarm: vi.fn(),
  RemoteGetFarm: vi.fn(),
}))

describe('UpdateFarmFactory', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve criar instância de RemoteUpdateFarm', () => {
    const factory = UpdateFarmFactory.create()
    
    expect(RemoteUpdateFarm).toHaveBeenCalled()
    expect(factory).toBeInstanceOf(RemoteUpdateFarm)
  })

  it('deve retornar nova instância a cada chamada', () => {
    const factory1 = UpdateFarmFactory.create()
    const factory2 = UpdateFarmFactory.create()
    
    expect(factory1).not.toBe(factory2)
    expect(RemoteUpdateFarm).toHaveBeenCalledTimes(2)
  })
})

describe('GetFarmFactory', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve criar instância de RemoteGetFarm', () => {
    const factory = GetFarmFactory.create()
    
    expect(RemoteGetFarm).toHaveBeenCalled()
    expect(factory).toBeInstanceOf(RemoteGetFarm)
  })

  it('deve retornar nova instância a cada chamada', () => {
    const factory1 = GetFarmFactory.create()
    const factory2 = GetFarmFactory.create()
    
    expect(factory1).not.toBe(factory2)
    expect(RemoteGetFarm).toHaveBeenCalledTimes(2)
  })
})
