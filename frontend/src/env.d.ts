/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** "true" = mode demo statis tanpa backend (data dilayani di browser oleh utils/demoApi.ts) */
  readonly VITE_DEMO_MODE?: string
  /** Base URL backend Golang; kosong = same-origin (proxy Vite / mode single-binary) */
  readonly VITE_API_BASE?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<object, object, any>
  export default component
}
