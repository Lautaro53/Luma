
<script lang="ts">
  import { onMount } from 'svelte';
  import { foldersApi, tasksApi } from '../api/client';
  import type { Folder, TaskPriority } from '../api/types';

  interface Props {
    onCreated: () => void;
    onClose: () => void;
  }

  let { onCreated, onClose }: Props = $props();

  let dialog: HTMLDialogElement;

  let title = $state('');
  let description = $state('');
  let priority = $state<TaskPriority>('media');
  let folderId = $state('');

  let folders = $state<Folder[]>([]);
  let loadingFolders = $state(true);
  let folderError = $state('');
  let saving = $state(false);
  let errorMessage = $state('');

  async function loadFolders(): Promise<void> {
    loadingFolders = true;
    folderError = '';

    try {
      folders = await foldersApi.list();
    } catch {
      folderError = 'No se pudieron cargar las carpetas.';
    } finally {
      loadingFolders = false;
    }
  }

  async function handleSubmit(
    event: SubmitEvent,
  ): Promise<void> {
    event.preventDefault();

    if (saving) return;

    const cleanTitle = title.trim();

    if (!cleanTitle || cleanTitle.length > 200) {
      errorMessage = 'Ingresá un título de hasta 200 caracteres.';
      return;
    }

    saving = true;
    errorMessage = '';

    try {
      await tasksApi.create({
        title: cleanTitle,
        description: description.trim(),
        priority,
        folder_id: folderId === ''
          ? null
          : Number(folderId),
      });

      onCreated();
      dialog.close();
    } catch (error) {
      errorMessage = error instanceof Error
        ? error.message
        : 'No se pudo crear la tarea.';
    } finally {
      saving = false;
    }
  }

  onMount(() => {
    dialog.showModal();
    void loadFolders();
  });
</script>

<dialog
  bind:this={dialog}
  class="task-dialog"
  aria-labelledby="create-task-title"
  onclose={onClose}
  oncancel={(event) => {
    if (saving) event.preventDefault();
  }}
>
  <form onsubmit={handleSubmit}>
    <header class="dialog-header">
      <h2 id="create-task-title">Nueva tarea</h2>
      <p>Organizá lo que querés hacer.</p>
    </header>

    <div class="fields">
      <label for="task-title">Título *</label>
      <input
        id="task-title"
        type="text"
        bind:value={title}
        maxlength="200"
        placeholder="Ej. Estudiar infraestructura"
        required
        disabled={saving}
      />

      <label for="task-description">Descripción</label>
      <textarea
        id="task-description"
        bind:value={description}
        rows="3"
        placeholder="Detalles opcionales"
        disabled={saving}
      ></textarea>

      <div class="field-grid">
        <div class="field-group">
          <label for="task-priority">Prioridad</label>
          <select
            id="task-priority"
            bind:value={priority}
            disabled={saving}
          >
            <option value="baja">Baja</option>
            <option value="media">Media</option>
            <option value="alta">Alta</option>
          </select>
        </div>

        <div class="field-group">
          <label for="task-folder">Carpeta</label>
          <select
            id="task-folder"
            bind:value={folderId}
            disabled={saving || loadingFolders}
          >
            <option value="">Sin carpeta</option>
            {#each folders as folder (folder.id)}
              <option value={String(folder.id)}>
                {folder.name}
              </option>
            {/each}
          </select>
        </div>
      </div>

      {#if folderError}
        <div class="folder-warning" role="status">
          <span>{folderError}</span>
          <button
            type="button"
            onclick={() => void loadFolders()}
            disabled={saving || loadingFolders}
          >
            Reintentar
          </button>
        </div>
      {/if}

      {#if errorMessage}
        <p class="form-error" role="alert">
          {errorMessage}
        </p>
      {/if}
    </div>

    <footer class="dialog-actions">
      <button
        type="button"
        class="cancel-button"
        onclick={() => dialog.close()}
        disabled={saving}
      >
        Cancelar
      </button>

      <button
        type="submit"
        class="save-button"
        disabled={saving || !title.trim()}
      >
        {saving ? 'Guardando...' : 'Crear tarea'}
      </button>
    </footer>
  </form>
</dialog>

<style>
  .task-dialog {
    box-sizing: border-box;
    width: min(460px, calc(100vw - 32px));
    max-height: calc(100dvh - 32px);
    padding: 0;
    overflow: auto;
    border: 1px solid
      color-mix(in srgb, currentColor 15%, transparent);
    border-radius: 20px;
    background: var(--surface, var(--card-bg, Canvas));
    color: var(--text-primary, CanvasText);
    box-shadow: 0 20px 65px rgba(0, 0, 0, 0.18);
    font: inherit;
  }

  :global(html[data-theme='forest']) .task-dialog,
  :global(html[data-theme='midnight']) .task-dialog {
    color-scheme: dark;
  }

  .task-dialog::backdrop {
    background: rgba(20, 30, 25, 0.5);
  }

  .task-dialog form {
    padding: 26px;
  }

  .dialog-header h2 {
    margin: 0;
    font-size: 1.3rem;
  }

  .dialog-header p {
    margin: 8px 0 24px;
    font-size: 0.88rem;
    opacity: 0.7;
  }

  .fields {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  label {
    font-size: 0.83rem;
    font-weight: 600;
  }

  input,
  textarea,
  select {
    box-sizing: border-box;
    width: 100%;
    padding: 11px 12px;
    border: 1px solid
      color-mix(in srgb, currentColor 20%, transparent);
    border-radius: 10px;
    background: transparent;
    color: inherit;
    font: inherit;
    font-size: 0.88rem;
  }

  textarea {
    resize: vertical;
  }

  input:focus-visible,
  textarea:focus-visible,
  select:focus-visible,
  button:focus-visible {
    outline: 2px solid currentColor;
    outline-offset: 2px;
  }

  .field-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px;
    margin-top: 6px;
  }

  .field-group {
    display: flex;
    flex-direction: column;
    gap: 10px;
    min-width: 0;
  }

  .folder-warning,
  .form-error {
    font-size: 0.8rem;
    line-height: 1.5;
  }

  .folder-warning button {
    margin-left: 8px;
    border: 0;
    background: transparent;
    color: inherit;
    text-decoration: underline;
    cursor: pointer;
  }

  .form-error {
    margin: 0;
    color: #ca6b59;
  }

  .dialog-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 26px;
  }

  .dialog-actions button {
    padding: 11px 16px;
    border-radius: 10px;
    border: 1px solid
      color-mix(in srgb, currentColor 20%, transparent);
    color: inherit;
    font: inherit;
    font-size: 0.85rem;
    cursor: pointer;
  }

  .cancel-button {
    background: transparent;
  }

  .save-button {
    background: color-mix(
      in srgb,
      currentColor 13%,
      transparent
    );
    font-weight: 600 !important;
  }

  .dialog-actions button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  @media (max-width: 420px) {
    .task-dialog form {
      padding: 20px;
    }

    .field-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
