export interface FarmData {
  id: number;
  logo: string;
  language: string;
  company_id: number;
  company?: {
    id: number;
    company_name: string;
    location: string;
    farm_cnpj: string;
  };
  created_at: string;
  updated_at: string;
}

export interface BackendFarmData {
  ID: number;
  Logo: string;
  Language: string;
  CompanyID: number;
  Company?: {
    ID: number;
    CompanyName: string;
    Location: string;
    FarmCNPJ: string;
  };
  CreatedAt: string;
  UpdatedAt: string;
}

export interface UpdateFarmParams {
  logo: string;
  language?: string;
}
