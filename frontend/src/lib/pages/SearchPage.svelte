<script lang="ts">
  import Search from "@lucide/svelte/icons/search";
  import Folder from "@lucide/svelte/icons/folder";
  import SlidersHorizontal from "@lucide/svelte/icons/sliders-horizontal";

  type Priority = "Alta" | "Media" | "Baja";

  type SearchTask = {
    title: string;
    description: string;
    priority: Priority;
    completed: boolean;
  };

  type TaskGroup = {
    card: string;
    tasks: SearchTask[];
  };

  // Datos exclusivamente visuales.
  const matchingCards = [
    {
      name: "Universidad",
      description: "Materias, parciales y entregas.",
      total: 10,
    },
    {
      name: "Proyecto Luma",
      description: "Diseño y desarrollo de la aplicación.",
      total: 5,
    },
  ];

  const taskGroups: TaskGroup[] = [
    {
      card: "Universidad",
      tasks: [
        {
          title: "Estudiar Bases de Datos",
          description: "Preparar los ejercicios del parcial.",
          priority: "Alta",
          completed: false,
        },
        {
          title: "Repasar normalización",
          description: "Repasar las formas normales.",
          priority: "Media",
          completed: false,
        },
      ],
    },
    {
      card: "Proyecto Luma",
      tasks: [
        {
          title: "Definir reglas funcionales",
          description: "Consolidar la nueva definición de Luma.",
          priority: "Alta",
          completed: true,
        },
        {
          title: "Diseñar pantalla de tarjetas",
          description: "Revisar tamaños y desplazamiento.",
          priority: "Media",
          completed: false,
        },
      ],
    },
  ];
</script>

<main class="content search-page">
  <header class="page-intro">
    <p class="eyebrow">TU ESPACIO</p>
    <h1>Buscar</h1>
    <p>Encontrá tus tarjetas y tareas en un solo lugar.</p>
  </header>

  <section class="search-panel">
    <div class="search-field">
      <Search size={20} strokeWidth={1.8} />
      <span>Buscar tarjetas o tareas...</span>
      <span class="shortcut">Vista previa</span>
    </div>

    <div class="filters">
      <div class="filter-heading">
        <SlidersHorizontal size={17} />
        <span>Filtros</span>
      </div>

      <div class="filter-list">
        <div class="filter-chip">
          Estado: Todos
          <span>⌄</span>
        </div>

        <div class="filter-chip">
          Prioridad: Todas
          <span>⌄</span>
        </div>

        <div class="filter-chip">
          Tarjeta: Todas
          <span>⌄</span>
        </div>
      </div>
    </div>
  </section>

  <section class="results">
    <div class="results-heading">
      <div>
        <h2>Resultados de ejemplo</h2>
        <p>Tarjetas y tareas organizadas por contexto.</p>
      </div>

      <span class="results-count"> 2 tarjetas · 4 tareas </span>
    </div>

    <div class="results-layout">
      <section class="results-section">
        <div class="section-title">
          <h3>Tarjetas</h3>
          <span>{matchingCards.length}</span>
        </div>

        <div class="card-results">
          {#each matchingCards as card (card.name)}
            <article class="matching-card">
              <div class="card-icon">
                <Folder size={19} strokeWidth={1.8} />
              </div>

              <div class="matching-card-info">
                <strong>{card.name}</strong>
                <p>{card.description}</p>
                <small>{card.total} tareas</small>
              </div>
            </article>
          {/each}
        </div>
      </section>

      <section class="results-section">
        <div class="section-title">
          <h3>Tareas</h3>
          <span>4</span>
        </div>

        <div class="task-groups">
          {#each taskGroups as group (group.card)}
            <div class="task-group">
              <div class="group-heading">
                <Folder size={15} strokeWidth={1.8} />
                <span>{group.card}</span>
              </div>

              <ul class="task-results">
                {#each group.tasks as task (task.title)}
                  <li class:completed={task.completed}>
                    <span
                      class="task-check"
                      class:checked={task.completed}
                      aria-hidden="true"
                    >
                      {task.completed ? "✓" : ""}
                    </span>

                    <div class="task-info">
                      <strong>{task.title}</strong>
                      <p>{task.description}</p>
                    </div>

                    <span
                      class={`priority priority-${task.priority.toLowerCase()}`}
                    >
                      {task.priority}
                    </span>
                  </li>
                {/each}
              </ul>
            </div>
          {/each}
        </div>
      </section>
    </div>
  </section>

  <p class="demo-note">
    Maqueta visual con resultados ficticios. Los filtros y la búsqueda todavía
    no tienen funcionalidad.
  </p>
</main>

<style>
  .search-page {
    display: block;

    width: 100%;
    max-width: 1400px;
    height: 100%;

    min-width: 0;
    min-height: 0;

    overflow-x: hidden;
    overflow-y: auto;
  }

  .page-intro {
    flex-shrink: 0;
    margin-bottom: 16px;
  }

  .page-intro .eyebrow {
    margin: 0 0 8px;
    color: var(--color-primary);
  }

  .page-intro h1 {
    margin: 0;
    font-size: clamp(28px, 3vw, 36px);
  }

  .page-intro > p:last-child {
    margin: 8px 0 0;
    font-size: 14px;
    color: var(--color-text-muted);
  }

  .search-panel,
  .results-section {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-card);
    box-shadow: var(--shadow-card);
  }

  .search-panel {
    flex-shrink: 0;
    padding: 16px 18px;
    margin-bottom: 16px;
  }
  .search-field {
    display: flex;
    align-items: center;
    gap: 13px;
    min-height: 54px;
    padding: 0 16px;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-control);
    background: var(--color-surface-soft);
    color: var(--color-text-muted);
    font-size: 14px;
  }

  .search-field > span:nth-child(2) {
    flex: 1;
  }

  .shortcut {
    font-size: 11px;
    color: var(--color-text-muted);
  }

  .filters {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    align-items: center;
    margin-top: 18px;
  }

  .filter-heading {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--color-text-muted);
    font-size: 13px;
  }

  .filter-list {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
  }

  .filter-chip {
    display: flex;
    align-items: center;
    gap: 15px;
    padding: 9px 12px;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-control);
    color: var(--color-text);
    font-size: 12px;
  }

  .filter-chip span {
    color: var(--color-text-muted);
  }

  .results {
    display: block;
    width: 100%;
    min-width: 0;
  }

  .results-heading {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    margin-bottom: 14px;
  }

  .results-heading h2 {
    margin: 0;
    font-size: 21px;
  }

  .results-heading p {
    margin: 7px 0 0;
    font-size: 12px;
    color: var(--color-text-muted);
  }

  .results-count {
    font-size: 12px;
    color: var(--color-text-muted);
  }

  .results-layout {
    display: grid;
    grid-template-columns: minmax(0, 0.8fr) minmax(0, 1.2fr);
    align-items: start;
    gap: 20px;

    width: 100%;
    min-width: 0;
  }

  .results-section {
    min-width: 0;
    padding: 20px;
  }

  .section-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
  }

  .section-title h3 {
    margin: 0;
    font-size: 17px;
  }

  .section-title span {
    color: var(--color-primary);
    font-size: 12px;
  }

  .card-results,
  .task-groups {
    display: grid;
    gap: 12px;
  }

  .matching-card {
    display: flex;
    gap: 13px;
    padding: 15px;
    background: var(--color-surface-soft);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-control);
  }

  .card-icon {
    display: grid;
    place-items: center;
    flex: 0 0 38px;
    height: 38px;
    border-radius: var(--radius-control);
    background: var(--color-surface);
    color: var(--color-primary);
  }

  .matching-card-info {
    min-width: 0;
  }

  .matching-card-info strong {
    font-size: 14px;
  }

  .matching-card-info p {
    margin: 6px 0;
    font-size: 12px;
    color: var(--color-text-muted);
  }

  .matching-card-info small {
    font-size: 11px;
    color: var(--color-primary);
  }

  .task-group {
    display: grid;
    gap: 9px;
  }

  .group-heading {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 0;
    color: var(--color-text-muted);
    font-size: 12px;
  }

  .task-results {
    display: grid;
    gap: 9px;
    padding: 0;
    margin: 0 0 10px;
    list-style: none;
  }

  .task-results li {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
    padding: 13px;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-control);
    background: var(--color-surface-soft);
  }

  .task-check {
    display: grid;
    place-items: center;
    flex: 0 0 18px;
    height: 18px;
    border: 1.5px solid var(--color-text-muted);
    border-radius: 5px;
    font-size: 12px;
  }

  .task-check.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: var(--color-background);
  }

  .task-info {
    flex: 1;
    min-width: 0;
  }

  .task-info strong {
    font-size: 13px;
    font-weight: 600;
    overflow-wrap: anywhere;
  }

  .task-info p {
    margin: 5px 0 0;
    font-size: 11px;
    color: var(--color-text-muted);
  }

  .completed .task-info strong {
    text-decoration: line-through;
    color: var(--color-text-muted);
  }

  .priority {
    flex-shrink: 0;
    padding: 5px 8px;
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

  .demo-note {
    flex-shrink: 0;
    margin: 10px 0 0;
    color: var(--color-text-muted);
    font-size: 12px;
  }

  @media (max-width: 1050px) {
    .results-layout {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 600px) {
    .search-panel,
    .results-section {
      padding: 15px;
    }

    .results-heading {
      align-items: flex-start;
      flex-direction: column;
    }

    .shortcut {
      display: none;
    }
  }
</style>
