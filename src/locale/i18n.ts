import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import pt from "./pt.json";
import es from "./es.json";
import en from "./en.json";
import { translations as dashboardTranslationsPt } from "../pages/contents/Dashboard/translation/pt";
import { translations as dashboardTranslationsEn } from "../pages/contents/Dashboard/translation/en";
import { translations as dashboardTranslationsEs } from "../pages/contents/Dashboard/translation/es";
import { translations as animalTableTranslationsPt } from "../pages/contents/AnimalTable/translation/pt";
import { translations as animalTableTranslationsEn } from "../pages/contents/AnimalTable/translation/en";
import { translations as animalTableTranslationsEs } from "../pages/contents/AnimalTable/translation/es";
import { settingsTranslations as settingsTranslationsPt } from "../pages/contents/Settings/translation/pt";
import { settingsTranslations as settingsTranslationsEn } from "../pages/contents/Settings/translation/en";
import { settingsTranslations as settingsTranslationsEs } from "../pages/contents/Settings/translation/es";
import { translations as salesTranslationsPt } from "./sales/pt";
import { translations as salesTranslationsEn } from "./sales/en";
import { translations as salesTranslationsEs } from "./sales/es";

const mergeTranslations = (base: Record<string, unknown>, ...overrides: Record<string, unknown>[]): Record<string, unknown> => {
  const merged = { ...base };
  overrides.forEach(override => {
    Object.keys(override).forEach(key => {
      if (key === 'common' && typeof override[key] === 'object' && override[key] !== null && 
          typeof merged[key] === 'object' && merged[key] !== null) {
        merged[key] = { ...(override[key] as Record<string, unknown>), ...(merged[key] as Record<string, unknown>) };
      } else if (typeof override[key] === 'object' && override[key] !== null && !Array.isArray(override[key]) && 
          typeof merged[key] === 'object' && merged[key] !== null && !Array.isArray(merged[key])) {
        merged[key] = mergeTranslations(merged[key] as Record<string, unknown>, override[key] as Record<string, unknown>);
      } else {
        merged[key] = override[key];
      }
    });
  });
  return merged;
};

const resources = {
  pt: {
    translation: mergeTranslations(
      pt.translation,
      dashboardTranslationsPt,
      animalTableTranslationsPt,
      settingsTranslationsPt,
      salesTranslationsPt
    )
  },
  es: {
    translation: mergeTranslations(
      es.translation,
      dashboardTranslationsEs,
      animalTableTranslationsEs,
      settingsTranslationsEs,
      salesTranslationsEs
    )
  },
  en: {
    translation: mergeTranslations(
      en.translation,
      dashboardTranslationsEn,
      animalTableTranslationsEn,
      settingsTranslationsEn,
      salesTranslationsEn
    )
  }
};

i18n
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'en',
    lng: "pt",
    interpolation: {
      escapeValue: false
    },
    react: {
      useSuspense: false
    },
    debug: false
  });

export default i18n;