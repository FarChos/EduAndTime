// * Función auxiliar para separar las etiquetas
export function formateadorEtiquetas(etiquetaCompleta: string): string[] {
  // Usamos una expresión regular para encontrar las etiquetas que comienzan con '#' y terminan con un espacio
  const regex = /#\S+/g; // Encuentra cualquier palabra que comience con '#' y no tenga espacios
  return etiquetaCompleta.match(regex) || []; // Devuelve el array de etiquetas o un array vacío si no hay coincidencias
}

export function formatoLibro(nombreDocumento: string): string {
  const regex = /\.\S+/; // Expresión regular que busca el punto seguido de caracteres no espacios
  const resultado = regex.exec(nombreDocumento); // Usamos exec() para obtener las coincidencias
  
  return resultado ? resultado[0].slice(1) : ""; // Si hay coincidencia, devuelve la extensión sin el punto '.'
}

