<script lang="ts">
  import ListTodo from "@lucide/svelte/icons/list-todo";
  import Plus from "@lucide/svelte/icons/plus";

  import TaskPreviewList from "./TaskPreviewList.svelte";
  import TaskCreateDialog from "./TaskCreateDialog.svelte";

  let showCreateDialog = $state(false);
  let taskListVersion = $state(0);
  let pendingCount = $state<number | null>(null);
</script>

<section class="panel task-panel" aria-labelledby="tasks-title">
  <header class="task-header">
    <div>
      <div class="task-heading">
        <ListTodo size={21} strokeWidth={1.8} />
        <h2 id="tasks-title">Mis tareas</h2>
      </div>

      <p>Un paso a la vez.</p>
    </div>

    <span class="task-count">
      {pendingCount === null ? "Pendientes" : `${pendingCount} pendientes`}
    </span>
  </header>

  {#key taskListVersion}
    <TaskPreviewList
      onCountChange={(count: number | null) => {
        pendingCount = count;
      }}
    />
  {/key}

  <footer class="task-footer">
    <button
      type="button"
      class="add-task"
      onclick={() => (showCreateDialog = true)}
    >
      <Plus size={18} strokeWidth={1.8} />
      Nueva tarea
    </button>
  </footer>
</section>

{#if showCreateDialog}
  <TaskCreateDialog
    onClose={() => {
      showCreateDialog = false;
    }}
    onCreated={() => {
      pendingCount = null;
      taskListVersion += 1;
    }}
  />
{/if}

<style>
  .task-panel {
    min-width: 0;
  }

  .task-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 24px;
  }

  .task-heading {
    display: flex;
    align-items: center;
    gap: 10px;
    color: var(--color-primary);
  }

  .task-heading h2 {
    margin: 0;
    color: var(--color-text);
  }

  .task-header p {
    margin: 7px 0 0;
    color: var(--color-text-muted);
    font-size: 13px;
  }

  .task-count {
    padding: 7px 12px;
    border-radius: 999px;
    background: var(--color-surface-soft);
    color: var(--color-primary);
    font-size: 12px;
    font-weight: 600;
  }

  .task-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 12px;
    margin-top: 28px;
    padding-top: 18px;
    border-top: 1px solid var(--color-border);
  }

  .add-task {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 14px;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-control);
    background: var(--color-surface-soft);
    color: var(--color-primary);
    font-size: 13px;
    font-weight: 600;
  }

  .add-task:disabled {
    cursor: not-allowed;
    opacity: 0.75;
  }

</style>
