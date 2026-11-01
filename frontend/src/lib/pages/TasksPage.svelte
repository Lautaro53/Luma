<script lang="ts">
    import Folder from "@lucide/svelte/icons/folder";

    type Priority = "Alta" | "Media" | "Baja";

    type Task = {
        title: string;
        priority: Priority;
        done: boolean;
    };

    type ContextCard = {
        id: number;
        name: string;
        description: string;
        tasks: Task[];
    };

    // Datos ficticios para probar el diseño.
    const cards: ContextCard[] = [
        {
            id: 1,
            name: "Universidad",
            description: "Materias, parciales y entregas.",
            tasks: [
                {
                    title: "Preparar parcial de Bases de Datos",
                    priority: "Alta",
                    done: false,
                },
                {
                    title: "Repasar normalización",
                    priority: "Media",
                    done: false,
                },
                {
                    title: "Resolver ejercicios de grafos",
                    priority: "Alta",
                    done: false,
                },
                {
                    title: "Revisar apuntes de clase",
                    priority: "Media",
                    done: false,
                },
                { title: "Estudiar árboles", priority: "Baja", done: false },
                {
                    title: "Preparar presentación",
                    priority: "Alta",
                    done: false,
                },
                { title: "Leer bibliografía", priority: "Baja", done: false },
                {
                    title: "Revisar fechas de entregas",
                    priority: "Media",
                    done: false,
                },
                { title: "Preparar resumen", priority: "Baja", done: false },
                { title: "Organizar apuntes", priority: "Baja", done: true },
            ],
        },
        {
            id: 2,
            name: "Proyecto Luma",
            description: "Diseño y desarrollo de la aplicación.",
            tasks: [
                {
                    title: "Diseñar pantalla de tarjetas",
                    priority: "Alta",
                    done: false,
                },
                {
                    title: "Preparar calendario semanal",
                    priority: "Media",
                    done: false,
                },
                { title: "Revisar navegación", priority: "Media", done: false },
                {
                    title: "Ajustar los temas visuales",
                    priority: "Baja",
                    done: false,
                },
                {
                    title: "Definir reglas funcionales",
                    priority: "Alta",
                    done: true,
                },
            ],
        },
        {
            id: 3,
            name: "Personal",
            description: "Pendientes y actividades cotidianas.",
            tasks: [
                {
                    title: "Actualizar currículum",
                    priority: "Media",
                    done: false,
                },
                { title: "Ordenar documentos", priority: "Baja", done: false },
            ],
        },
        {
            id: 4,
            name: "Trabajo",
            description: "Objetivos y actividades profesionales.",
            tasks: [
                {
                    title: "Organizar tareas de la semana",
                    priority: "Alta",
                    done: false,
                },
                {
                    title: "Preparar informe de avance",
                    priority: "Media",
                    done: false,
                },
            ],
        },
    ];

    let boardElement: HTMLElement | undefined;

    function moveBoard(direction: -1 | 1) {
        if (!boardElement) return;

        boardElement.scrollBy({
            left: direction * Math.max(280, boardElement.clientWidth * 0.75),
            behavior: "smooth",
        });
    }
</script>

<main class="content tasks-page">
    <header class="page-header">
        <div>
            <h1>Mis tareas</h1>
            <p>Organizá tus tareas en tarjetas de contexto.</p>
        </div>

        <div class="board-actions">
            <button
                type="button"
                class="board-arrow"
                aria-label="Ver tarjetas anteriores"
                onclick={() => moveBoard(-1)}
            >
                ‹
            </button>

            <button
                type="button"
                class="board-arrow"
                aria-label="Ver tarjetas siguientes"
                onclick={() => moveBoard(1)}
            >
                ›
            </button>
        </div>
    </header>

    <section
        class="board"
        bind:this={boardElement}
        aria-label="Tablero de tarjetas de contexto"
    >
        {#each cards as card (card.id)}
            {@const completed = card.tasks.filter((task) => task.done).length}
            {@const progress = card.tasks.length
                ? Math.round((completed * 100) / card.tasks.length)
                : 0}

            <article class="context-card">
                <header class="card-header">
                    <div class="folder-icon">
                        <Folder size={20} strokeWidth={1.8} />
                    </div>

                    <div class="card-heading">
                        <h2>{card.name}</h2>
                        <p>{card.description}</p>
                    </div>

                    <span class="card-more" aria-hidden="true">···</span>
                </header>

                <div
                    class="progress"
                    aria-label={`${completed} de ${card.tasks.length} tareas completadas`}
                >
                    <div class="progress-label">
                        <span
                            >{completed} de {card.tasks.length} completadas</span
                        >
                        <span>{progress}%</span>
                    </div>

                    <div class="progress-track">
                        <div
                            class="progress-fill"
                            style:width={`${progress}%`}
                        ></div>
                    </div>
                </div>

                <div
                    class="task-scroll"
                    role="region"
                    aria-label={`Tareas de ${card.name}`}
                >
                    <ul class="tasks">
                        {#each [...card.tasks.filter((task) => !task.done), ...card.tasks.filter((task) => task.done)] as task, i (`${card.id}-${i}`)}
                            <li class:done={task.done}>
                                <span
                                    class:checked={task.done}
                                    class="checkbox"
                                    aria-hidden="true"
                                >
                                    {task.done ? "✓" : ""}
                                </span>

                                <span class="task-title">{task.title}</span>

                                <span
                                    class={`priority priority-${task.priority.toLowerCase()}`}
                                >
                                    {task.priority}
                                </span>
                            </li>
                        {/each}
                    </ul>
                </div>

                <footer class="card-footer">
                    <div class="add-preview">
                        <span class="plus">+</span>
                        <span>Añadir una tarea...</span>
                    </div>
                </footer>
            </article>
        {/each}

        <div class="new-card-preview">
            <span class="plus">+</span>
            Nueva tarjeta
        </div>
    </section>

    <p class="board-hint">
        Desplazá el tablero hacia los lados y las tareas dentro de cada tarjeta
        hacia abajo.
    </p>
</main>

<style>
    .tasks-page {
        display: flex;
        flex-direction: column;
        min-width: 0;
        min-height: 0;
        max-width: none;
        height: 100%;
        padding: 22px 26px 16px;
        overflow: hidden;
    }

    .page-header {
        flex-shrink: 0;
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 20px;
        margin-bottom: 14px;
    }

    .page-header h1 {
        margin: 0;
        font-size: clamp(26px, 3vw, 34px);
    }

    .page-header p:last-child {
        margin: 7px 0 0;
        font-size: 14px;
        color: var(--color-text-muted);
    }

    .prototype-tag {
        flex-shrink: 0;
        border: 1px solid var(--color-border);
        border-radius: var(--radius-control);
        padding: 8px 10px;
        font-size: 12px;
        color: var(--color-text-muted);
    }

    /* Desplazamiento horizontal del tablero */
    .board {
        flex: 1;
        display: flex;
        align-items: stretch;
        gap: 18px;
        min-height: 0;
        min-width: 0;
        overflow-x: auto;
        overflow-y: hidden;
        padding: 2px 2px 14px;
        scrollbar-width: auto;
        overscroll-behavior-x: contain;
        -webkit-overflow-scrolling: touch;
        scrollbar-color: var(--color-border) transparent;
    }

    .context-card {
        display: flex;
        flex-direction: column;
        flex: 0 0 385px;
        min-width: 0;
        min-height: 0;
        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        background: var(--color-surface);
        box-shadow: var(--shadow-card);
        padding: 18px;
    }

    .card-header {
        display: flex;
        gap: 12px;
        align-items: flex-start;
        flex-shrink: 0;
        min-height: 58px;
    }

    .folder-icon {
        display: grid;
        place-items: center;
        flex: 0 0 40px;
        height: 40px;
        border-radius: var(--radius-control);
        background: var(--color-surface-soft);
        color: var(--color-primary);
    }

    .card-heading {
        flex: 1;
        min-width: 0;
    }

    .card-heading h2 {
        font-size: 18px;
        margin: 0 0 5px;
    }

    .card-heading p {
        font-size: 12px;
        line-height: 1.45;
        color: var(--color-text-muted);
        margin: 0;
    }

    .card-more {
        color: var(--color-text-muted);
        font-size: 20px;
        line-height: 1;
    }

    .progress {
        margin: 12px 0;
        flex-shrink: 0;
    }

    .progress-label {
        display: flex;
        justify-content: space-between;
        gap: 8px;
        font-size: 12px;
        color: var(--color-text-muted);
        margin-bottom: 8px;
    }

    .progress-track {
        height: 6px;
        border-radius: 99px;
        overflow: hidden;
        background: var(--color-surface-soft);
    }

    .progress-fill {
        height: 100%;
        border-radius: inherit;
        background: var(--color-primary);
    }

    /* Desplazamiento vertical interno */
    .task-scroll {
        flex: 1;
        min-height: 0;
        overflow-y: auto;
        overscroll-behavior: contain;
        padding-right: 5px;
        scrollbar-width: thin;
        scrollbar-color: var(--color-border) transparent;
    }

    .tasks {
        margin: 0;
        padding: 0;
        display: grid;
        gap: 9px;
        list-style: none;
    }

    .tasks li {
        display: flex;
        align-items: center;
        gap: 9px;
        min-height: 54px;
        padding: 11px;
        min-width: 0;
        border: 1px solid var(--color-border);
        background: var(--color-surface-soft);
        border-radius: var(--radius-control);
    }

    .checkbox {
        display: grid;
        place-items: center;
        flex: 0 0 18px;
        height: 18px;
        border: 1.5px solid var(--color-text-muted);
        border-radius: 5px;
        font-size: 12px;
    }

    .checkbox.checked {
        background: var(--color-primary);
        color: var(--color-background);
        border-color: var(--color-primary);
    }

    .task-title {
        flex: 1;
        min-width: 0;
        overflow-wrap: anywhere;
        font-size: 13px;
        line-height: 1.4;
    }

    .done .task-title {
        text-decoration: line-through;
        color: var(--color-text-muted);
    }

    .priority {
        flex-shrink: 0;
        font-size: 11px;
        font-weight: 600;
        padding: 4px 7px;
        border-radius: 999px;
    }

    .priority-alta {
        border-left: 3px solid #d45d50;
        background: rgba(212, 93, 80, 0.16);
    }

    .priority-media {
        border-left: 3px solid #c29240;
        background: rgba(194, 146, 64, 0.16);
    }

    .priority-baja {
        border-left: 3px solid #599a71;
        background: rgba(89, 154, 113, 0.16);
    }

    /* El botón visual permanece debajo del scroll */
    .card-footer {
        flex-shrink: 0;
        padding-top: 14px;
    }

    .add-preview {
        display: flex;
        align-items: center;
        gap: 9px;
        padding: 13px;
        border: 1px dashed var(--color-border);
        border-radius: var(--radius-control);
        color: var(--color-text-muted);
        font-size: 12px;
    }

    .plus {
        color: var(--color-primary);
        font-size: 20px;
        line-height: 1;
    }

    .new-card-preview {
        align-self: flex-start;
        display: flex;
        align-items: center;
        gap: 10px;
        flex: 0 0 245px;
        padding: 18px;
        border: 1px dashed var(--color-border);
        background: var(--color-surface);
        border-radius: var(--radius-card);
        font-size: 13px;
        color: var(--color-text-muted);
    }

    .board-hint {
        display: none;
    }

    .board-actions {
        display: flex;
        align-items: center;
        flex-shrink: 0;
        gap: 8px;
    }

    .board-arrow {
        display: grid;
        place-items: center;

        width: 36px;
        height: 36px;

        border: 1px solid var(--color-border);
        border-radius: var(--radius-control);

        background: var(--color-surface);
        color: var(--color-text);

        font-size: 24px;
        cursor: pointer;
    }

    .board-arrow:hover {
        background: var(--color-surface-soft);
    }

    @media (max-width: 850px) {
        .tasks-page {
            height: 100%;
            min-height: 0;
            padding: 14px 12px 10px;
        }

        .context-card {
            flex-basis: min(350px, calc(100vw - 32px));
        }

        .board {
            padding-bottom: 14px;
        }
    }
</style>
