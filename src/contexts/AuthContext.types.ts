export interface DecodedToken {
  exp: number;
  iat: number;
  sub: string;
  email: string;
}

export interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: DecodedToken | null;
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => Promise<void>;
  updateToken: (newToken: string) => void;
  token: string | null;
}
