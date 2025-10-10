import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import pt from "./pt.json";
import es from "./es.json";
import en from "./en.json";
import { translations as dashboardTranslations } from "../pages/contents/Dashboard/translation/pt";
import { animalTableTranslations } from "../pages/contents/AnimalTable/translation";
import { settingsTranslations } from "../pages/contents/Settings/translation";

const resources = {
  pt: {
    translation: {
      ...pt.translation,
      ...dashboardTranslations,
      ...animalTableTranslations,
      ...settingsTranslations
    }
  },
  es: {
    translation: {
      ...es.translation,
      ...dashboardTranslations,
      ...animalTableTranslations,
      ...settingsTranslations
    }
  },
  en: {
    translation: {
      ...en.translation,
      ...dashboardTranslations,
      ...animalTableTranslations,
      ...settingsTranslations
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
    }
  });

export default i18n;