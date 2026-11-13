"use client";
import * as React from "react";
import { DEFAULT_UI_LOCALE, parseUiLocale, UI_LOCALE_KEY, type UiLocale } from "./ui-locale";
import { getLocalStorageValue, setLocalStorageValue } from "./local-storage.client";
const listeners = new Set<() => void>();
let current: UiLocale | undefined;
function snapshot(): UiLocale { return current ?? parseUiLocale(getLocalStorageValue(UI_LOCALE_KEY)); }
function serverSnapshot(): UiLocale { return DEFAULT_UI_LOCALE; }
function subscribe(listener: () => void) {
  listeners.add(listener);
  const onStorage = (event: StorageEvent) => {
    if (event.key === UI_LOCALE_KEY || event.key === null) { current = undefined; listener(); }
  };
  window.addEventListener("storage", onStorage);
  return () => { listeners.delete(listener); window.removeEventListener("storage", onStorage); };
}
export function useUiLocale(): UiLocale {
  return React.useSyncExternalStore(subscribe, snapshot, serverSnapshot);
}
export function setUiLocale(locale: UiLocale) {
  current = locale;
  setLocalStorageValue(UI_LOCALE_KEY, locale);
  for (const listener of listeners) listener();
}
export function UiLocaleDocument() {
  const locale = useUiLocale();
  React.useEffect(() => { document.documentElement.lang = locale; }, [locale]);
  return null;
}
