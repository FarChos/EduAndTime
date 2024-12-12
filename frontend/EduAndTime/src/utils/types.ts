export enum Categoria {
  Sociedad= 'sociedad',
  Geografia= '  geografia',
  Tecnologia= '  tecnologia',
  Ciencia= '  ciencia',
  Economia= '  economia',
  Bienestar= '  bienestar',
  Politica= '  politica',
  Arte= '  arte',
  Filosofia= '  filosofia',
  Exotica= '  exotica',
}

export enum Formato {
  PDF= 'pdf',
  EPUB= 'epub',
  MOVI= 'mobi',
  TXT= 'txt',
  AZW= 'azw',
  AZW3= 'azw3',
  FB2= 'fb2',
  DJVU= 'djvu',
  DOCX= 'docx',
  ODT= 'odt',
}

export type GraphQLResponse<T> = {
  data?: T;
  errors?: Array<{
    message: string;
  }>;
};


export type RecursoMuestra = {
  Id: number;
  Titulo: string;
  Autor: string;
  Categoria: Categoria; // ENUM'S
  Formato: Formato; // ENUM'S
  Archivo: string;
  Etiquetas: string[];
  Calificacion: number;
}

export type Recurso = {
  Id: number;
  Titulo: string;
  Autor: string;
  Categoria: Categoria;
  IdUsuario: number;
  Formato: Formato;
  Descripcion: string;
  Archivo: string;
  FechaOrigen: string;
  Etiquetas: string[];
  Calificacion: number;
  NumDescargas: number;
}