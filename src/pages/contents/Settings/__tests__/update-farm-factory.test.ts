import { describe, it, expect, vi, beforeEach } from 'vitest'
import { UpdateFarmFactory, GetFarmFactory } from '../factories/usecases/update-farm-factory'

const mockRemoteUpdateFarm = vi.fn()
const mockRemoteGetFarm = vi.fn()

vi.mock('../data/usecases/remote-update-farm', () => ({
  RemoteUpdateFarm: class {
    constructor() {
      mockRemoteUpdateFarm()
    }
  },
  RemoteGetFarm: class {
    constructor() {
      mockRemoteGetFarm()
    }
  },
}))

describe('UpdateFarmFactory', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve criar instância de RemoteUpdateFarm', () => {
    const factory = UpdateFarmFactory.create()
    
    expect(mockRemoteUpdateFarm).toHaveBeenCalled()
    expect(factory).toBeDefined()
  })

  it('deve retornar nova instância a cada chamada', () => {
    const factory1 = UpdateFarmFactory.create()
    const factory2 = UpdateFarmFactory.create()
    
    expect(factory1).not.toBe(factory2)
    expect(mockRemoteUpdateFarm).toHaveBeenCalledTimes(2)
  })
})

describe('GetFarmFactory', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve criar instância de RemoteGetFarm', () => {
    const factory = GetFarmFactory.create()
    
    expect(mockRemoteGetFarm).toHaveBeenCalled()
    expect(factory).toBeDefined()
  })

  it('deve retornar nova instância a cada chamada', () => {
    const factory1 = GetFarmFactory.create()
    const factory2 = GetFarmFactory.create()
    
    expect(factory1).not.toBe(factory2)
    expect(mockRemoteGetFarm).toHaveBeenCalledTimes(2)
  })
})
