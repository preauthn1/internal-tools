import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import en from './i18n/en.json';
import jp from './i18n/jp.json';
import de from './i18n/de.json';

const resources = {
  en: {
    translation: en,
  },
  'jp':{
    translation: jp,
  },
  'de':{
    translation: de,
  }
};

i18n.use(initReactI18next).init({
  resources,
  lng: 'en',            
  fallbackLng: 'en',     
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;