import { useState, useEffect } from "react";
import { ThemeProviderProps } from "../types/theme-props";
import { ThemeContext } from "./theme-context";

export const ThemeProvider = ({
  children,
  defaultDarkMode = false
}: ThemeProviderProps) => {
  const [isDarkMode, setIsDarkMode] = useState(() => {
    const savedTheme = localStorage.getItem('darkMode');
    return savedTheme !== null ? JSON.parse(savedTheme) : defaultDarkMode;
  });

  const toggleTheme = () => {
    setIsDarkMode(!isDarkMode);
  };

  useEffect(() => {
    localStorage.setItem('darkMode', JSON.stringify(isDarkMode));
    if (isDarkMode) {
      document.body.classList.add('dark-mode');
    } else {
      document.body.classList.remove('dark-mode');
    }
  }, [isDarkMode]);

  return (
    <ThemeContext.Provider value={{ isDarkMode, setIsDarkMode, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
};