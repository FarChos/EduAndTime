import { defineStore } from 'pinia';
import { authInstance } from '../utils/axios';

// Store para recursos calificados
export const useCalificadosStore = defineStore('recursosCalificados', {
  state: () => ({
    Ids: [] as number[], // Inicializa como un array vacío
    calificaciones: [] as number[],
  }),
  actions: {
    setCalificados(data: { ids: number[]; calificaciones: number[] }) {
      this.Ids = data.ids;
      this.calificaciones = data.calificaciones;
    },
  },
});

// Store para mis recursos
export const useMisRecursosStore = defineStore('idesMisRecursos', {
  state: () => ({
    Ids: [] as number[], // Inicializa como un array vacío
  }),
  actions: {
    setMisRecursos(ids: number[]) {
      this.Ids = ids;
    },
  },
});

// Store para favoritos
export const useFavoritosStore = defineStore('idesFavoritos', {
  state: () => ({
    Ids: [] as number[], // Inicializa como un array vacío
  }),
  actions: {
    setFavoritos(ids: number[]) {
      this.Ids = ids;
    },
  },
});

export const cargarNuevoDatoUsuarioXRecurso = async (
  id: number,
  ideFavorito?: number | null,
  ideMiRecurso?: number | null,
  idCalificado?: number | null,
  calificacion?: number | null
) => {
  // Convertir undefined a null si es necesario
  if (ideFavorito === undefined) ideFavorito = null;
  if (ideMiRecurso === undefined) ideMiRecurso = null;
  if (idCalificado === undefined) idCalificado = null;
  if (calificacion === undefined) calificacion = null;

  try {
    const mutation = `
      mutation actualizarUsuarioXRecursos($id: Int!, $input: UsuarioXRecursoInput!) {
        actualizarUsuarioXRecursos(id: $id, input: $input) {
          Exito
        }
      }
    `;

    // Construcción dinámica del input
    const input: any = {};
    if (ideFavorito !== null) input.ideFavorito = ideFavorito;
    if (ideMiRecurso !== null) input.ideMiRecurso = ideMiRecurso;
    if (idCalificado !== null && calificacion !== null) {
      input.recursoCalificado = {
        idCalificado,
        calificacion,
      };
    }

    const variables = { id, input };

    const response = await authInstance.post('/query', {
      query: mutation,
      variables,
    });

    const resultado = response?.data?.data?.actualizarUsuarioXRecursos?.Exito;

    if (resultado) {
      console.log('Recurso subido correctamente.');
    } else {
      console.log('Hubo un error al subir el recurso.');
    }
  } catch (error) {
    console.error('Error al cargar el recurso:', error);
  }
};


// Función para cargar todos los datos desde el backend
export const datosUsuarioXrecurso = async (id: number) => {
  try {
    const query = `
      query tomarUsuarioXRecursos($id: Int) {
        tomarUsuarioXRecursos(id: $id) {
          idesFavoritos
          idesMisRecursos
          recursosCalificados {
            id
            calificacion
          }
        }
      }
    `;
    const variables = { id: id };
    const response = await authInstance.post('/query', { query, variables });
    const tomarUsuarioXRecursos = response?.data?.data?.tomarUsuarioXRecursos;

    if (!tomarUsuarioXRecursos) {
      throw new Error('No se encontraron datos en la respuesta del servidor.');
    }

    // Desestructurar datos de la respuesta
    const favoritos: number[] = tomarUsuarioXRecursos?.idesFavoritos || [];
    const misRecursos: number[] = tomarUsuarioXRecursos?.idesMisRecursos || [];
    const recursosCalificados = tomarUsuarioXRecursos?.recursosCalificados || [];

    const idsCalificados: number[] = recursosCalificados.map((item: any) => item.id);
    const calificaciones: number[] = recursosCalificados.map((item: any) => item.calificacion);

    // Actualizar cada store con los datos obtenidos
    const favoritosStore = useFavoritosStore();
    const misRecursosStore = useMisRecursosStore();
    const calificadosStore = useCalificadosStore();

    favoritosStore.setFavoritos(favoritos);
    misRecursosStore.setMisRecursos(misRecursos);
    calificadosStore.setCalificados({ ids: idsCalificados, calificaciones });

    console.log('Datos cargados exitosamente en los stores.');
  } catch (error) {
    console.error('Error al cargar los favoritos:', error);
  }
};


export const eliminarDatoUsuarioXRecurso = async (
  id: number,
  ideFavorito?: number,
  ideMiRecurso?: number,
  idCalificado?: number,
  calificacion?: number
) => {
  try {
    const mutation = `
      mutation eliminarUsuarioXRecurso($id: Int!, $input: UsuarioXRecursoInput!) {
        actualizarUsuarioXRecursos(id: $id, input: $input) {
          Exito
        }
      }
    `;

    // Construir el input dinámicamente con los datos disponibles
    const input: any = {};
    if (ideFavorito) input.ideFavorito = ideFavorito;
    if (ideMiRecurso) input.ideMiRecurso = ideMiRecurso;
    if (idCalificado && calificacion !== undefined) {
      input.recursoCalificado = {
        idCalificado: idCalificado,
        calificacion: calificacion,
      };
    }

    const variables = { id, input };

    // Realizar la petición
    const response = await authInstance.post('/query', {
      query: mutation,
      variables,
    });

    const resultado = response?.data?.data?.actualizarUsuarioXRecursos?.Exito;

    if (resultado) {
      console.log('Recurso eliminado correctamente.');
    } else {
      console.log('Hubo un error al eliminar el recurso.');
    }
  } catch (error) {
    console.error('Error al eliminar el recurso el recurso:', error);
  }
};