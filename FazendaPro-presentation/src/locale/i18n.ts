import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import pt from "./pt.json";
import es from "./es.json";
import en from "./en.json";

i18n.use(initReactI18next).init({
  resources: {
    pt: {
      ...pt,
    },
    es: {
      ...es,
    },
    en: {
      ...en,
    },
  },
  lng: "pt",
});

export default i18n;