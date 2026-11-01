<script lang="ts">
  const today = new Date();

  const todayLabel = new Intl.DateTimeFormat("es-UY", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(today);

  // Datos ficticios para la maqueta.
  const todayTasks = [
    {
      title: "Estudiar Bases de Datos",
      card: "Universidad",
      time: "10:00",
      priority: "Alta",
    },
    {
      title: "Revisar diseño del calendario",
      card: "Proyecto Luma",
      time: "16:00",
      priority: "Media",
    },
  ];

  const attention = [
    {
      title: "Preparar presentación",
      context: "Universidad",
      status: "Vencida",
    },
    {
      title: "Revisar apuntes",
      context: "Universidad",
      status: "Reprogramar",
    },
  ];

  const agenda = [
    {
      offset: 1,
      entries: [{ time: "09:00", title: "Actualizar currículum" }],
    },
    {
      offset: 2,
      entries: [{ time: "14:00", title: "Preparar parcial" }],
    },
    {
      offset: 3,
      entries: [{ time: "11:00", title: "Organizar apuntes" }],
    },
  ];

  function dayLabel(offset: number): string {
    if (offset === 0) return "Hoy";

    const day = new Date(
      today.getFullYear(),
      today.getMonth(),
      today.getDate() + offset,
    );

    return new Intl.DateTimeFormat("es-UY", {
      weekday: "short",
      day: "numeric",
      month: "short",
    }).format(day);
  }
</script>

<main class="content home-page">
  <header class="home-header">
    <div>
      <h1>Inicio</h1>
      <p class="header-subtitle">Un vistazo a lo que tenés por delante.</p>
    </div>

    <span class="today-date">{todayLabel}</span>
  </header>

  <section class="summary-stats" aria-label="Resumen de tareas de ejemplo">
    <article class="stat-card">
      <span class="stat-label">Pendientes para hoy</span>
      <strong>2</strong>
      <span class="stat-note">Con horario programado</span>
    </article>

    <article class="stat-card">
      <span class="stat-label">Completadas hoy</span>
      <strong>1</strong>
      <span class="stat-note">Actividades finalizadas</span>
    </article>

    <article class="stat-card">
      <span class="stat-label">Requieren atención</span>
      <strong>2</strong>
      <span class="stat-note"> Vencidas o por reprogramar </span>
    </article>
  </section>

  <section class="overview-grid" aria-label="Panel de resumen">
    <div class="overview-column">
      <article class="overview-panel">
        <div class="section-heading">
          <div>
            <h2>Tareas de hoy</h2>
            <p>Lo que tenés programado para este día</p>
          </div>

          <span class="count-label">2 pendientes</span>
        </div>

        <ul class="today-list">
          {#each todayTasks as task (task.title)}
            <li>
              <span class="task-check" aria-hidden="true"></span>

              <div class="task-info">
                <strong>{task.title}</strong>
                <small>{task.card} · {task.time}</small>
              </div>

              <span
                class:high={task.priority === "Alta"}
                class="priority-label"
              >
                {task.priority}
              </span>
            </li>
          {/each}
        </ul>
      </article>

      <article class="overview-panel">
        <div class="section-heading">
          <div>
            <h2>Requieren atención</h2>
            <p>Sin modificaciones automáticas</p>
          </div>

          <span class="count-label">2 tareas</span>
        </div>

        <ul class="attention-list">
          {#each attention as task (task.title)}
            <li>
              <div class="task-info">
                <strong>{task.title}</strong>
                <small>{task.context}</small>
              </div>

              <span
                class="attention-status"
                class:overdue={task.status === "Vencida"}
              >
                {task.status}
              </span>
            </li>
          {/each}
        </ul>
      </article>
    </div>

    <article class="overview-panel agenda-panel">
      <div class="section-heading">
        <div>
          <h2>Próximos días</h2>
          <p>Tu planificación a partir de mañana</p>
        </div>

        <span class="count-label">3 días</span>
      </div>

      <div class="agenda-list">
        {#each agenda as day (day.offset)}
          <section class="agenda-day">
            <h3>{dayLabel(day.offset)}</h3>

            {#each day.entries as event (event.title)}
              <div class="agenda-event">
                <time>{event.time}</time>
                <span>{event.title}</span>
              </div>
            {/each}
          </section>
        {/each}
      </div>
    </article>
  </section>

  <p class="demo-disclaimer">
    Maqueta visual: todas las tareas, cifras y fechas programadas son ejemplos.
    No hay conexión con tus datos reales.
  </p>
</main>

<style>
  .home-page {
    max-width: 1360px;
    min-width: 0;
  }

  .home-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    gap: 20px;
    margin-bottom: 26px;
  }

  .home-header h1 {
    margin: 0;
    font-size: clamp(28px, 3vw, 36px);
    letter-spacing: -1px;
  }

  .header-subtitle {
    margin: 8px 0 0;
    font-size: 14px;
    color: var(--color-text-muted);
  }

  .today-date {
    color: var(--color-text-muted);
    font-size: 13px;
    text-transform: capitalize;
    text-align: right;
  }

  .summary-stats {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 18px;
    margin-bottom: 22px;
  }

  .stat-card,
  .overview-panel {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-card);
    box-shadow: var(--shadow-card);
  }

  .stat-card {
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 22px;
  }

  .stat-label {
    font-size: 13px;
    color: var(--color-text-muted);
  }

  .stat-card strong {
    font-size: 36px;
    line-height: 1.1;
    font-variant-numeric: tabular-nums;
  }

  .stat-note {
    font-size: 12px;
    color: var(--color-text-muted);
  }

  .overview-grid {
    display: grid;
    grid-template-columns: minmax(0, 1.15fr) minmax(0, 0.85fr);
    gap: 20px;
    align-items: start;
  }

  .overview-column {
    display: grid;
    gap: 20px;
    min-width: 0;
  }

  .overview-panel {
    min-width: 0;
    padding: 22px;
  }

  .section-heading {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 18px;
  }

  .section-heading h2 {
    margin: 0;
    font-size: 19px;
  }

  .section-heading p {
    margin: 6px 0 0;
    font-size: 12px;
    color: var(--color-text-muted);
  }

  .count-label {
    white-space: nowrap;
    color: var(--color-primary);
    font-size: 12px;
  }

  .today-list,
  .attention-list {
    list-style: none;
    display: grid;
    gap: 9px;
    padding: 0;
    margin: 0;
  }

  .today-list li,
  .attention-list li {
    display: flex;
    align-items: center;
    gap: 11px;
    min-width: 0;
    padding: 13px;
    border-radius: var(--radius-control);
    background: var(--color-surface-soft);
  }

  .task-check {
    flex: 0 0 18px;
    height: 18px;
    border: 1.5px solid var(--color-text-muted);
    border-radius: 5px;
  }

  .task-info {
    min-width: 0;
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .task-info strong {
    font-size: 13px;
    line-height: 1.4;
    font-weight: 600;
    overflow-wrap: anywhere;
  }

  .task-info small {
    color: var(--color-text-muted);
    font-size: 11px;
  }

  .priority-label {
    flex-shrink: 0;
    font-size: 11px;
    padding: 5px 8px;
    border-left: 3px solid #c29240;
    border-radius: 5px;
    background: rgba(194, 146, 64, 0.16);
  }

  .priority-label.high {
    border-left-color: #d45d50;
    background: rgba(212, 93, 80, 0.16);
  }

  .attention-status {
    flex-shrink: 0;
    font-size: 11px;
    color: var(--color-text-muted);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-control);
    padding: 5px 7px;
  }

  .attention-status.overdue {
    color: #e79b8c;
    border-color: rgba(212, 93, 80, 0.55);
  }

  .agenda-list {
    display: grid;
    gap: 20px;
  }

  .agenda-day {
    border-top: 1px solid var(--color-border);
    padding-top: 13px;
  }

  .agenda-day h3 {
    margin: 0 0 10px;
    font-size: 13px;
    text-transform: capitalize;
    color: var(--color-text-muted);
  }

  .agenda-event {
    display: flex;
    align-items: baseline;
    gap: 11px;
    padding: 10px 12px;
    margin-top: 7px;
    border-left: 3px solid var(--color-primary);
    border-radius: var(--radius-control);
    background: var(--color-surface-soft);
  }

  .agenda-event time {
    flex-shrink: 0;
    color: var(--color-primary);
    font-size: 11px;
    font-variant-numeric: tabular-nums;
  }

  .agenda-event span {
    font-size: 12px;
    line-height: 1.4;
    overflow-wrap: anywhere;
  }

  .demo-disclaimer {
    margin: 18px 0 0;
    color: var(--color-text-muted);
    font-size: 12px;
    line-height: 1.5;
  }

  @media (max-width: 1050px) {
    .overview-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 650px) {
    .home-header {
      align-items: flex-start;
      flex-direction: column;
    }

    .today-date {
      text-align: left;
    }

    .summary-stats {
      grid-template-columns: 1fr;
      gap: 10px;
    }

    .stat-card {
      padding: 16px;
    }

    .overview-panel {
      padding: 17px;
    }
  }
</style>
