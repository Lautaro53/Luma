
export type TaskPriority = 'baja' | 'media' | 'alta';

export type TaskStatus = 'pendiente' | 'completada';

export type TaskSort = 'newest' | 'oldest' | 'priority';

export interface Task {
  id: number;
  title: string;
  description: string;
  priority: TaskPriority;
  status: TaskStatus;
  folder_id: number | null;
  created_at: string;
  updated_at: string;
  completed_at: string | null;
}

export interface Folder {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface TaskInput {
  title: string;
  description?: string;
  priority?: TaskPriority;
  folder_id?: number | null;
}

export interface FolderInput {
  name: string;
}

export interface TaskFilters {
  status?: TaskStatus;
  priority?: TaskPriority;
  folder_id?: number;
  without_folder?: boolean;
  search?: string;
  sort?: TaskSort;
}
