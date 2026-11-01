<script lang="ts">
  import { onMount } from "svelte";
  import { tasksApi } from "../api/client";
  import type { Task } from "../api/types";

  interface Props {
    onCountChange?: (count: number | null) => void;
  }

  let { onCountChange }: Props = $props();

  let tasks = $state<Task[]>([]);
  let loading = $state(true);
  let errorMessage = $state<string | null>(null);

  let mounted = false;

  async function loadTasks(): Promise<void> {
    loading = true;
    errorMessage = null;

    try {
      const result = await tasksApi.list({
        status: "pendiente",
        sort: "newest",
      });

      if (mounted) {
        tasks = result;
        onCountChange?.(result.length);
      }
    } catch (error) {
      if (mounted) {
        onCountChange?.(null);

        errorMessage =
          error instanceof Error
            ? error.message
            : "No se pudieron cargar las tareas.";
      }
    } finally {
      if (mounted) {
        loading = false;
      }
    }
  }

  onMount(() => {
    mounted = true;
    void loadTasks();

    return () => {
      mounted = false;
    };
  });
</script>

<div class="task-preview-content">
  {#if loading}
    <p class="task-state" role="status">Cargando tareas...</p>
  {:else if errorMessage}
    <div class="task-state" role="alert">
      <p>{errorMessage}</p>

      <button
        type="button"
        class="retry-button"
        onclick={() => void loadTasks()}
      >
        Reintentar
      </button>
    </div>
  {:else if tasks.length === 0}
    <div class="task-state">
      <p>No tenés tareas pendientes.</p>
      <span>Un buen momento para empezar algo nuevo.</span>
    </div>
  {:else}
    <ul class="task-list" aria-label="Tareas pendientes">
      {#each tasks.slice(0, 3) as task (task.id)}
        <li class="task-item">
          <span class="task-circle" aria-hidden="true"></span>

          <div class="task-information">
            <span class="task-title">{task.title}</span>

            {#if task.description}
              <span class="task-description">
                {task.description}
              </span>
            {/if}
          </div>

          <span class="task-priority">
            {task.priority}
          </span>
        </li>
      {/each}
    </ul>

    {#if tasks.length > 3}
      <p class="task-more">
        +{tasks.length - 3} tareas pendientes
      </p>
    {/if}
  {/if}
</div>

<style>
  .task-preview-content {
    width: 100%;
    color: inherit;
  }

  .task-list {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .task-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 0;
    border-bottom: 1px solid color-mix(in srgb, currentColor 10%, transparent);
  }

  .task-item:last-child {
    border-bottom: 0;
  }

  .task-circle {
    width: 19px;
    height: 19px;
    flex-shrink: 0;
    border: 1.5px solid currentColor;
    border-radius: 50%;
    opacity: 0.4;
  }

  .task-information {
    display: flex;
    flex-direction: column;
    gap: 4px;
    flex: 1;
    min-width: 0;
  }

  .task-title {
    font-size: 0.9rem;
    font-weight: 500;
    overflow-wrap: anywhere;
  }

  .task-description {
    font-size: 0.77rem;
    opacity: 0.65;
    overflow-wrap: anywhere;
  }

  .task-priority {
    flex-shrink: 0;
    font-size: 0.7rem;
    text-transform: capitalize;
    padding: 5px 9px;
    border-radius: 999px;
    background: color-mix(in srgb, currentColor 8%, transparent);
  }

  .task-state {
    padding: 25px 12px;
    text-align: center;
    font-size: 0.85rem;
    opacity: 0.8;
  }

  .task-state p {
    margin: 0 0 6px;
  }

  .task-state span {
    font-size: 0.78rem;
    opacity: 0.7;
  }

  .retry-button {
    margin-top: 12px;
    padding: 7px 14px;
    border: 1px solid currentColor;
    border-radius: 8px;
    background: transparent;
    color: inherit;
    cursor: pointer;
  }

  .task-more {
    margin: 12px 0 0;
    font-size: 0.76rem;
    opacity: 0.7;
  }
</style>
