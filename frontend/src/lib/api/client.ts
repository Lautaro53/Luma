
import type {
  Task,
  TaskInput,
  TaskStatus,
  TaskFilters,
  Folder,
  FolderInput,
} from './types';

// Error HTTP reconocido por la aplicación.
export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// Cliente HTTP centralizado.
async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const headers = new Headers(options.headers);

  if (options.body != null && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  let response: Response;

  try {
    response = await fetch(`/api${path}`, {
      ...options,
      headers,
    });
  } catch (error) {
    if (
      error instanceof DOMException &&
      error.name === 'AbortError'
    ) {
      throw error;
    }

    throw new ApiError(
      'No se pudo conectar con Luma. Verificá que el backend esté iniciado.',
      0,
    );
  }

  if (!response.ok) {
    let message = `Error HTTP ${response.status}`;

    try {
      const data: unknown = await response.json();

      if (
        typeof data === 'object' &&
        data !== null &&
        'error' in data &&
        typeof data.error === 'string'
      ) {
        message = data.error;
      }
    } catch {
      // Se conserva el mensaje HTTP predeterminado.
    }

    throw new ApiError(message, response.status);
  }

  // DELETE devuelve 204 sin cuerpo.
  if (response.status === 204) {
    return undefined as T;
  }

  try {
    return (await response.json()) as T;
  } catch {
    throw new ApiError(
      'El servidor devolvió una respuesta JSON inválida.',
      response.status,
    );
  }
}

// Genera los parámetros de consulta de tareas.
function buildTaskQuery(filters: TaskFilters): string {
  const params = new URLSearchParams();

  if (filters.status) {
    params.set('status', filters.status);
  }

  if (filters.priority) {
    params.set('priority', filters.priority);
  }

  if (filters.folder_id !== undefined) {
    params.set('folder_id', String(filters.folder_id));
  }

  if (filters.without_folder !== undefined) {
    params.set(
      'without_folder',
      String(filters.without_folder),
    );
  }

  if (filters.search) {
    params.set('search', filters.search);
  }

  if (filters.sort) {
    params.set('sort', filters.sort);
  }

  const query = params.toString();

  return query ? `?${query}` : '';
}

// Operaciones sobre tareas.
export const tasksApi = {
  list(filters: TaskFilters = {}): Promise<Task[]> {
    return request<Task[]>(
      `/tasks${buildTaskQuery(filters)}`,
    );
  },

  create(input: TaskInput): Promise<Task> {
    return request<Task>('/tasks', {
      method: 'POST',
      body: JSON.stringify(input),
    });
  },

  update(id: number, input: TaskInput): Promise<Task> {
    return request<Task>(`/tasks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    });
  },

  setStatus(id: number, status: TaskStatus): Promise<Task> {
    return request<Task>(`/tasks/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    });
  },

  remove(id: number): Promise<void> {
    return request<void>(`/tasks/${id}`, {
      method: 'DELETE',
    });
  },
};

// Operaciones sobre carpetas.
export const foldersApi = {
  list(): Promise<Folder[]> {
    return request<Folder[]>('/folders');
  },

  create(input: FolderInput): Promise<Folder> {
    return request<Folder>('/folders', {
      method: 'POST',
      body: JSON.stringify(input),
    });
  },

  rename(id: number, input: FolderInput): Promise<Folder> {
    return request<Folder>(`/folders/${id}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    });
  },

  remove(id: number): Promise<void> {
    return request<void>(`/folders/${id}`, {
      method: 'DELETE',
    });
  },
};
