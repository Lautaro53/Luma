<script lang="ts">
  import House from "@lucide/svelte/icons/house";
  import ListTodo from "@lucide/svelte/icons/list-todo";
  import CalendarDays from "@lucide/svelte/icons/calendar-days";
  import Search from "@lucide/svelte/icons/search";
  import Trash2 from "@lucide/svelte/icons/trash-2";
  import Settings from "@lucide/svelte/icons/settings";

  import type { PageId } from "../navigation";

  interface Props {
    currentPage: PageId;
    onNavigate: (page: PageId) => void;
  }

  let { currentPage, onNavigate }: Props = $props();
</script>

<aside class="sidebar">
  <div class="brand">
    <img class="brand-logo" src="/luma.svg" alt="" width="36" height="36" />

    <span>
      luma<span class="brand-period">.</span>
    </span>
  </div>

  <p class="eyebrow">Tu espacio</p>

  <nav aria-label="Navegación principal">
    <button
      type="button"
      class="nav-item"
      class:active={currentPage === "home"}
      aria-current={currentPage === "home" ? "page" : undefined}
      onclick={() => onNavigate("home")}
    >
      <House size={19} strokeWidth={1.8} />
      <span>Inicio</span>
    </button>

    <button
      type="button"
      class="nav-item"
      class:active={currentPage === "tasks"}
      aria-current={currentPage === "tasks" ? "page" : undefined}
      onclick={() => onNavigate("tasks")}
    >
      <ListTodo size={19} strokeWidth={1.8} />
      <span>Tareas</span>
    </button>

    <button
      type="button"
      class="nav-item"
      class:active={currentPage === "calendar"}
      aria-current={currentPage === "calendar" ? "page" : undefined}
      onclick={() => onNavigate("calendar")}
    >
      <CalendarDays size={19} strokeWidth={1.8} />
      <span>Calendario</span>
    </button>

    <button
      type="button"
      class="nav-item"
      class:active={currentPage === "search"}
      aria-current={currentPage === "search" ? "page" : undefined}
      onclick={() => onNavigate("search")}
    >
      <Search size={19} strokeWidth={1.8} />
      <span>Buscar</span>
    </button>

    <button
      type="button"
      class="nav-item"
      class:active={currentPage === "trash"}
      aria-current={currentPage === "trash" ? "page" : undefined}
      onclick={() => onNavigate("trash")}
    >
      <Trash2 size={19} strokeWidth={1.8} />
      <span>Papelera</span>
    </button>

    <button
      type="button"
      class="nav-item"
      class:active={currentPage === "settings"}
      aria-current={currentPage === "settings" ? "page" : undefined}
      onclick={() => onNavigate("settings")}
    >
      <Settings size={19} strokeWidth={1.8} />
      <span>Configuración</span>
    </button>
  </nav>

  <div class="sidebar-meta">
    <span class="brand-version" aria-label="Versión 1.0"> v1.0 </span>
  </div>
</aside>

<style>
  .sidebar {
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .sidebar-meta {
    display: flex;
    align-items: center;
    gap: 9px;

    margin-top: auto;
    padding-top: 16px;

    border-top: 1px solid var(--color-border);
    min-width: 0;
    justify-content: flex-start;
  }

  .brand-version {
    display: inline-flex;
    align-items: center;
    justify-content: center;

    margin-left: 0;
    flex-shrink: 0;
    padding: 5px 7px;

    border: 1px solid var(--color-border);
    border-radius: 7px;

    background: var(--color-surface-soft);
    color: var(--color-primary);

    font-size: 10px;
    font-weight: 700;
  }

  @media (max-width: 850px) {
    .sidebar-meta {
      margin-top: 12px;
      padding-top: 12px;
    }
  }
</style>
