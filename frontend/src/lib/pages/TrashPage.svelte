<script lang="ts">
  import Trash2 from "@lucide/svelte/icons/trash-2";
  import RotateCcw from "@lucide/svelte/icons/rotate-ccw";
  import Clock3 from "@lucide/svelte/icons/clock-3";

  type Priority = "Alta" | "Media" | "Baja";

  type DeletedTask = {
    id: number;
    title: string;
    context: string;
    priority: Priority;
    deletedDaysAgo: number;
  };

  // Datos ficticios para la maqueta.
  const deletedTasks: DeletedTask[] = [
    {
      id: 1,
      title: "Revisar documentación antigua",
      context: "Proyecto Luma",
      priority: "Media",
      deletedDaysAgo: 3,
    },
    {
      id: 2,
      title: "Preparar informe anterior",
      context: "Universidad",
      priority: "Alta",
      deletedDaysAgo: 14,
    },
    {
      id: 3,
      title: "Organizar archivos personales",
      context: "Personal",
      priority: "Baja",
      deletedDaysAgo: 7,
    },
  ];

  function remainingDays(task: DeletedTask): number {
    return Math.max(0, 15 - task.deletedDaysAgo);
  }

  function dayLabel(days: number): string {
    return `${days} ${days === 1 ? "día" : "días"}`;
  }
</script>

<main class="content trash-page">
  <header class="page-intro">
    <div>
      <h1>Papelera</h1>
      <p>Recuperá tareas eliminadas antes de que desaparezcan.</p>
    </div>
  </header>

  <section class="trash-summary">
    <div class="summary-icon">
      <Trash2 size={24} strokeWidth={1.8} />
    </div>

    <div class="summary-text">
      <h2>Conservación de tareas</h2>
      <p>
        Las tareas eliminadas permanecen en la Papelera durante 15 días. Después
        se eliminan definitivamente.
      </p>
    </div>

    <div class="summary-count">
      <strong>{deletedTasks.length}</strong>
      <span>tareas</span>
    </div>
  </section>

  <section class="trash-panel">
    <header class="panel-heading">
      <div>
        <h2>Tareas eliminadas</h2>
        <p>Podrás restaurarlas o eliminarlas definitivamente.</p>
      </div>

      <button
        type="button"
        class="empty-trash"
        disabled
        title="Disponible en una etapa posterior"
      >
        <Trash2 size={16} strokeWidth={1.8} />
        Vaciar papelera
      </button>
    </header>

    <ul class="deleted-list">
      {#each deletedTasks as task (task.id)}
        <li class="deleted-task">
          <div class="task-main">
            <div class="task-icon">
              <Trash2 size={19} strokeWidth={1.8} />
            </div>

            <div class="task-info">
              <div class="task-title-row">
                <h3>{task.title}</h3>

                <span
                  class={`priority priority-${task.priority.toLowerCase()}`}
                >
                  {task.priority}
                </span>
              </div>

              <p>Tarjeta original: {task.context}</p>

              <div class="retention-info">
                <Clock3 size={14} strokeWidth={1.8} />

                <span>
                  Eliminada hace {dayLabel(task.deletedDaysAgo)}
                </span>

                <span class="separator">·</span>

                <span class:expiring={remainingDays(task) <= 2}>
                  Se elimina en {dayLabel(remainingDays(task))}
                </span>
              </div>
            </div>
          </div>

          <div class="task-actions">
            <button
              type="button"
              class="restore-action"
              disabled
              title="Disponible en una etapa posterior"
            >
              <RotateCcw size={15} strokeWidth={1.8} />
              Restaurar
            </button>

            <button
              type="button"
              class="delete-action"
              disabled
              title="Disponible en una etapa posterior"
            >
              Eliminar
            </button>
          </div>
        </li>
      {/each}
    </ul>
  </section>

  <p class="demo-note">
    Los elementos son ficticios. Ninguna acción modifica o elimina información
    de SQLite.
  </p>
</main>

<style>
  .trash-page {
    max-width: 1360px;
    min-width: 0;
  }

  .page-intro {
    margin-bottom: 16px;
  }

  .page-intro h1 {
    margin: 0;
    font-size: clamp(28px, 3vw, 36px);
  }

  .page-intro p:last-child {
    margin: 8px 0 0;
    color: var(--color-text-muted);
    font-size: 14px;
  }

  .demo-note {
    color: var(--color-text-muted);
    font-size: 12px;
  }

  .trash-summary {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 18px 20px;
    margin-bottom: 16px;
  }

  .summary-icon {
    flex-shrink: 0;
    display: grid;
    place-items: center;
    width: 54px;
    height: 54px;
    border-radius: var(--radius-control);
    background: var(--color-surface-soft);
    color: var(--color-primary);
  }

  .summary-text {
    flex: 1;
    min-width: 0;
  }

  .summary-text h2 {
    margin: 0 0 7px;
    font-size: 18px;
  }

  .summary-text p {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 13px;
    line-height: 1.6;
  }

  .summary-count {
    display: flex;
    flex-direction: column;
    align-items: center;
    min-width: 80px;
    gap: 3px;
  }

  .summary-count strong {
    font-size: 30px;
    line-height: 1;
    color: var(--color-primary);
  }

  .summary-count span {
    font-size: 12px;
    color: var(--color-text-muted);
  }

  .trash-panel {
    padding: 20px;
  }

  .panel-heading {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 20px;
    margin-bottom: 16px;
  }

  .panel-heading h2 {
    margin: 0;
    font-size: 20px;
  }

  .panel-heading p {
    margin: 7px 0 0;
    color: var(--color-text-muted);
    font-size: 12px;
  }

  .empty-trash,
  .restore-action,
  .delete-action {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    gap: 8px;
    padding: 10px 13px;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-control);
    background: var(--color-surface-soft);
    color: var(--color-text);
    font-size: 12px;
  }

  /* Acciones solo ilustrativas */
  button:disabled {
    opacity: 0.75;
    cursor: not-allowed;
  }

  .empty-trash,
  .delete-action {
    color: #d98b80;
  }

  .deleted-list {
    display: grid;
    gap: 12px;
    padding: 0;
    margin: 0;
    list-style: none;
  }

  .deleted-task {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    min-width: 0;
    padding: 18px;
    background: var(--color-surface-soft);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-control);
  }

  .task-main {
    display: flex;
    align-items: flex-start;
    gap: 13px;
    flex: 1;
    min-width: 0;
  }

  .task-icon {
    flex-shrink: 0;
    display: grid;
    place-items: center;
    width: 37px;
    height: 37px;
    color: var(--color-text-muted);
    border-radius: var(--radius-control);
    background: var(--color-surface);
  }

  .task-info {
    flex: 1;
    min-width: 0;
  }

  .task-title-row {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
  }

  .task-title-row h3 {
    margin: 0;
    font-size: 14px;
    line-height: 1.4;
    overflow-wrap: anywhere;
  }

  .task-info > p {
    margin: 7px 0;
    font-size: 12px;
    color: var(--color-text-muted);
  }

  .priority {
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 11px;
  }

  .priority-alta {
    background: rgba(212, 93, 80, 0.16);
    border-left: 3px solid #d45d50;
  }

  .priority-media {
    background: rgba(194, 146, 64, 0.16);
    border-left: 3px solid #c29240;
  }

  .priority-baja {
    background: rgba(89, 154, 113, 0.16);
    border-left: 3px solid #599a71;
  }

  .retention-info {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 7px;
    margin-top: 11px;
    font-size: 11px;
    color: var(--color-text-muted);
  }

  .retention-info .expiring {
    color: #d98b80;
    font-weight: 600;
  }

  .task-actions {
    display: flex;
    flex-shrink: 0;
    gap: 8px;
  }

  .restore-action {
    color: var(--color-primary);
  }

  .demo-note {
    margin-top: 16px;
  }

  @media (max-width: 1050px) {
    .deleted-task {
      flex-direction: column;
      align-items: stretch;
    }

    .task-actions {
      justify-content: flex-end;
    }
  }

  @media (max-width: 650px) {
    .page-intro {
      align-items: flex-start;
      flex-direction: column;
    }

    .trash-summary,
    .trash-panel {
      padding: 16px;
    }

    .summary-count {
      min-width: 50px;
    }

    .panel-heading {
      flex-direction: column;
      align-items: flex-start;
    }

    .task-actions {
      flex-wrap: wrap;
    }
  }
</style>
