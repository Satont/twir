import type { WritableComputedRef } from 'vue';
import { ComposerTranslation, createI18n, useI18n } from 'vue-i18n';

import type ILandingLocale from '@/types/landingLocaleInterface.js';
import type { Locale } from '@/types/locale';

type Languages = { name: string; locale: Locale }[];
type LocaleType = keyof LocaleTypes;

interface LocaleTypes {
  landing: ILandingLocale;
  app: ReturnType<() => typeof import('@/locales/app/en.json')>;
}

export const languages: Languages = [
  { name: 'English', locale: 'en' },
  { name: 'Русский', locale: 'ru' },
];

export const locales: Locale[] = ['en', 'ru'];
export const defaultLocale = 'en';

export const useTranslation = <L extends LocaleType>() => {
  const { t } = useI18n();

  return t as ComposerTranslation<{ en: LocaleTypes[L] }, Locale>;
};

export async function loadLocaleMessages<L extends LocaleType>(
  localeType: L,
  locale: Locale,
): Promise<LocaleTypes[L]> {
  return (await import(`../locales/${localeType}/${locale}.ts`)).default;
}

export async function setupI18n(locale: Locale, localeType: LocaleType) {
  const i18n = createI18n({ locale, legacy: false });

  (i18n.global.locale as WritableComputedRef<Locale>).value = locale;
  const messages = await loadLocaleMessages(localeType, locale);
  i18n.global.setLocaleMessage<any>(locale, messages);

  return i18n;
}
