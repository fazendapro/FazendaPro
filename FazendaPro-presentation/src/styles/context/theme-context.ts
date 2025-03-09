import { createContext } from "react";
import { ThemeContextType } from "../types/theme-props";

export const ThemeContext = createContext<ThemeContextType>({
  isDarkMode: false,
  setIsDarkMode: () => {},
  toggleTheme: () => {},
});