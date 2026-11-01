<script lang="ts">
    type CalendarView = "week" | "month";

    let calendarView = $state<CalendarView>("week");

    const weekdays = ["Lun", "Mar", "Mié", "Jue", "Vie", "Sáb", "Dom"];

    let selectedDate = $state(new Date(2026, 10, 2));

    const monthFormatter = new Intl.DateTimeFormat("es-UY", {
        month: "long",
        year: "numeric",
    });

    const shortFormatter = new Intl.DateTimeFormat("es-UY", {
        day: "numeric",
        month: "short",
    });

    function addDays(date: Date, amount: number): Date {
        const result = new Date(date);
        result.setDate(result.getDate() + amount);
        return result;
    }

    function startOfWeek(date: Date): Date {
        const result = new Date(
            date.getFullYear(),
            date.getMonth(),
            date.getDate(),
        );

        result.setDate(result.getDate() - ((result.getDay() + 6) % 7));
        return result;
    }

    const weekStart = $derived(startOfWeek(selectedDate));

    const days = $derived(
        Array.from({ length: 7 }, (_, index) => {
            const date = addDays(weekStart, index);

            return date.toLocaleDateString("es-UY", {
                weekday: "short",
                day: "numeric",
            });
        }),
    );

    const weekTitle = $derived(
        `${shortFormatter.format(weekStart)} – ${shortFormatter.format(
            addDays(weekStart, 6),
        )}`,
    );

    const monthTitle = $derived(monthFormatter.format(selectedDate));

    const isDemoWeek = $derived(
        weekStart.getFullYear() === 2026 &&
            weekStart.getMonth() === 10 &&
            weekStart.getDate() === 2,
    );

    const isDemoMonth = $derived(
        selectedDate.getFullYear() === 2026 && selectedDate.getMonth() === 10,
    );

    const monthCells = $derived.by(() => {
        const first = new Date(
            selectedDate.getFullYear(),
            selectedDate.getMonth(),
            1,
        );

        const offset = (first.getDay() + 6) % 7;

        return Array.from({ length: 42 }, (_, index) => {
            const date = addDays(first, index - offset);

            return {
                day: date.getDate(),
                currentMonth:
                    date.getMonth() === selectedDate.getMonth() &&
                    date.getFullYear() === selectedDate.getFullYear(),
            };
        });
    });

    function movePeriod(direction: number): void {
        if (calendarView === "week") {
            selectedDate = addDays(weekStart, direction * 7);
        } else {
            selectedDate = new Date(
                selectedDate.getFullYear(),
                selectedDate.getMonth() + direction,
                1,
            );
        }
    }

    function goToToday(): void {
        selectedDate = new Date();
    }

    type CalendarTask = {
        day: number;
        slot: number;
        title: string;
        time: string;
        priority: "alta" | "media" | "baja";
    };

    const hours = ["08:00", "10:00", "12:00", "14:00", "16:00", "18:00"];

    // Información ficticia para el diseño.
    const tasks: CalendarTask[] = [
        {
            day: 0,
            slot: 1,
            title: "Estudiar Bases de Datos",
            time: "10:00 · 90 min",
            priority: "alta",
        },
        {
            day: 1,
            slot: 0,
            title: "Diseñar Proyecto Luma",
            time: "08:00 · 2 h",
            priority: "media",
        },
        {
            day: 2,
            slot: 2,
            title: "Preparar presentación",
            time: "12:00 · 60 min",
            priority: "alta",
        },
        {
            day: 3,
            slot: 3,
            title: "Actualizar currículum",
            time: "14:00 · 45 min",
            priority: "baja",
        },
        {
            day: 4,
            slot: 1,
            title: "Repasar apuntes",
            time: "10:00 · 90 min",
            priority: "media",
        },
    ];
</script>

<main class="content calendar-page">
    <header class="calendar-header">
        <div>
            <p class="eyebrow">TU ESPACIO</p>
            <h1>Calendario</h1>
            <p>Planificá tu semana y organizá tu tiempo.</p>
        </div>

        <div class="view-selector" aria-label="Vista del calendario">
            <button
                type="button"
                class:active={calendarView === "week"}
                aria-pressed={calendarView === "week"}
                onclick={() => (calendarView = "week")}
            >
                Semana
            </button>

            <button
                type="button"
                class:active={calendarView === "month"}
                aria-pressed={calendarView === "month"}
                onclick={() => (calendarView = "month")}
            >
                Mes
            </button>
        </div>
    </header>

    <section class="calendar-panel">
        <div class="calendar-toolbar">
            <div>
                <h2>
                    {calendarView === "week" ? weekTitle : monthTitle}
                </h2>
                <p>Calendario de ejemplo · Datos ficticios</p>
            </div>

            <div class="calendar-actions">
                <button
                    type="button"
                    aria-label="Período anterior"
                    onclick={() => movePeriod(-1)}
                >
                    ‹
                </button>

                <button type="button" onclick={goToToday}> Hoy </button>

                <button
                    type="button"
                    aria-label="Período siguiente"
                    onclick={() => movePeriod(1)}
                >
                    ›
                </button>
            </div>
        </div>

        {#if calendarView === "week"}
            <div class="calendar-scroll">
                <div class="calendar-grid">
                    <div class="time-heading">Hora</div>

                    {#each days as day}
                        <div class="day-heading">{day}</div>
                    {/each}

                    {#each hours as hour, row}
                        <div class="hour-label">{hour}</div>

                        {#each days as day, column}
                            {@const events = tasks.filter(
                                (task) =>
                                    isDemoWeek &&
                                    task.day === column &&
                                    task.slot === row,
                            )}

                            <div class="time-slot">
                                {#each events as task}
                                    <div
                                        class={`calendar-task priority-${task.priority}`}
                                    >
                                        <strong>{task.title}</strong>
                                        <span>{task.time}</span>
                                    </div>
                                {/each}
                            </div>
                        {/each}
                    {/each}
                </div>
            </div>
        {:else}
            <div class="calendar-scroll">
                <div class="month-grid">
                    {#each weekdays as weekday}
                        <div class="month-heading">{weekday}</div>
                    {/each}

                    {#each monthCells as cell, index (index)}
                        <div
                            class="month-day"
                            class:outside={!cell.currentMonth}
                        >
                            <span class="month-number">{cell.day}</span>

                            {#if cell.currentMonth}
                                {#each tasks.filter((task) => isDemoMonth && task.day + 2 === cell.day) as task}
                                    <div
                                        class={`month-event priority-${task.priority}`}
                                    >
                                        <strong>{task.title}</strong>
                                        <span>{task.time}</span>
                                    </div>
                                {/each}
                            {/if}
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    </section>

    <p class="calendar-note">
        Vista visual de ejemplo. Las tareas todavía no están conectadas a las
        tarjetas ni al backend.
    </p>
</main>

<style>
    .calendar-page {
        display: flex;
        flex-direction: column;

        width: 100%;
        max-width: 1450px;
        height: 100%;

        min-width: 0;
        min-height: 0;

        padding-top: 12px;

        overflow-x: hidden;
        overflow-y: auto;
    }

    .calendar-header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        gap: 16px;
        margin-bottom: 10px;
    }

    .calendar-header h1 {
        margin: 0 0 4px;
        font-size: clamp(26px, 2.5vw, 32px);
    }

    .calendar-header .eyebrow {
        display: none;
    }

    .calendar-header p:last-child {
        margin: 0;
        font-size: 14px;
        color: var(--color-text-muted);
    }

    .view-selector {
        display: flex;
        gap: 4px;
        padding: 4px;
        border: 1px solid var(--color-border);
        border-radius: var(--radius-control);
        background: var(--color-surface);
    }

    .view-selector button {
        padding: 9px 16px;
        border: none;
        border-radius: var(--radius-control);
        background: transparent;
        color: var(--color-text-muted);
        font-size: 13px;
        cursor: pointer;
    }

    .view-selector button.active {
        background: var(--color-surface-soft);
        color: var(--color-primary);
        font-weight: 600;
    }

    .view-selector button:hover:not(.active) {
        background: var(--color-surface-soft);
    }

    .calendar-panel {
        display: flex;
        flex-direction: column;
        flex: 0 0 auto;
        min-height: 0;
        min-width: 0;

        overflow: hidden;

        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        background: var(--color-surface);
        box-shadow: var(--shadow-card);
    }

    .calendar-toolbar {
        display: flex;
        align-items: center;
        justify-content: space-between;

        padding: 12px 16px;
        gap: 12px;

        border-bottom: 1px solid var(--color-border);
    }

    .calendar-toolbar h2 {
        margin: 0;
        font-size: 18px;
    }

    .calendar-toolbar p {
        display: none;
    }

    .demo-label,
    .calendar-note {
        color: var(--color-text-muted);
        font-size: 12px;
    }

    .calendar-scroll {
        flex: none;

        height: clamp(420px, calc(100dvh - 190px), 860px);

        min-width: 0;
        min-height: 0;

        overflow: auto;
        overscroll-behavior: contain;

        scrollbar-width: thin;
        scrollbar-color: var(--color-border) transparent;
    }
    
    .calendar-grid {
        display: grid;

        grid-template-columns: 65px repeat(7, minmax(110px, 1fr));
        grid-template-rows: 44px repeat(6, minmax(100px, 1fr));

        min-width: 840px;
        min-height: 644px;
        height: 100%;
    }

    .time-heading,
    .day-heading {
        padding: 15px 8px;
        text-align: center;
        font-size: 12px;
        font-weight: 600;
        background: var(--color-surface-soft);
        border-bottom: 1px solid var(--color-border);
    }

    .day-heading {
        border-left: 1px solid var(--color-border);
    }

    .hour-label {
        padding: 9px 6px;
        font-size: 11px;
        text-align: center;
        color: var(--color-text-muted);
        border-bottom: 1px solid var(--color-border);
    }

    .time-slot {
        min-width: 0;
        min-height: 94px;
        padding: 5px;
        border-left: 1px solid var(--color-border);
        border-bottom: 1px solid var(--color-border);
    }

    .calendar-task {
        display: flex;
        flex-direction: column;
        gap: 6px;
        padding: 9px;
        border-radius: var(--radius-control);
        background: var(--color-surface-soft);
        overflow-wrap: anywhere;
    }

    .calendar-task strong {
        font-size: 12px;
        line-height: 1.35;
        font-weight: 600;
    }

    .calendar-task span {
        font-size: 11px;
        color: var(--color-text-muted);
    }

    .priority-alta {
        border-left: 3px solid #d45d50;
    }

    .priority-media {
        border-left: 3px solid #c29240;
    }

    .priority-baja {
        border-left: 3px solid #599a71;
    }

    .calendar-note {
        flex-shrink: 0;
        margin: 8px 0 0;
        color: var(--color-text-muted);
        font-size: 12px;
    }

    .month-grid {
        display: grid;

        grid-template-columns: repeat(7, minmax(118px, 1fr));
        grid-template-rows: 42px repeat(6, minmax(76px, 1fr));

        min-width: 826px;
        height: 100%;
        min-height: 498px;
    }

    .month-heading {
        padding: 15px 8px;
        text-align: center;
        font-size: 12px;
        font-weight: 600;
        background: var(--color-surface-soft);
        border-bottom: 1px solid var(--color-border);
    }

    .month-day {
        min-width: 0;
        min-height: 0;
        overflow-y: auto;
        padding: 9px;
        border-right: 1px solid var(--color-border);
        border-bottom: 1px solid var(--color-border);
    }

    .month-day.outside {
        opacity: 0.4;
    }

    .month-number {
        display: block;
        margin-bottom: 9px;
        font-size: 12px;
        font-weight: 600;
        color: var(--color-text-muted);
    }

    .month-event {
        display: flex;
        flex-direction: column;
        gap: 5px;
        padding: 8px;
        margin-bottom: 5px;
        border-radius: var(--radius-control);
        background: var(--color-surface-soft);
        overflow-wrap: anywhere;
    }

    .month-event strong {
        font-size: 11px;
        line-height: 1.4;
    }

    .month-event span {
        font-size: 10px;
        color: var(--color-text-muted);
    }

    .calendar-actions {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .calendar-actions button {
        min-width: 38px;
        height: 36px;
        padding: 0 12px;
        border: 1px solid var(--color-border);
        border-radius: var(--radius-control);
        background: var(--color-surface-soft);
        color: var(--color-text);
        font-size: 14px;
        cursor: pointer;
    }

    .calendar-actions button:hover {
        color: var(--color-primary);
        border-color: var(--color-primary);
    }

    @media (max-width: 600px) {
        .calendar-toolbar {
            flex-wrap: wrap;
        }
    }

    @media (max-width: 700px) {
        .calendar-header {
            flex-direction: column;
            align-items: flex-start;
        }

        .calendar-toolbar {
            padding: 16px;
        }

        .demo-label {
            display: none;
        }
    }
</style>
