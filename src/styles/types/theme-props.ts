export type ThemeProviderProps = {
  children: React.ReactNode;
  defaultDarkMode?: boolean;
};

export type ThemeContextType = {
  isDarkMode: boolean;
  setIsDarkMode: (isDarkMode: boolean) => void;
  toggleTheme: () => void;
};