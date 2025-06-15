import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import pt from "./pt.json";
import es from "./es.json";
import en from "./en.json";
import { translations as dashboardTranslations } from "../pages/contents/Dashboard/translation/pt";

const resources = {
  pt: {
    translation: {
      ...pt.translation,
      ...dashboardTranslations
    }
  },
  es: {
    translation: {
      ...es.translation,
      ...dashboardTranslations
    }
  },
  en: {
    translation: {
      ...en.translation,
      ...dashboardTranslations
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