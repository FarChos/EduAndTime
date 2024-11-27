/// <reference types="vue" />

declare namespace JSX {
  interface IntrinsicElements {
    [elem: string]: any; // Esto permite cualquier elemento JSX
  }
}
