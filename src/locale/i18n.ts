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
import { translations as salesTranslations } from "./sales/pt";

const resources = {
  pt: {
    translation: {
      ...pt.translation,
      ...dashboardTranslationsPt,
      ...animalTableTranslationsPt,
      ...settingsTranslationsPt,
      ...salesTranslations
    }
  },
  es: {
    translation: {
      ...es.translation,
      ...dashboardTranslationsEs,
      ...animalTableTranslationsEs,
      ...settingsTranslationsEs,
      ...salesTranslations
    }
  },
  en: {
    translation: {
      ...en.translation,
      ...dashboardTranslationsEn,
      ...animalTableTranslationsEn,
      ...settingsTranslationsEn,
      ...salesTranslations
    }
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