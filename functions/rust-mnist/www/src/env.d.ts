/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_SLS_FUNCTION_URL: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}